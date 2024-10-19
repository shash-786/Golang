package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type WordsOut struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type mockclient struct {
	response *http.Response
}

func (m *mockclient) Get(url string) (resp *http.Response, err error) {
	return m.response, nil
}

func TestDoGetRequest(t *testing.T) {
	var (
		output       []byte
		err          error
		testresponse Response
	)

	wordsout := WordsOut{
		Page:  "words",
		Input: "ao",
		Words: []string{"apples", "oranges"},
	}

	if output, err = json.Marshal(wordsout); err != nil {
		t.Errorf("marshal error: %v", err)
	}

	api := API_instance{
		Client: &mockclient{
			response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(output)),
			},
		},
	}

	if testresponse, err = api.DoGetrequest("http:/localdummyhost:8080/put"); err != nil {
		t.Errorf("Failed DoGetRequest: %v", err)
	}

	if testresponse.GetResponse() != fmt.Sprintf("Words -> %s\n", strings.Join([]string{"apples", "oranges"}, ",")) {
		t.Errorf("Reponse not matching\ngot: %s", testresponse.GetResponse())
	}
}
