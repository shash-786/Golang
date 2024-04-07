/*
Context

> Includes Deadlines, Cancellation Signals
> Typically used across API Boundaries and Between Processes

Incoming server should create a context and outgoing calls to server should
accept context

Contexts form an immutable tree structure

Contexts are of three types

1) Context.Background
	When a new request comes in, you start with a background
	context (context.Background()), which serves as the parent
	context for the request.

	The background context is used as the initial context for
	handling the request.

2) Context With Value
	Send Request scoped data across Different Parts of your program

	example:

	func auth(w , r*) {
		usrname = r.URL.Query().Get("Username")
		passwrd = r.URL.Query().Get("Password")
		ctx := context.WithValue(r.Context(), authkey("username"), usrname)
		ctx = context.WithValue(ctx, authkey("password"), passwrd)
		nextHandler(w, r.WithContext(ctx))
	}

	func nextHandler(w http.ResponseWriter, r *http.Request) {
		Retrieve the authentication information from the context
		username := r.Context().Value(authKey("username")).(string)
		password := r.Context().Value(authKey("password")).(string)

		Perform some operation with the authentication information
		fmt.Fprintf(w, "Authenticated user: %s\n", username)
	}

	3) Context with Timeouts

	This function creates a new context with the given timeout duration.

	It's useful for specifying how long an operation should be allowed to
	take before it's canceled.

	The timeout duration determines when the context and
	any operations associated with it should be canceled

*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var websites = []string{
	"https://google.com",
	"https://amazon.com",
	"https://fb.com",
	"http://localhost:8080/",
}

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, website string, channel chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, website, nil)
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		channel <- result{website, err, 0}
	} else {
		channel <- result{website, nil, time.Since(start).Round(time.Millisecond)}
		response.Body.Close()
	}
}

func main() {
	channel := make(chan result)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// Go get websites
	for _, website := range websites {
		go get(ctx, website, channel)
	}

	// reading from the channel
	for range websites {
		r := <-channel

		if r.err != nil {
			fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.err)
		} else {
			fmt.Printf("%-20v %-20s %s\n", time.Now().Format("02-01-2006 15:04:05"), r.url, r.latency)
		}
	}

}
