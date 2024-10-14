package main

import (
	"flag"
	// "fmt"
	"log"
	"net/url"
	// "os"
)

func main() {
	var (
		requestUrl, password, token string
		parsedURL                   *url.URL
		err                         error
	)

	flag.StringVar(&requestUrl, "url", "", "URL to Access")
	flag.StringVar(&password, "password", "", "enter passsword for login")
	flag.Parse()

	if requestUrl == "" || password == "" {
		log.Fatalln("No URL or Password Given")
	}

	if parsedURL, err = url.ParseRequestURI(requestUrl); err != nil {
		log.Fatalf("Cannot Parse requestURL: %v", err)
	}

	if token, err = doLoginRequest(parsedURL.Scheme+"://"+parsedURL.Host+"/login", password); err != nil {
		log.Fatalf("Cannot Process login request: %v", err)
	}

	// fmt.Printf("Token --> %s\n", token)
	// os.Exit(0)

	// if res, err = doRequest(parsedURL.String()); err != nil {
	// 	log.Fatalf("Cannot Process Request: %v", err)
	// }
}
