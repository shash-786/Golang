package api

import (
	"fmt"
	"log"
	"net/http"
)

type myjwt_transport struct {
	token, loginURL, password string
	transport                 http.RoundTripper
	HTTPclient                ClientIface
}

func (m *myjwt_transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.password != "" {
		var (
			err   error
			token string
		)
		if token, err = DoLoginRequest(m.HTTPclient, m.loginURL, m.password); err != nil {
			log.Fatalf("LoginRequest Failed: %v", err)
		}
		m.token = token
		fmt.Println(token)
	}
	if m.token != "" {
		req.Header.Add("Authorization", "Bearer "+m.token)
	}

	return m.transport.RoundTrip(req)
}
