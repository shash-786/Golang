package api

import "net/http"

type API interface {
	DoRequest(requestURL string) (Response, error)
}

type Options struct {
	Password string
	LoginURL string
}

type API_instance struct {
	Option *Options
	Client http.Client
}

func New(o *Options) API {
	client := http.Client{}
	client.Transport = &MyJWTTransport{
		transport: http.DefaultTransport,
		loginURL:  o.LoginURL,
		password:  o.Password,
	}

	created_api := API_instance{
		o,
		client,
	}
	return created_api
}
