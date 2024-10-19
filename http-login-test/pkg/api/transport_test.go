package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type mockroundtripper struct {
	round_tripper_response *http.Response
}

func (m *mockroundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer xqc" {
		return nil, fmt.Errorf("Token different %s", req.Header.Get("Authorization"))
	}
	return m.round_tripper_response, nil
}

func (m *mockclient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return m.Postresponse, nil
}

func TestRoundTrip(t *testing.T) {
	loginresponse := LoginResponse{
		Token: "xqc",
	}
	response, err := json.Marshal(loginresponse)
	if err != nil {
		t.Errorf("Mashal error: %v", err)
	}

	mocktransport := myjwt_transport{
		transport: &mockroundtripper{
			round_tripper_response: &http.Response{
				StatusCode: 200,
			},
		},

		HTTPclient: &mockclient{
			Postresponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(response)),
			},
		},
		password: "abc",
	}

	request := &http.Request{
		Header: make(http.Header),
	}

	res, err := mocktransport.RoundTrip(request)
	if err != nil {
		t.Errorf("RoundTrip Fail: %v", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("RoundTrip Fail status code not 200 but %d", res.StatusCode)
	}
}
