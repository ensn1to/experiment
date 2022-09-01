package demo

import (
	"os"

	"github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"
)

// 写入文件，并写到redis
func Demo6() {
	file1, err := os.OpenFile("./access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	file2, err := log.NewRedisWriter("log_list", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}

	tops := []log.TeeOption{
		{
			W: file1,
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
		{
			W: file2,
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
	}

	logger := log.NewTee(tops)
	log.ResetDefault(logger)

	for i := 0; i < 200; i++ {
		log.Info("demo4:", log.String("app", "start ok"),
			log.Int("major version", i))
	}
}
