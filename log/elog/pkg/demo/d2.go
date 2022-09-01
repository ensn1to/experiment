package demo

import (
	"os"

	"zapDemo/pkg/log"
)

// write into file
func Demo2() {
	file, err := os.OpenFile("./demo.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}

	l := log.New(file, log.InfoLevel, log.WithCaller(true), log.AddStacktrace(log.InfoLevel))

	log.ResetDefault(l)

	defer log.Sync()

	log.Info("default", log.String("key", "value"),
		log.Int("count", 4))
}
