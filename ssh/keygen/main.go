package main

import (
	"fmt"
	"os"

	"github.com/shash-786/Golang/ssh"
)

func main() {
	var (
		pub, priv []byte
		err       error
	)

	if priv, pub, err = ssh.GenerateKeys(); err != nil {
		fmt.Printf("Error in ssh GenerateKeys: %v", err)
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
}
