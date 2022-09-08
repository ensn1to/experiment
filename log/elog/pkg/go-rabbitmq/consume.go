package gorabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Action is an action that occurs after processed this delivery
type Action int

const (
	// Ack default ack this msg after you have successfully processed this delivery.
	Ack Action = iota
	// NackDiscard the message will be dropped or delivered to a server configured dead-letter queue.
	NackDiscard
	// NackRequeue deliver this message to a different consumer.
	NackRequeue
)

type Consumer struct {
	chManager     *channelManager
	globalOptions []ConsumeOptFunc
}
type ConsumerOptions struct {
	ReconnectInterval time.Duration
}

type (
	Config     amqp.Config
	OptionFunc func(*ConsumerOptions)
)

type Delivery struct {
	amqp.Delivery
}

type Handler func(d Delivery) (action Action)

func NewConsumer(url string, config Config, optsFuncs ...OptionFunc) (*Consumer, error) {
	options := &ConsumerOptions{
		ReconnectInterval: time.Second * 5,
	}
	for _, f := range optsFuncs {
		f(options)
	}

	chManager, err := newChannelManager(url, amqp.Config(config), options.ReconnectInterval)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		chManager:     chManager,
		globalOptions: make([]ConsumeOptFunc, 0),
	}, nil
}

func (c *Consumer) Disconnect() (err error) {
	if err = c.chManager.channel.Close(); err != nil {
		return err
	}
	if c.chManager.connection.Close(); err != nil {
		return err
	}

	return nil
}

func WithConsumerOptionsReconnectInterval(ReconnectInterval time.Duration) OptionFunc {
	return func(co *ConsumerOptions) {
		co.ReconnectInterval = ReconnectInterval
	}
}

func (c *Consumer) Consume(ctx context.Context, queue string, routingKey string, handler Handler, optsFuncs ...ConsumeOptFunc) error {
	routingKeys := []string{routingKey}
	return c.ConsumeMultiple(ctx, queue, routingKeys, handler, optsFuncs...)
}

func (c *Consumer) AddCommonOption(optionFunc ...ConsumeOptFunc) {
	for _, optsFunc := range optionFunc {
		if optsFunc != nil {
			c.globalOptions = append(c.globalOptions, optsFunc)
		}
	}
}

func (c *Consumer) ConsumeMultiple(ctx context.Context, queue string, routingKeys []string, handler Handler, optsFuncs ...ConsumeOptFunc) error {
	consumOpts := newDefaultConsumeOptions()

	if len(c.globalOptions) > 0 {
		optsFuncs = append(optsFuncs, c.globalOptions...)
	}
	for _, f := range optsFuncs {
		f(consumOpts)
	}

	if err := c.startConsuming(queue, routingKeys, handler, consumOpts); err != nil {
		return err
	}

	// blocked
	go func() {
		for err := range c.chManager.notifyCancelOrClose {
			fmt.Printf("successful recover from: %s", err.Error())
			// restart consume
			if err := c.startConsuming(queue, routingKeys, handler, consumOpts); err != nil {
				fmt.Printf("error restarting consumer goroutines after cancel or close: %s", err.Error())
			}
		}
	}()

	return nil
}

func (c *Consumer) StopConsume(consumerName string, noWait bool) {
	c.chManager.channel.Cancel(consumerName, noWait)
}

// operation: QueueDeclare -> ExchangeDecalre -> QueueBind -> Consume
func (c *Consumer) startConsuming(queue string, routingKeys []string, handler Handler, consumeOpts *ConsumeOptions) error {
	c.chManager.channelMux.RLock()
	defer c.chManager.channelMux.RUnlock()

	if consumeOpts.QueueDeclare {
		if _, err := c.chManager.channel.QueueDeclare(
			queue,
			consumeOpts.QueueDurable,
			consumeOpts.QueueAutoDelete,
			consumeOpts.QueueExclusive,
			consumeOpts.QueueNoWait,
			consumeOpts.ConsumerArgs,
		); err != nil {
			return err
		}
	}

	if consumeOpts.BindingExchange != nil {
		exchange := consumeOpts.BindingExchange
		if exchange.Name == "" {
			return errors.New("binding to exchange but name not specified")
		}

		if exchange.Declare {
			if err := c.chManager.channel.ExchangeDeclare(
				exchange.Name,
				exchange.Kind,
				exchange.Durable,
				exchange.AutoDelete,
				exchange.Internal,
				exchange.NoWait,
				amqp.Table(exchange.ExchangeArgs)); err != nil {
				return err
			}
		}

		for _, routingKey := range routingKeys {
			if err := c.chManager.channel.QueueBind(
				queue,
				routingKey,
				exchange.Name,
				consumeOpts.BindingNoWait,
				consumeOpts.BindingArgs); err != nil {
				return err
			}
		}
	}

	err := c.chManager.channel.Qos(
		consumeOpts.QOSPrefetch,
		0,
		consumeOpts.QOSGlobal,
	)
	if err != nil {
		return err
	}

	msgs, err := c.chManager.channel.Consume(
		queue,
		consumeOpts.ConsumerName,
		consumeOpts.ConsumerAutoAck,
		consumeOpts.ConsumerExclusive,
		consumeOpts.ConsumerNoLocal,
		consumeOpts.ConsumerNoWait,
		consumeOpts.ConsumerArgs,
	)
	if err != nil {
		return err
	}

	for i := 0; i < consumeOpts.Concurrency; i++ {
		go handlerMsg(msgs, consumeOpts, handler)
	}

	fmt.Printf("Processing messages on %v goroutines", consumeOpts.Concurrency)

	return nil
}

func handlerMsg(msgs <-chan amqp.Delivery, consumeOptions *ConsumeOptions, handler Handler) {
	for msg := range msgs {
		if consumeOptions.ConsumerAutoAck {
			handler(Delivery{msg})
			continue
		}

		switch handler(Delivery{msg}) {
		case Ack:
			if err := msg.Ack(false); err != nil {
				fmt.Printf("cann't ack message: %s", err.Error())
			}
		case NackDiscard:
			err := msg.Nack(false, false)
			if err != nil {
				fmt.Printf("can't nack message: %v", err)
			}
		case NackRequeue:
			err := msg.Nack(false, true)
			if err != nil {
				fmt.Printf("can't nack message: %v", err)
			}
		}

	}

	fmt.Printf("rabbit consumer handler goroutine closed")
}
