package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response interface {
	GetResponse() string
}

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

func (w Words) GetResponse() string {
	return fmt.Sprintf("Words:%s", strings.Join(w.Words, ","))
}

func (o Occur) GetResponse() string {
	str := make([]string, 0)
	for key, val := range o.Freq {
		str = append(str, fmt.Sprintf("Key:%s Val:%d", key, val))
	}
	return fmt.Sprintf("%s", strings.Join(str, ","))
}

func (api API_instance) DoRequest(requestURL string) (Response, error) {
	var (
		err      error
		response *http.Response
		body     []byte
		page     Page
	)

	if response, err = api.Client.Get(requestURL); err != nil {
		return nil, fmt.Errorf("./usage http.get: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error Code not 200 but %d\nbody:%s", response.StatusCode, string(body))
	}

	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	if err = json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("Unmarhsal Error: %v", err)
	}

	switch page.Name {
	case "words":
		var words Words
		if err = json.Unmarshal(body, &words); err != nil {
			return nil, fmt.Errorf("unmarshal error words: %v", err)
		}
		return words, nil

	case "occurence":
		var occur Occur
		if err = json.Unmarshal(body, &occur); err != nil {
			return nil, fmt.Errorf("unmarshal error occur: %v", err)
		}
		return occur, nil

	default:
		return nil, nil
	}
}
