package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func main1() {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("file", filepath.Base("\"/path/to/file\""))
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Open("\"/path/to/file\"")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatal(err)
	}

	formField, err := writer.CreatePart(textproto.MIMEHeader{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = formField.Write([]byte(`"59795289-176e-4d2d-ab40-41cd4ecb05d6"`))

	formField, err = writer.CreateFormField("name")
	if err != nil {
		log.Fatal(err)
	}
	_, err = formField.Write([]byte(`"John"`))

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://admin-api.investment.imaninvest.com/v1/files", form)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}