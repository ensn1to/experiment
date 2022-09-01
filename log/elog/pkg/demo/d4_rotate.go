package demo

import "github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"

// 日志分片
func Demo4() {
	tops := []log.TeeOptionWithRotate{
		{
			Filename: "access.log",
			Ropt: log.RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
		{
			Filename: "error.log",
			Ropt: log.RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(lvl log.Level) bool {
				return lvl > log.InfoLevel
			},
		},
	}

	logger := log.NewTeeWithRotate(tops)
	log.ResetDefault(logger)

	for i := 0; i < 20000; i++ {
		log.Info("demo4:", log.String("app", "start ok"),
			log.Int("major version", 3))
		log.Error("demo4:", log.String("app", "crash"),
			log.Int("reason", -1))
	}
}
