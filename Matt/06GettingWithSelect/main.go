package main

import (
	"fmt"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

var websites = []string{
	"https://google.com",
	"https://fb.com",
	"https://amazon.com",
	"http://localhost:8080/", // NOTE: This is defined in the HTTP Server file and it will basically
	// sleep for 6 seconds before returning Status Code
	"https://lazyvim.org",
}

func get(website string, channel chan<- result) {
	start := time.Now()
	res, err := http.Get(website)
	if err != nil {
		channel <- result{website, err, 0}
		return
	}
	defer res.Body.Close()
	channel <- result{website, nil, time.Since(start).Round(time.Millisecond)}
}

func main() {
	ch := make(chan result)
	time_stop_notify_channel := time.After(4 * time.Second)

	for _, website := range websites {
		go get(website, ch)
	}

	for range websites {
		select {
		case r := <-ch:
			if r.err != nil {
				fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.err)
			} else {
				fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.latency)
			}
		case <-time_stop_notify_channel:
			fmt.Printf("%-20v Connection Timed Out!\n", time.Now().Format("02-01-2006 15:04:05"))
		}
	}
}
