package ssh

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func StartServer(serverPK []byte, userPubK []byte) error {
	var (
		err              error
		authorizedKeyMap map[string]bool
		config           *ssh.ServerConfig
		priv             ssh.Signer
	)

	authorizedKeyMap = make(map[string]bool)
	for len(userPubK) > 0 {
		var (
			pubkey ssh.PublicKey
			rest   []byte
		)
		if pubkey, _, _, rest, err = ssh.ParseAuthorizedKey(userPubK); err != nil {
			return fmt.Errorf("ParseAuthorizedKey error: %v", err)
		}

		authorizedKeyMap[string(pubkey.Marshal())] = true
		userPubK = rest
		/*
		   NOTE: what the userPubK file might look like

		   ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEArfV3J2j1rZNCj+DZVRuK0r+yI8J4E...
		   ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICoaWjSmZhFKECmAZHpdsXkTPxH/...
		   ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAsJ+6Md4cVV9oqfSmK9FSQJME98W...
		*/
	}

	config = &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeyMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	if priv, err = ssh.ParsePrivateKey(serverPK); err != nil {
		return fmt.Errorf("ssh.ParsePrivateKey error: %v", err)
	}
	config.AddHostKey(priv)

	listner, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		return fmt.Errorf("net.Listen error: %v", err)
	}

	for {
		ncon, err := listner.Accept()
		if err != nil {
			fmt.Printf("failed to accept incoming connection: %v", err)
		}

		conn, chans, reqs, err := ssh.NewServerConn(ncon, config)
		if err != nil {
			fmt.Printf("NewServerConn error: %v", err)
		}

		if conn != nil && conn.Permissions != nil {
			log.Printf("logged in with key %s", conn.Permissions.Extensions["pubkey-fp"])
		}

		go ssh.DiscardRequests(reqs)
		go handleConnection(conn, chans)
	}
}

func handleConnection(conn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unkown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Printf("Request Type made by client: %s\n", req.Type)
				switch req.Type {
				case "pty-req":
					createTerminal(conn, channel)
				case "shell":
					req.Reply(true, nil)

				default:
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}

func createTerminal(conn *ssh.ServerConn, channel ssh.Channel) {
	termInstance := term.NewTerminal(channel, "> ")

	go func() {
		defer channel.Close()
		for {
			line, err := termInstance.ReadLine()
			if err != nil {
				fmt.Printf("ReadLinde error: %s", err)
				break
			}

			switch line {
			case "whoami":
				termInstance.Write([]byte(fmt.Sprintf("You are: %s\n", conn.Conn.User())))
			case "":
			default:
				termInstance.Write([]byte("Command not found\n"))
			}
		}
	}()
}
