package api

import (
	"fmt"
	"net/http"
)

type MyJWTTransport struct {
	token     string
	transport http.RoundTripper
	loginURL  string
	password  string
}

func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.password != "" {
		var (
			token string
			err   error
		)
		token, err = DoLoginRequest(http.DefaultClient, m.loginURL, m.password)
		if err != nil {
			return nil, err
		}
		m.token = token
		fmt.Printf("token --> %s\n", token)
	}
	if m.token != "" {
		req.Header.Add("Authorization", "Bearer "+m.token)
	}
	return m.transport.RoundTrip(req)
}
