package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occur struct {
	Freq map[string]int `json:"freq"`
}

type Response interface {
	GetResponse() string
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("Words -> %s\n", strings.Join(w.Words, ","))
}

func (o Occur) GetResponse() string {
	map_slice := make([]string, 0)
	for key, value := range o.Freq {
		map_slice = append(map_slice, fmt.Sprintf("%s:%d", key, value))
	}
	return fmt.Sprintf("Occurences -> %s\n", strings.Join(map_slice, ","))
}

func (a *API_instance) DoGetrequest(requestURL string) (Response, error) {
	var (
		err      error
		response *http.Response
		body     []byte
		page     Page
	)

	if response, err = a.Client.Get(requestURL); err != nil {
		return nil, fmt.Errorf("./usage http.get: %v", err)
	}
	defer response.Body.Close()

	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("./usage io.Readall: %v", err)
	}

	if err = json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("page name unmarshal error: %v", err)
	}

	switch page.Name {
	case "words":
		var words Words
		if err = json.Unmarshal(body, &words); err != nil {
			return nil, fmt.Errorf("words unmrashal error: %v", err)
		}
		return words, nil

	case "occurence":
		var occur Occur
		if err = json.Unmarshal(body, &occur); err != nil {
			return nil, fmt.Errorf("occur unmarshal error: %v", err)
		}
		return occur, nil

	default:
		return nil, nil
	}
}
