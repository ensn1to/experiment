package main

import (
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	fp "path/filepath"
	"time"
)

var (
	filepath string
	addr     string
)

func init() {
	flag.StringVar(&filepath, "filepath", "", "the file to upload")
	flag.StringVar(&addr, "addr", "localhost:18080", "file server")
	flag.Parse()
}

func main() {
	if filepath == "" {
		fmt.Println("filepath is empty")
		return
	}

	if addr == "" {
		fmt.Println("file server addr is empty")
		return
	}
	start := time.Now()
	err := doUpload(addr, filepath)
	if err != nil {
		return
	}
	total := time.Since(start)
	fmt.Printf("upload file [%s] ok, total: %v\n", filepath, total)
}

func createReqBody(filepath string) (string, io.Reader, error) {
	var err error

	pr, pw := io.Pipe()

	// buf := new(bytes.Buffer)
	bw := multipart.NewWriter(pw)

	f, err := os.Open(filepath)
	if err != nil {
		return "", nil, err
	}
	// io.Pipe() based on no-buffer channel, need anthoer gorotine in case of main goroutine blocked
	go func() {
		defer f.Close()

		// text part1
		p1w, _ := bw.CreateFormField("name")
		p1w.Write([]byte("Hello, world!"))

		// text part2
		p2w, _ := bw.CreateFormField("age")
		p2w.Write([]byte("22"))

		// file part1
		_, fileName := fp.Split(filepath)
		fw1, _ := bw.CreateFormFile("file1", fileName)

		fw1, _ = bw.CreateFormField(fileName)
		io.Copy(fw1, f) // default buffer size is 32k, use io.CopyBuffer instead

		bw.Close() // close and write the tail boundry
	}()
	return bw.FormDataContentType(), pr, nil
}

func doUpload(addr, filepath string) error {
	// create http request body
	conType, reader, err := createReqBody(filepath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/upload", addr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}

	// add headers
	req.Header.Add("Content-Type", conType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request send err", err.Error())
		return err
	}
	resp.Body.Close()
	return nil
}
