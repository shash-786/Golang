package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getString(mp map[string]float64) string {
	str := make([]string, 0)
	for key, val := range mp {
		str = append(str, fmt.Sprintf("%s:%v", key, val))
	}
	return strings.Join(str, " ")
}

func main() {
	var (
		err        error
		parsedURL  *url.URL
		requestURL string
		client     *http.Client
		response   *Response
	)

	flag.StringVar(&requestURL, "u", "", "enter the url to obtain info")
	flag.Parse()

	if requestURL == "" {
		flag.Usage()
		log.Fatalln("no url recieved!")
	}

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		log.Fatalf("cannot parse request url: %v", err)
	}

	client = http.DefaultClient

	if response, err = DoGetRequest(client, parsedURL.String()); err != nil {
		log.Fatalf("Error in DoGetRequest: %v", err)
	}

	fmt.Printf("Page:%s\n", response.Page)
	fmt.Printf("Words:%s\n", strings.Join(response.Words, ","))
	fmt.Printf("Map:%s\n", getString(response.Percentages))
	fmt.Printf("Special:%v\n", response.Special)
	fmt.Printf("ExtraSpecial:%v\n", response.ExtraSpecial)
}
