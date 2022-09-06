package demo

import (
	"os"

	"github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"
)

func Demo8() {
	file1, err := os.OpenFile("./access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	p, err := log.NewKafkaAsyncProducer([]string{"127.0.0.1:29092"})
	kafkaWriter := log.NewKafakSyncer(p, "test", file1)
	logger := log.New(kafkaWriter, log.Level(log.InfoLevel))
	log.ResetDefault(logger)

	log.Info("demo8:", log.String("app", "start ok"),
		log.Int("major version", 8))
}
