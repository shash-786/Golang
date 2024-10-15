package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/shash-786/Golang/http-login/pkg/api"
)

func main() {
	var (
		password, requestURL string
		err                  error
		parsedURL            *url.URL
		response             api.Response
	)
	flag.StringVar(&password, "p", "", "Password to Access API")
	flag.StringVar(&requestURL, "u", "", "Request url")
	flag.Parse()

	if requestURL == "" {
		flag.Usage()
		log.Fatal("No Request Url Given!")
	}

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		log.Fatalf("Error in parsing uri: %v", err)
	}

	new_api := api.New(&api.Options{
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
		Password: password,
	})

	if response, err = new_api.DoRequest(parsedURL.String()); err != nil {
		log.Fatalf("error in processing request: %v", err)
	}

	if response == nil {
		log.Fatal("page switched to default")
	}

	fmt.Printf("Result\n%s", response.GetResponse())
}
