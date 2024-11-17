package main

import (
	"fmt"
	"os"

	"github.com/shash-786/Golang/ssh"
)

func main() {
	var (
		publicKey, token string
		err              error
	)

	if publicKey, err = generatekeys(); err != nil {
		fmt.Printf("generatekeys error: %v", err)
		os.Exit(1)
	}

	if token, err = generateToken(); err != nil {
		fmt.Printf("generatekeys error: %v", err)
		os.Exit(1)
	}

	if err = launchInstance(); err != nil {
		fmt.Printf("launchInstance error: %v", err)
		os.Exit(1)
	}
}

func generatekeys() (string, error) {
	var priv, pub []byte
	var err error

	if priv, pub, err = ssh.GenerateKeys(); err != nil {
		fmt.Printf("GenerateKeys error: %v ", err)
		os.Exit(1)
	}

	if err = os.WriteFile("./key/myKey.pem", priv, 0600); err != nil {
		fmt.Printf("priv writefile error: %v", err)
		os.Exit(1)
	}

	if err = os.WriteFile("./key/myKey.pub", pub, 0644); err != nil {
		fmt.Printf("pub writefile error: %v", err)
		os.Exit(1)
	}

	return string(pub), err
}

func generateToken() (string, error) {
	return "", nil
}

func launchInstance() error {
	return nil
}
