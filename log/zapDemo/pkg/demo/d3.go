package demo

import (
	"os"

	"zapDemo/pkg/log"
)

// 写入多个文件
func Demo3() {
	file1, err := os.OpenFile("./access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	file2, err := os.OpenFile("./error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
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
				return lvl > log.InfoLevel
			},
		},
	}

	logger := log.NewTee(tops)
	log.ResetDefault(logger)

	log.Info("demo3:", log.String("app", "start ok"),
		log.Int("major version", 3))

	log.Error("demo3:", log.String("app", "crash"),
		log.Int("reason", -1))
}
