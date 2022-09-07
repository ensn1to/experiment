package gorabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishOption struct {
	ReconnectInterval time.Duration
}

type Publisher struct {
	chManager *channelManager

	notfiyReturnChan  chan amqp.Return
	notfiyPublishChan chan amqp.Confirmation

	disablePublishDueToFlow    bool
	disablePublishDueToFlowMux *sync.RWMutex

	options PublishOption
}

func WithPublisherOptionsReconnectInterval(reconnectInterval time.Duration) func(options *PublishOption) {
	return func(options *PublishOption) {
		options.ReconnectInterval = reconnectInterval
	}
}

func NewPublisher(url string, config Config, optsFuncs ...func(*PublishOption)) (*Publisher, error) {
	opts := &PublishOption{
		ReconnectInterval: time.Second * 5,
	}
	for _, f := range optsFuncs {
		f(opts)
	}

	chManager, err := newChannelManager(url, amqp.Config(config), opts.ReconnectInterval)
	if err != nil {
		return nil, err
	}

	publish := &Publisher{
		chManager:                  chManager,
		disablePublishDueToFlow:    false,
		disablePublishDueToFlowMux: &sync.RWMutex{},
		options:                    *opts,
		notfiyReturnChan:           nil,
		notfiyPublishChan:          nil,
	}

	go publish.startNotifyFlowHandler()

	go publish.handleRestarts()

	return publish, nil
}

func (p *Publisher) handleRestarts() {
	for err := range p.chManager.notifyCancelOrClose {
		fmt.Printf("successful publish recovery from : %s", err.Error())
		go p.startNotifyFlowHandler()

		if p.notfiyReturnChan != nil {
			go p.startNotifyReturnHandler()
		}

		if p.notfiyPublishChan != nil {
			go p.startNotifyPublishHandler()
		}
	}
}

func (p *Publisher) startNotifyFlowHandler() {
	notifyFlowChan := p.chManager.channel.NotifyFlow(make(chan bool))
	p.disablePublishDueToFlowMux.Lock()
	p.disablePublishDueToFlow = false
	p.disablePublishDueToFlowMux.Unlock()

	//
	for ok := range notifyFlowChan {
		p.disablePublishDueToFlowMux.Lock()
		if ok {
			fmt.Print("pasuing publishing due to flow request from server")
			p.disablePublishDueToFlow = true
		} else {
			p.disablePublishDueToFlow = false
			fmt.Print("resuming publishing due to flow request from server")
		}
		p.disablePublishDueToFlowMux.Unlock()
	}
}

func (p *Publisher) startNotifyReturnHandler() {
	retrunCh := p.chManager.channel.NotifyReturn(make(chan amqp.Return, 1))

	for ret := range retrunCh {
		p.notfiyReturnChan <- ret
	}
}

func (p *Publisher) NotifyReturn() <-chan amqp.Return {
	p.notfiyReturnChan = make(chan amqp.Return)
	go p.startNotifyReturnHandler()
	return p.notfiyReturnChan
}

func (p *Publisher) NotifyPublish() <-chan amqp.Confirmation {
	p.notfiyPublishChan = make(chan amqp.Confirmation)
	go p.startNotifyPublishHandler()
	return p.notfiyPublishChan
}

func (p *Publisher) startNotifyPublishHandler() {
	p.chManager.channel.Confirm(false)
	publishCh := p.chManager.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	for conf := range publishCh {
		p.notfiyPublishChan <- conf
	}
}

func (p *Publisher) Pulish(ctx context.Context, data []byte, routingKey string, optsFunc ...PublishOptFunc) error {
	routingKeys := []string{routingKey}
	return p.PublishMultiple(ctx, data, routingKeys, optsFunc...)
}

func (p *Publisher) PublishMultiple(ctx context.Context, data []byte, routingKeys []string, optsFunc ...PublishOptFunc) error {
	p.disablePublishDueToFlowMux.RLock()
	if p.disablePublishDueToFlow {
		return errors.New("publishing blocked due to high flow on the server")
	}
	p.disablePublishDueToFlowMux.RUnlock()

	options := &PublishOptions{}
	for _, optionFunc := range optsFunc {
		optionFunc(options)
	}
	if options.DeliveryMode == 0 {
		options.DeliveryMode = amqp.Transient
	}
	for _, routingKey := range routingKeys {
		message := amqp.Publishing{}
		message.ContentType = options.ContentType
		message.DeliveryMode = options.DeliveryMode
		message.Body = data
		message.Headers = options.Headers
		message.Expiration = options.Expiration
		message.ContentEncoding = options.ContentEncoding
		message.Priority = options.Priority
		message.CorrelationId = options.CorrelationID
		message.ReplyTo = options.ReplyTo
		message.MessageId = options.MessageID
		message.Timestamp = options.Timestamp
		message.Type = options.Type
		message.UserId = options.UserID
		message.AppId = options.AppID

		// Actual publish.
		err := p.chManager.channel.Publish(
			options.Exchange,
			routingKey,
			options.Mandatory,
			options.Immediate,
			message,
		)
		if err != nil {
			fmt.Printf("Message publishing falied: %v", err)
			return err
		}
	}
	return nil
}

func (p *Publisher) StopPublish() error {
	if err := p.chManager.close(); err != nil {
		return err
	}
	return nil
}
