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

func get(website string) result {
	time_now := time.Now()
	resp, err := http.Get(website)
	if err != nil {
		log.Printf("Website %s not reachable %s", website, err)
		return result{url: website, err: err, latency: 0}
	}
	defer resp.Body.Close()

	return result{url: website, err: nil, latency: time.Since(time_now).Round(time.Millisecond)}
}

func main() {
	var results []result
	start := time.Now()
	for _, website := range websites {
		results = append(results, get(website))
	}

	for _, r := range results {
		fmt.Printf("URL:%s\tLATENCY:%s", r.url, r.latency)
		fmt.Println()
	}

	fmt.Printf("\nTime to run the program : %v\n", time.Since(start).Round(time.Second))
}
