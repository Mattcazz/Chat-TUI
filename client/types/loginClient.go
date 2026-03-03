package types

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/pkg"
)

type LoginClient struct {
	Client BaseClient
}

func (c *LoginClient) Login(pk []byte, signature []byte) error {
	login_req := pkg.LoginRequest{
		PublicKey: string(pk),
		Signature: string(signature),
	}

	// req.Header.Add("Content-Type", "application/json")

	logger.Log.Printf("Attempting to log in with:")
	logger.Log.Printf("\tPublic Key: %s", login_req.PublicKey)
	logger.Log.Printf("\tSignature: %s", login_req.Signature)

	resp, err := c.Client.doRequest("POST", "login", login_req, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	logger.Log.Printf("Response Body: %s", body)
	if err != nil {
		log.Panic("Trying to log in: " + err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var login_resp pkg.LoginResponse
		json.Unmarshal(body, &login_resp)

		c.Client.Config.SetJWT(login_resp.Token)

		return nil
	default:
		return errors.New("Received an unsupported response status: " + resp.Status);
	}
}

func (c *LoginClient) RequestChallenge(pk []byte) ([]byte, error) {
	login_req := pkg.LoginRequest{
		PublicKey: string(pk),
	}

	logger.Log.Printf("Attempting to request a challenge with:")
	logger.Log.Printf("\tPublic Key: %s", login_req.PublicKey)

	resp, err := c.Client.doRequest("POST", "login", login_req, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	logger.Log.Printf("Response Body: %s", body)
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		// Nonce coming
		var challenge_resp pkg.ChallengeResponse
		json.Unmarshal(body, &challenge_resp)
		logger.Log.Printf("Challenge nonce received: %s", challenge_resp.Nonce)

		return []byte(challenge_resp.Nonce), nil
	default:
		return nil, errors.New("Received an unsupported response status: " + resp.Status);
	}

}

func (c *LoginClient) Register(pk []byte, username string) {
	register_req := pkg.RegisterRequest{
		PublicKey: string(pk),
		Username: username,
	}

	logger.Log.Printf("Attempting to register with:")
	logger.Log.Printf("\tPublic Key: %s", register_req.PublicKey)
	logger.Log.Printf("\tUsername: %s", register_req.Username)

	resp, err := c.Client.doRequest("POST", "register", register_req, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	logger.Log.Printf("Response Body: %s", body)
	if err != nil {
		// TODO
		log.Panic(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Register ok
	default:
		log.Panic("Could not register")
	}
}
