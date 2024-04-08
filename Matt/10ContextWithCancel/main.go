package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
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
	"https://wsj.com",
}

func get(ctx context.Context, website string, results chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, website, nil)
	client := http.DefaultClient

	ticker := time.NewTicker(2 * time.Second).C
	response, err := client.Do(req)

	var r result
	if err != nil {
		r = result{website, err, 0}
	} else {
		r = result{website, nil, time.Since(start).Round(time.Millisecond)}
		defer response.Body.Close()
		// You were closing the Response Body in the case where there was response dummy???
	}

	for {
		select {
		case results <- r:
			return
		case <-ticker:
			fmt.Printf("Tick --> %-20v %v\n", time.Now().Format("02-01-2006 15:04:05"), r)
		}
	}
}

func first(ctx context.Context) (*result, error) {
	results := make(chan result)
	// TO AVOID the blocking of goroutines you can makethe channel buffered
	// results := make(chan result, 4)

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	for _, website := range websites {
		go get(ctx, website, results)
	}

	select {
	case r := <-results:
		return &r, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func main() {
	r, _ := first(context.Background())

	// for range websites {
	if r.err != nil {
		fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.err)
	} else {
		fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.latency)
	}
	// }

	time.Sleep(6 * time.Second)
	fmt.Printf("Number of Running Goroutines %d", runtime.NumGoroutine())
}
