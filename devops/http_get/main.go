package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args[:]
	if len(args) < 2 {
		fmt.Printf("No Arguments to get")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Println("Cannot Parse Url")
		os.Exit(1)
	}

	var response *http.Response
	var err error

	if response, err = http.Get(args[1]); err != nil {
		fmt.Printf("Error in get: %v", err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println("Get Not success")
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Code: %d\nBody: %s\n", response.StatusCode, string(body))
}
