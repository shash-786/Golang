package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var websites = []string{
	"https://google.com",
	"https://fb.com",
	"https://youtube.com",
	"http://localhost:8080/",
}

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(website string, ch chan<- result) {
	time_now := time.Now()
	resp, err := http.Get(website)
	if err != nil {
		ch <- result{website, err, 0}
		return
	}
	defer resp.Body.Close()
	ch <- result{website, nil, time.Since(time_now).Round(time.Millisecond)}
}

func main() {
	ch := make(chan result)
	start := time.Now()
	for _, website := range websites {
		go get(website, ch)
	}

	for range websites {
		r := <-ch

		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}

	fmt.Printf("\n Time To Run the Program %v", time.Since(start).Round(time.Second))
}
