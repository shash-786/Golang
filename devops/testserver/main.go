package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WordOut struct {
  Doc   string `json:"doc"`
  Input string `json:"input"`
  Pages []string `json:"pages"`
}

type WordsHandler struct {
  words []string
  password string
  tokenSecret []byte
}

func (wh *WordsHandler) Page_handler(w http.ResponseWriter , r *http.Request) {
  input := r.URL.Query().Get("input")
  if input != "" {
    wh.words = append(wh.words, input)
  }

  output := WordOut{
    Doc: "words",
    Input: input,
    Pages:  wh.words,
  }

  out, err := json.Marshal(output)
  if err != nil {
    fmt.Println("Error in marhsalling")
    return;
  }

  fmt.Fprint(w, string(out))
} 

func main() {
	w := &WordsHandler{
    words: []string{},
	}

  http.HandleFunc("/put", w.Page_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
