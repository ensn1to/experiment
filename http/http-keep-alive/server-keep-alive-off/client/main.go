package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	c := http.DefaultClient

	for i := 0; i < 5; i++ {
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
