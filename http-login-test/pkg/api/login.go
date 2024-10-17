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

func DoLoginRequest(loginURL, password string) (string, error) {
	var (
		err            error
		response       *http.Response
		to_login, body []byte
		resBody        LoginResponse
	)

	loginrequest := LoginRequest{
		Password: password,
	}

	if to_login, err = json.Marshal(loginrequest); err != nil {
		return "", fmt.Errorf("loginRequest Marshal error: %v", err)
	}
	if response, err = http.Post(loginURL, "application/json", bytes.NewBuffer(to_login)); err != nil {
		return "", fmt.Errorf("./usage http.post: %v", err)
	}
	defer response.Body.Close()

	if body, err = io.ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("ioReadall error: %v", err)
	}

	if err = json.Unmarshal(body, &resBody); err != nil {
		return "", fmt.Errorf("loginresponse unmarshal error: %v", err)
	}
	return resBody.Token, nil
}
