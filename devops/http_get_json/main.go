package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Words struct {
  Page string `json:"page"`
  Input string `json:"input"`
  Words []string `json:"words"`
}

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

  // fmt.Printf("Code: %d\nBody: %s\n", response.StatusCode, string(body))
  var words Words

  if err = json.Unmarshal(body, &words); err != nil {
    log.Fatal(err)
  }
  
  fmt.Printf("Doc:%s\nPages:%v", words.Page, words.Words)
}
