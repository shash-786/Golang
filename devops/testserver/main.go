package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WordOut struct {
	doc   string
	input string
	pages []string
}

func (*wh WordOut) page_handler(w http.ResponseWriter, r *http.Request) {
  input := r.URL.Query().Get("input")
  if input != "" {
    wh.pages = append(wh.pages, input)
  }

  Output := &WordOut{
    doc: "words",
    input: input,
    pages:  wh.pages,
  }

  out, err = json.Marshal(Output)
  if err != nil {
    fmt.Println("Error in marhsalling")
  }

  fmt.Fprintf(w, string(out))
}

func main() {
	w := &WordOut{
    pages: make([]string, 0),
	}

  http.HandleFunc("/put", w.page_handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
