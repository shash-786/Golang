package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Page         string             `json:"page"`
	Words        []string           `json:"words"`
	Percentages  map[string]float64 `json:"percentages"`
	Special      []*string          `json:"special"`
	ExtraSpecial []any              `json:"extraSpecial"`
}

func DoGetRequest(client *http.Client, request string) (*Response, error) {
	var (
		err      error
		response *http.Response
		body     []byte
		output   Response
	)

	if response, err = client.Get(request); err != nil {
		return nil, fmt.Errorf("./usage http.get %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Status code not 200 but %d", response.StatusCode)
	}

	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("./usage io.Readall %v", err)
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("not a valid json body")
	}

	if err = json.Unmarshal(body, &output); err != nil {
		return nil, fmt.Errorf("unmarshal error %v", err)
	}

	return &output, nil
}
