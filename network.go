package main

import (
	"encoding/json"
	"net/http"
)

var Net network

type network struct{}

func (n network) Call(httpMethod string, url string, output interface{}) error {

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body := resp.Body
	defer body.Close()

	err = json.NewDecoder(body).Decode(&output)
	if err != nil {
		return err
	}
	return nil
}
