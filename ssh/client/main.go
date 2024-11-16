package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	var (
		priv, pub, output        []byte
		err                      error
		parsed_public_server_key ssh.PublicKey
		parsed_private_key       ssh.Signer
	)

	if priv, err = os.ReadFile("./key/myKey.pem"); err != nil {
		fmt.Printf("ReadFile error for mykey: %v", err)
		os.Exit(1)
	}

	if pub, err = os.ReadFile("./key/server.pub"); err != nil {
		fmt.Printf("ReadFile error for server: %v", err)
		os.Exit(1)
	}

	if parsed_public_server_key, _, _, _, err = ssh.ParseAuthorizedKey(pub); err != nil {
		fmt.Printf("ParseAuthorizedKey error: %v", err)
		os.Exit(1)
	}

	if parsed_private_key, err = ssh.ParsePrivateKey(priv); err != nil {
		fmt.Printf("ParsePrivateKey error: %v", err)
		os.Exit(1)
	}

	config := &ssh.ClientConfig{
		User: "Shashank",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(parsed_private_key),
		},
		HostKeyCallback: ssh.FixedHostKey(parsed_public_server_key),
	}

	client, err := ssh.Dial("tcp", "localhost:2022", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	output, err = session.Output("whoami")
	if err != nil {
		fmt.Printf("session.Output error: %v", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
