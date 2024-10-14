package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	// "fmt"
	"log"
	"net/http"
	"net/url"
	// "os"
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
	var (
		requestUrl, password, token string
		parsedURL                   *url.URL
		err                         error
		res                         Response
	)

	flag.StringVar(&requestUrl, "url", "", "URL to Access")
	flag.StringVar(&password, "password", "", "enter passsword for login")
	flag.Parse()

	if requestUrl == "" {
		flag.Usage()
		log.Fatalln("No URL or Password Given")
	}

	if parsedURL, err = url.ParseRequestURI(requestUrl); err != nil {
		log.Fatalf("Cannot Parse requestURL: %v", err)
	}

	client := http.Client{}

	if password != "" {
		if token, err = doLoginRequest(client, parsedURL.Scheme+"://"+parsedURL.Host+"/login", password); err != nil {
			log.Fatalf("Cannot Process login request: %v", err)
		}

		client.Transport = &MyJWTTransport{
			token:     token,
			transport: http.DefaultTransport,
		}
		fmt.Printf("Token --> %s\n", token)
	}
	// fmt.Printf("url --> %s\n", parsedURL)
	// os.Exit(0)

	res, err = doRequest(client, parsedURL.String())
	if err != nil {
		log.Fatalf("cannot process request: %v", err)
	}

	if res == nil {
		fmt.Println("not a valid page ")
		return
	}
	fmt.Println(res.GetResponse())
}

func doRequest(client http.Client, requestURL string) (Response, error) {
	var (
		err      error
		response *http.Response
		body     []byte
	)

	if response, err = client.Get(requestURL); err != nil {
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
