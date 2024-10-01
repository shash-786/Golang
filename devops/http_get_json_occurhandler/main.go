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


type Page struct {
  Name string `json:"page"`
}

type Words struct {
  Input string `json:"input"`
  Words []string `json:"words"`
}

type Occur struct {
  Freq map[string]int `json:"freq"`
}

func main() {
  uri := os.Args[1]

  var err error
  var response *http.Response
  var body []byte

  if _, err = url.ParseRequestURI(uri); err != nil {
    log.Fatalf("cannot parse request url: %v",err)
  }

  if response, err = http.Get(uri); err != nil {
    log.Fatalf("./usage get: %v", err)
  }
  
  if response.StatusCode != 200 {
    log.Fatalf("./usage get status code: %d", response.StatusCode)
  }
  defer response.Body.Close()

  if body, err = io.ReadAll(response.Body); err != nil {
    log.Fatalf("./usage readall: %v", err)
  }

  var page Page
  if err = json.Unmarshal(body, &page); err != nil {
    log.Fatalf("/.usage unmarshal: %v", err)
  }

  switch page.Name {
  case "words":
    var words Words
    if err = json.Unmarshal(body, &words); err != nil {
      log.Fatal(err)
    }
    fmt.Printf("Doc:%s\nPages:%v", words.Input, words.Words)

  case "occurence":
    var occur Occur
    if err = json.Unmarshal(body, &occur); err != nil {
      log.Fatal(err)
    }
    for key, value := range occur.Freq {
      fmt.Printf("Key:%s\tValue:%d\n", key, value)
    }

  default:
    fmt.Println("not a valid page")
  }
}
