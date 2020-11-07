package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpGet(url string, headers [][]string) (response string, code int) {
	req, _ := http.NewRequest("GET", url, nil)
	for i := range headers {
		req.Header.Set(headers[i][0], headers[i][1])
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.StatusCode
}

func httpPostJSON(url string, headers [][]string, json []byte) (response string, code int) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	for i := range headers {
		req.Header.Set(headers[i][0], headers[i][1])
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.StatusCode
}
