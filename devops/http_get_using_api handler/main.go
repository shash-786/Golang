package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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
	for k, v := range o.Freq {
		str = append(str, fmt.Sprintf("Key:%s Value:%d\n", k, v))
	}
	return fmt.Sprintf("%s", strings.Join(str, "\n"))
}

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalln("no url given")
	}

	res, err := doRequest(args[1])
	if err != nil {
		log.Fatalf("cannot process request: %v", err)
	}

	if res == nil {
		fmt.Println("not a valid page ")
		return
	}
	fmt.Println(res.GetResponse())
}

func doRequest(requestURL string) (Response, error) {
	var (
		err      error
		response *http.Response
		body     []byte
	)

	if _, err = url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("cannot parse request url: %v", err)
	}

	if response, err = http.Get(requestURL); err != nil {
		return nil, fmt.Errorf("./usage get error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response Code:%d\nError:%v", response.StatusCode, err)
	}

	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("cannot read the response body: %v", err)
	}

	var page Page

	if err = json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("Unmarshal error page: %v", err)
	}

	switch page.Name {
	case "words":
		var words Words
		if err = json.Unmarshal(body, &words); err != nil {
			return nil, fmt.Errorf("Unmarshal error words: %v", err)
		}
		return words, nil

	case "occurence":
		var occur Occur
		if err = json.Unmarshal(body, &occur); err != nil {
			return nil, fmt.Errorf("Unmarshal error occur: %v", err)
		}
		return occur, nil

	default:
		return nil, nil
	}
}
