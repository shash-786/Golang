package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type V struct {
	requestURL string
	statusCode int
}

func get(requestURL string, ch chan *V) {
	var parsedURL *url.URL
	var httpRes *http.Response
	var err error

	response := new(V)
	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		response.requestURL = requestURL
		response.statusCode = 400
		ch <- response
		return
	}

	response.requestURL = parsedURL.String()
	if httpRes, err = http.Get(parsedURL.String()); err != nil {
		// NOTE: DOING httpRes.statusCode when err != nil
		// is a nil dereference
		// response.statusCode = httpRes.StatusCode

		response.statusCode = 400
		ch <- response
		return
	}
	defer httpRes.Body.Close()

	response.statusCode = httpRes.StatusCode
	ch <- response
}

func main() {
	ch := make(chan *V, 4)

	websites := []string{"https://google.com", "https://fb.com", "http://localhost:8080", "https://www.wsj.com"}

	for _, website := range websites {
		go get(website, ch)
	}
	time.Sleep(3 * time.Second)
	for i := 0; i < 4; i++ {
		v := <-ch
		fmt.Println(*v)
	}
}
