package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func DoLoginRequest(client *http.Client, requestURL, password string) (string, error) {
	var (
		err            error
		response       *http.Response
		to_login, body []byte
	)

	loginrequest := LoginRequest{
		Password: password,
	}

	if to_login, err = json.Marshal(loginrequest); err != nil {
		return "", fmt.Errorf("Marshalling error: %v", err)
	}
	if response, err = client.Post(requestURL, "application/json", bytes.NewBuffer(to_login)); err != nil {
		return "", fmt.Errorf("./usage http.post: %v", err)
	}
	defer response.Body.Close()
	if body, err = io.ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("ReadAll error Login: %v", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid Status Code: %d\nResponse Body: %s", response.StatusCode, string(body))
	}
	if !json.Valid(body) {
		return "", fmt.Errorf("Not Valid JSON")
	}

	var loginresponse LoginResponse
	if err = json.Unmarshal(body, &loginresponse); err != nil {
		return "", fmt.Errorf("Login response unmarshall error: %v", err)
	}
	return loginresponse.Token, nil
}
