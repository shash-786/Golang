package main

import (
	"log"
	"net/http"
	"time"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(7 * time.Second)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", MyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
