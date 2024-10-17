package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/shash-786/Golang/http-login-test/pkg/api"
)

func main() {
	var (
		err                  error
		password, requestURL string
		parsedURL            *url.URL
		response             api.Response
	)
	flag.StringVar(&password, "p", "", "Password to access the API")
	flag.StringVar(&requestURL, "u", "", "URL to Access")
	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		log.Fatalf("main: ParseRequestURI error %v", err)
	}

	api := api.New(api.Options{
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
		Password: password,
	})

	if response, err = api.DoGetrequest(parsedURL.String()); err != nil {
		log.Fatalf("DoRequest Failed: %v", err)
	}

	if response == nil {
		log.Fatal("couln't fetch response because default page name")
	}
	fmt.Printf("Reponse\n%s", response.GetResponse())
}
