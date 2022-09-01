package demo

import "zapDemo/pkg/log"

// 日志写到文件，并发送到http服务
func Demo5() {
	l := log.NewProduction([]string{"./demo.log", "http://localhost:8080/"})

	log.ResetDefault(l)

	defer log.Sync()

	log.Info("default", log.String("key", "value"),
		log.Int("count", 4))
}
