package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	c := http.DefaultClient

	// http.Client底层的数据连接建立和维护是由http.Transport实现的
	c.Transport = &http.Transport{
		DisableKeepAlives: true,
	}
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

	}
}
