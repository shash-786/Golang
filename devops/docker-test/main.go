package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Booting Server!")

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s", r.URL.String())
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi %s", r.URL.String())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
