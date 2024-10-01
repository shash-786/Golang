package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WordsOut struct {
  Page string `json:"page"`
  Input string `json:"input"`
  Words []string `json:"words"`
}

type OccurenceOut struct {
  Page string `json:"page"`
  Freq map[string]int `json:"freq"`
}

type database struct {
  words []string
  // password string
  // tokenSecret []byte
}

func (db *database) insert_handler(w http.ResponseWriter , r *http.Request) {
  input := r.URL.Query().Get("input")
  if input != "" {
    db.words = append(db.words, input)
  }

  output := WordsOut{
    Page: "words",
    Input: input,
    Words:  db.words,
  }

  out, err := json.Marshal(output)
  if err != nil {
    fmt.Println("Error in marhsalling")
    return;
  }

  fmt.Fprint(w, string(out))
}

func (db *database) occurence_handler(w http.ResponseWriter, r *http.Request) {
  mp := make(map[string]int)
  for _, v := range db.words {
    if _, ok := mp[v]; !ok {
      mp[v] = 1;
    } else {
      mp[v]++
    }
  }

  occurenceout := OccurenceOut{
    Page: "occurence",
    Freq: mp,
  }
 
  out, err := json.Marshal(occurenceout)
  if err != nil {
    log.Print(err)
  }
  fmt.Fprint(w, string(out))
}

func main() {
	db := &database{
    words: []string{},
	}

  http.HandleFunc("/put", db.insert_handler)
  http.HandleFunc("/occur", db.occurence_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
