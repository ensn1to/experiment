package log

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap/zapcore"
)

type kafkaWriterSyncer struct {
	topic          string
	producer       sarama.AsyncProducer
	fallbackSyncer zapcore.WriteSyncer
}

func NewKafkaAsyncProducer(addrs []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true

	return sarama.NewAsyncProducer(addrs, config)
}

func NewKafakSyncer(producer sarama.AsyncProducer, topic string, fallbackWs zapcore.WriteSyncer) zapcore.WriteSyncer {
	w := &kafkaWriterSyncer{
		topic:          topic,
		producer:       producer,
		fallbackSyncer: zapcore.AddSync(fallbackWs),
	}

	go func() {
		for e := range producer.Errors() {
			val, err := e.Msg.Value.Encode()
			if err != nil {
				continue
			}

			// 把失败的数据写入文件
			fallbackWs.Write(val)
		}
	}()

	return w
}

func (ws *kafkaWriterSyncer) Write(b []byte) (n int, err error) {
	b1 := make([]byte, len(b))
	copy(b1, b) // b is reused, we must pass its copy b1 to sarama
	msg := &sarama.ProducerMessage{
		Topic: ws.topic,
		Value: sarama.ByteEncoder(b1),
	}
	ws.producer.Input() <- msg

	// sarama会hang住，如果不会可直接用以上的方式
	// https://github.com/Shopify/sarama/pull/2133
	/*
		select {
		case ws.producer.Input() <- msg:
		default:
			// if producer block on input channel, write log entry to default fallbackSyncer
			return ws.fallbackSyncer.Write(b1)
		}
	*/

	return len(b1), nil
}

func (ws *kafkaWriterSyncer) Sync() error {
	ws.producer.AsyncClose()
	return ws.fallbackSyncer.Sync()
}
