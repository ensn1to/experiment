package demo

import (
	"os"

	"github.com/ensn1to/experiment/tree/master/log/elog/pkg/log"
)

// 运行过程中动态调整日志级别：info > error
func Demo7() {
	l := log.New(os.Stdout, log.ErrorLevel)

	log.ResetDefault(l)

	defer log.Sync()

	log.Info("info", log.String("d7", "demo7 before setLevel"))

	l.SetLevel(log.InfoLevel)

	log.Info("info", log.String("d7", "demo7 after setLevel"))
}
