package api

import (
	"fmt"
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
	return fmt.Sprintf("Words -> %s", strings.Join(w.Words, ","))
}

func (o Occur) GetResponse() string {
	map_slice := make([]string, 0)
	for key, value := range o.Freq {
		map_slice = append(map_slice, fmt.Sprintf("Key:%s Value:%d", key, value))
	}
	return fmt.Sprintf("Occurence\n", strings.Join(map_slice, ""))
}

func (a *API_instance) DoGetrequest(requestURL string) (Response, error) {
}
