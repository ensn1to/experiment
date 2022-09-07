package demo

import (
	"os"

	"github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"
)

func Demo9() {
	file1, err := os.OpenFile("./access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	p, err := log.NewRabbitAsyncProducer("127.0.0.1:5672")
	kafkaWriter := log.NewRabbitSyncer(p, "logstash_processing_queue", file1)
	logger := log.New(kafkaWriter, log.Level(log.InfoLevel))
	log.ResetDefault(logger)

	log.Info("demo9:", log.String("app", "start ok"),
		log.Int("major version", 9))
}
