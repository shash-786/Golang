package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
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

	return nil
}
