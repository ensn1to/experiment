package log

import (
	"context"

	gorabbit "github.com/ensn1to/experiment/tree/master/log/elog/pkg/go-rabbitmq"
	"go.uber.org/zap/zapcore"
)

type RabbitWriterSyncer struct {
	routingKey     string
	producer       *gorabbit.Publisher
	fallbackSyncer zapcore.WriteSyncer
}

func NewRabbitAsyncProducer(addr string) (*gorabbit.Publisher, error) {
	cfg := gorabbit.Config{
		Vhost: "/log",
	}
	return gorabbit.NewPublisher(addr, cfg)
}

func NewRabbitSyncer(producer *gorabbit.Publisher, routingKey string, fallbackWs zapcore.WriteSyncer) zapcore.WriteSyncer {
	w := &RabbitWriterSyncer{
		routingKey:     routingKey,
		producer:       producer,
		fallbackSyncer: zapcore.AddSync(fallbackWs),
	}

	return w
}

func (ws *RabbitWriterSyncer) Write(b []byte) (n int, err error) {
	b1 := make([]byte, len(b))
	copy(b1, b) // b is reused

	if err := ws.producer.Pulish(context.TODO(), b1, ws.routingKey); err != nil {
		ws.fallbackSyncer.Write(b1)
	}

	return len(b1), nil
}

func (ws *RabbitWriterSyncer) Sync() error {
	ws.producer.StopPublish()
	return ws.fallbackSyncer.Sync()
}
