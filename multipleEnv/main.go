package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadFromEnv() {
	viper.AutomaticEnv()
}

func main() {
	ReadFromEnv()

	// 环境变量要全大写，viper才能读到
	fmt.Println("SERVER_MODE: ", viper.Get("SERVER_MODE"))
	fmt.Println("CONSUL_ADDR ", viper.Get("CONSUL_ADDR"))
}
