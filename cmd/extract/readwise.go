package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Highlight struct {
	DocumentID   string `json:"document_id"`
	Text         string `json:"text"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	SourceType   string `json:"source_type"`
	Category     string `json:"category"`
	LocationType string `json:"location_type"`
	Location     int    `json:"location"`
}
type ReadWise struct {
	Highlights []Highlight `json:"highlights"`
}

const url = "https://readwise.io/api/v2/highlights/"

func SendRequest(data []byte, token string) (body []byte, err error) {
	client := http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorizaton", "Token "+token)
	res, err := client.Do(req)

	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("not ok, Status: %d", res.StatusCode)
		return
	}

	return body, nil
}
