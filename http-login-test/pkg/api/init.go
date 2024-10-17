package api

import "net/http"

type API interface {
	DoGetrequest(requestURL string) (Response, error)
}

// type ClientIface interface{}

type Options struct {
	LoginURL, Password string
}

type API_instance struct {
	Option Options
	Client *http.Client
}

func New(o Options) API {
	client := http.Client{}
	client.Transport = &myjwt_transport{
		password:  o.Password,
		loginURL:  o.LoginURL,
		transport: http.DefaultTransport,
	}

	api := &API_instance{
		Option: o,
		Client: &client,
	}

	return api
}
