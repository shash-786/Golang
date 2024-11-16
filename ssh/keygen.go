package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

const (
	PEM_BLOCK_TYPE = "RSA PRIVATE KEY"
	KEY_SIZE       = 4096
)

func GenerateKeys() ([]byte, []byte, error) {
	var (
		privatekey *rsa.PrivateKey
		pemBlock   *pem.Block
		publickey  ssh.PublicKey
		err        error
	)

	if privatekey, err = rsa.GenerateKey(rand.Reader, KEY_SIZE); err != nil {
		return nil, nil, fmt.Errorf("GenerateKeys error : %v", err)
	}

	pemBlock = &pem.Block{
		Type:  PEM_BLOCK_TYPE,
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}

	if publickey, err = ssh.NewPublicKey(&privatekey.PublicKey); err != nil {
		return nil, nil, fmt.Errorf("ssh.NewPublicKey error: %v", err)
	}

	return pem.EncodeToMemory(pemBlock), ssh.MarshalAuthorizedKey(publickey), nil
}
