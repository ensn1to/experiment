package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	c := http.DefaultClient

	for i := 0; i < 2; i++ {
		resp1, err := c.Get("http://localhost:18081")
		if err != nil {
			panic(err)
		}

		defer resp1.Body.Close()

		b1, err := ioutil.ReadAll(resp1.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("resp from server-1%s\n", b1)

		resp2, err := c.Get("http://localhost:18082")
		if err != nil {
			panic(err)
		}

		defer resp1.Body.Close()

		b2, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("resp from server-2%s\n", b2)

	}
}
