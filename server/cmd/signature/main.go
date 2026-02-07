package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Usage: go run cmd/tools/signer.go <NONCE_STRING>
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the nonce as an argument")
		return
	}
	nonce := os.Args[1]

	// 1. Load your private key (Change path if needed)
	keyPath := os.Getenv("HOME") + "/.ssh/id_ed25519"
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic("Could not read private key: " + err.Error())
	}

	// 2. Parse Private Key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		panic(err)
	}

	// 3. Sign the Nonce
	sig, err := signer.Sign(rand.Reader, []byte(nonce))
	if err != nil {
		panic(err)
	}

	// 4. Encode to Base64 (This is what you paste into Curl)
	sigBytes := ssh.Marshal(sig) // Important: Marshal to wire format first!
	b64Sig := base64.StdEncoding.EncodeToString(sigBytes)

	// 5. Get Public Key String (for the request)
	pubKey := ssh.MarshalAuthorizedKey(signer.PublicKey())

	fmt.Println("\n--- COPY THESE FOR CURL ---")
	fmt.Printf("Public Key: %s", pubKey) // Contains newline
	fmt.Printf("Signature:  %s\n", b64Sig)
}
