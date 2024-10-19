package api

import (
	"io"
	"net/http"
)

type API interface {
	DoGetrequest(requestURL string) (Response, error)
}

type ClientIface interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Options struct {
	LoginURL, Password string
}

type API_instance struct {
	Option Options
	Client ClientIface
}

func New(o Options) API {
	client := &http.Client{}
	client.Transport = &myjwt_transport{
		password:   o.Password,
		loginURL:   o.LoginURL,
		transport:  http.DefaultTransport,
		HTTPclient: &http.Client{},
	}

	api := &API_instance{
		Option: o,
		Client: client,
	}

	return api
}
