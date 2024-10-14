package main

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

func doLoginRequest(client http.Client, requestURL, password string) (string, error) {
	var (
		err            error
		to_login, body []byte
		response       *http.Response
		loginresponse  LoginResponse
	)
	loginrequest := LoginRequest{
		Password: password,
	}

	if to_login, err = json.Marshal(loginrequest); err != nil {
		return "", fmt.Errorf("json Marshal error: %v", err)
	}

	if response, err = client.Post(requestURL, "application/json", bytes.NewBuffer(to_login)); err != nil {
		return "", fmt.Errorf("/usage http.post: %v", err)
	}
	defer response.Body.Close()

	if body, err = io.ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("ReadAll Error: %v", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid Status Code: %d\nResponse Body: %s", response.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return "", fmt.Errorf("Not Valid JSON")
	}

	if err = json.Unmarshal(body, &loginresponse); err != nil {
		return "", fmt.Errorf("Unmarshal Login response error: %v", err)
	}

	return loginresponse.Token, err
}
