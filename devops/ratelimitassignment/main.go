package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	var (
		err        error
		args       []string
		parsedURL  *url.URL
		requestURL string
	)

	args = os.Args
	if len(args) < 2 {
		requestURL = "http://localhost:8080/ratelimit"
	} else {
		requestURL = args[1]
	}
	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		log.Fatalf("cannot parse request url: %v", err)
	}

	hit := time.NewTicker(200 * time.Millisecond)
	stop := time.After(time.Second * 5)

	for {
		select {
		case <-hit.C:
			go hit_rate_limit(parsedURL)
		case <-stop:
			return
		}
	}
}

func hit_rate_limit(parsedURL *url.URL) {
	response, err := http.Get(parsedURL.String())
	if err != nil {
		fmt.Printf("./usage http.get: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("./usage io.readall: %v", err)
	}
	fmt.Println(response.Status)
	fmt.Print(string(body))
}
