package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	c := http.DefaultClient

	for i := 0; i < 2; i++ {
		resp, err := c.Get("http://localhost:18081")
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s\n", b)

		// 等待7s
		time.Sleep(7 * time.Second)
	}
}
