package login

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"golang.org/x/crypto/ssh"
)

func getSSHKeys() ([]byte, []byte) {
	ssh_key_path := os.Getenv("HOME") + "/.ssh/" + config.Configuration.SSH_key_name

	pkPath := ssh_key_path + ".pub"
	pkBytes, err := os.ReadFile(pkPath)
	if err != nil {
		panic("Could not read public key: " + err.Error())
	}
	pkBytes = bytes.TrimSpace(pkBytes)

	skPath := ssh_key_path
	skBytes, err := os.ReadFile(skPath)
	if err != nil {
		panic("Could not read private key: " + err.Error())
	}
	skBytes = bytes.TrimSpace(skBytes)

	return pkBytes, skBytes
}

func createSignature(nonce string, sk []byte, passphrase []byte) ([]byte, error) {
	// Parse Private Key
	signer, err := ssh.ParsePrivateKey(sk)

	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			if passphrase == nil {
				return nil, err
			}
			signer, err = ssh.ParsePrivateKeyWithPassphrase(sk, passphrase)
			if err != nil {
				panic(err) // probably wrong password
			}
		} else {
			panic(err)
		}
	}

	// Sign the Nonce
	sig, err := signer.Sign(rand.Reader, []byte(nonce))
	if err != nil {
		panic(err)
	}

	// Encode to Base64 (This is what you paste into Curl)
	sigBytes := ssh.Marshal(sig) // Important: Marshal to wire format first!
	b64Sig := base64.StdEncoding.EncodeToString(sigBytes)

	return []byte(b64Sig), nil
}

