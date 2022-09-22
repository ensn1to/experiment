package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadFromEnv() {
	// viper.AutomaticEnv()
	viper.AutomaticEnv()
}

func main() {
	ReadFromEnv()

	fmt.Println("SERVER_MODE: ", viper.Get("SERVER_MODE"))
	fmt.Println("CONSUL_ADDR ", viper.Get("CONSUL_ADDR"))
}
