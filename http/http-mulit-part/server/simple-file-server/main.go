package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const uploadPath = "../upload"

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100)
	mForm := r.MultipartForm

	for k := range mForm.Value {
		v := r.FormValue(k)
		fmt.Printf("v: %s\n", v)
	}

	for k := range mForm.File {
		file, fileheader, err := r.FormFile(k)
		if err != nil {
			fmt.Print("invoke Formfile failed: ", err.Error())
			return
		}
		defer file.Close()

		fmt.Printf("upload file: name: %s, size: %d, header:%#v\n",
			fileheader.Filename, fileheader.Size, fileheader.Header)

		localFileName := uploadPath + "/" + fileheader.Filename
		out, err := os.Create(localFileName)
		if err != nil {
			fmt.Printf("failed to open the upload file %s to writing, error: %s", localFileName, err.Error())
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println("copy file failed:", err.Error())
			return
		}

		fmt.Printf("file %s upload ok\n", localFileName)
	}
}

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.ListenAndServe(":18080", nil)
}
