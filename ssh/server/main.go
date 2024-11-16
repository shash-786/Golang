package main

import (
	"fmt"
	"os"

	"github.com/shash-786/Golang/ssh"
)

func main() {
	var (
		priv, pub []byte
		err       error
	)

	if priv, err = os.ReadFile("./key/server.pem"); err != nil {
		fmt.Printf("io.ReadAll error %v", err)
		os.Exit(1)
	}

	if pub, err = os.ReadFile("./key/myKey.pub"); err != nil {
		fmt.Printf("io.ReadAll error %v", err)
		os.Exit(1)
	}

	if err = ssh.StartServer(priv, pub); err != nil {
		fmt.Printf("Start Server error --> %v", err)
		os.Exit(1)
	}
}
