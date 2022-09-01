package demo

import "github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"

// 封装成类似log方式
func Demo1() {
	log.Info("call demo", log.String("url", "http://localhost"),
		log.Int("retries", 3))
}
