package gorabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumeOptFunc func(*ConsumeOptions)

// includes the options of the queue, queueBind, exchange, consume and qos
type ConsumeOptions struct {
	QueueDurable      bool
	QueueAutoDelete   bool
	QueueExclusive    bool
	QueueNoWait       bool
	QueueDeclare      bool
	QueueArgs         amqp.Table
	BindingExchange   *BindingExchangeOptions
	BindingNoWait     bool
	BindingArgs       amqp.Table
	ConsumerName      string
	ConsumerAutoAck   bool
	ConsumerExclusive bool
	ConsumerNoWait    bool
	ConsumerNoLocal   bool
	ConsumerArgs      amqp.Table

	// numbers of the goroutines to consume
	Concurrency int
	QOSPrefetch int
	QOSGlobal   bool
}

func newDefaultConsumeOptions() *ConsumeOptions {
	return &ConsumeOptions{
		QueueDurable:      false,
		QueueAutoDelete:   false,
		QueueExclusive:    false,
		QueueNoWait:       false,
		QueueDeclare:      true,
		QueueArgs:         nil,
		BindingExchange:   nil,
		BindingNoWait:     false,
		BindingArgs:       nil,
		ConsumerName:      "",
		ConsumerAutoAck:   false,
		ConsumerExclusive: false,
		ConsumerNoWait:    false,
		ConsumerNoLocal:   false,
		ConsumerArgs:      nil,
		Concurrency:       2,
		QOSPrefetch:       0,
		QOSGlobal:         false,
	}
}

type BindingExchangeOptions struct {
	Name         string
	Kind         string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	ExchangeArgs Table
	Declare      bool
}

func newDefaultBindingExchangeOptions(consumerOpts *ConsumeOptions) *BindingExchangeOptions {
	if consumerOpts.BindingExchange == nil {
		consumerOpts.BindingExchange = &BindingExchangeOptions{
			Name:         "",
			Kind:         "topic",
			Durable:      false,
			AutoDelete:   false,
			Internal:     false,
			NoWait:       false,
			ExchangeArgs: nil,
			Declare:      false,
		}
	}
	return consumerOpts.BindingExchange
}

func WithConsumeOptionsQueueDurable(options *ConsumeOptions) {
	options.QueueDurable = true
}

func WithConsumeOptionsQueueNoWait(options *ConsumeOptions) {
	options.QueueNoWait = true
}

func WithConsumeOptionsBindingExchangeName(name string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		newDefaultBindingExchangeOptions(options).Name = name
	}
}

func WithConsumeOptionsBindingExchangeKind(kind string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		newDefaultBindingExchangeOptions(options).Kind = kind
	}
}

func WithConsumeOptionsBindingExchangeDurable(options *ConsumeOptions) {
	newDefaultBindingExchangeOptions(options).Durable = true
}

func WithConsumeOptionsBindingExchangeAutoDelete(options *ConsumeOptions) {
	newDefaultBindingExchangeOptions(options).AutoDelete = true
}

func WithConsumeOptionsBindingExchangeInternal(options *ConsumeOptions) {
	newDefaultBindingExchangeOptions(options).Internal = true
}

func WithConsumeOptionsBindingExchangeNoWait(options *ConsumeOptions) {
	newDefaultBindingExchangeOptions(options).NoWait = true
}

func WithConsumeOptionsBindingExchangeArgs(args Table) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		newDefaultBindingExchangeOptions(options).ExchangeArgs = args
	}
}

func WithConsumeOptionsBindingNoWait(options *ConsumeOptions) {
	options.BindingNoWait = true
}

func WithConsumeOptionsConcurrency(concurrency int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.Concurrency = concurrency
	}
}

func WithConsumeOptionsQOSPrefetch(prefetchCount int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.QOSPrefetch = prefetchCount
	}
}

func WithConsumeOptionsQOSGlobal(options *ConsumeOptions) {
	options.QOSGlobal = true
}

func WithConsumeOptionsConsumerName(consumerName string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerName = consumerName
	}
}

func WithConsumeOptionsConsumerAutoAck(autoAck bool) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerAutoAck = autoAck
	}
}

func WithConsumeOptionsConsumerExclusive(options *ConsumeOptions) {
	options.ConsumerExclusive = true
}

func WithConsumeOptionsConsumerNoWait(options *ConsumeOptions) {
	options.ConsumerNoWait = true
}

func WithConsumeOptionsQueueArgs(args Table) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.QueueArgs = amqp.Table(args)
	}
}
