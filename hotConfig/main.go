package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func prepare() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("prepare host : %s\n", viper.Get("redis"))
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("检测到配置更改...")
		fmt.Printf("filechange host : %s\n", viper.Get("redis"))
	})
}

func main() {
	prepare()

	watchConfig()

	time.Sleep(time.Second * 30)
}
