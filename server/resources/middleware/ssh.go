package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// isValidSshSignature verifies that the provided signature signs the nonce
// using the given public key.
//
// pubKeyStr:  The authorized key string (e.g., "ssh-ed25519 AAAA...")
// nonce:      The random string we sent the user (the data that was signed)
// sigStr:     The base64 encoded signature blob received from the client
func IsValidSshSignature(pubKeyStr, nonce, sigStr string) error {

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKeyStr))
	if err != nil {
		return fmt.Errorf("invalid public key format: %w", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sigStr)
	if err != nil {
		return fmt.Errorf("signature is not valid base64: %w", err)
	}

	signature := &ssh.Signature{}
	if err := ssh.Unmarshal(sigBytes, signature); err != nil {
		return fmt.Errorf("failed to unmarshal ssh signature: %w", err)
	}

	if signature.Format != pubKey.Type() {
		return fmt.Errorf("signature type %s does not match public key type %s", signature.Format, pubKey.Type())
	}

	err = pubKey.Verify([]byte(nonce), signature)
	if err != nil {
		return errors.New("signature verification failed")
	}

	return nil
}
