package api

import (
	"log"
	"net/http"
)

type myjwt_transport struct {
	token, loginURL, password string
	transport                 http.RoundTripper
}

func (m *myjwt_transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.password != "" {
		var (
			err   error
			token string
		)
		if token, err = DoLoginRequest(m.loginURL, m.password); err != nil {
			log.Fatalf("LoginRequest Failed: %v", err)
		}
		m.token = token
	}
	if m.token != "" {
		req.Header.Add("Authorization", "Bearer "+m.token)
	}

	return m.transport.RoundTrip(req)
}
