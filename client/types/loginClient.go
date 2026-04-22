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

func (c *LoginClient) Login(pk []byte, signature []byte) (pkg.LoginResponse, error) {
	loginRequest := pkg.LoginRequest{
		PublicKey: string(pk),
		Signature: string(signature),
	}

	// req.Header.Add("Content-Type", "application/json")

	logger.Log.Printf("Attempting to log in with:")
	logger.Log.Printf("\tPublic Key: %s", loginRequest.PublicKey)
	logger.Log.Printf("\tSignature: %s", loginRequest.Signature)

	resp, err := c.Client.doRequest("POST", "login", loginRequest, nil)
	if err != nil {
		logger.Log.Panicf("Error while logging in: %s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	logger.Log.Printf("Response Body: %s", body)
	if err != nil {
		logger.Log.Panicf("Error trying to log in: %s", err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var loginResponse pkg.LoginResponse
		json.Unmarshal(body, &loginResponse)

		c.Client.Config.SetJWT(loginResponse.Token)

		return loginResponse, nil
	default:
		return pkg.LoginResponse{}, errors.New("Received an unsupported response status: " + resp.Status);
	}
}

func (c *LoginClient) RequestChallenge(pk []byte) ([]byte, error) {
	loginRequest := pkg.LoginRequest{
		PublicKey: string(pk),
	}

	logger.Log.Printf("Attempting to request a challenge with:")
	logger.Log.Printf("\tPublic Key: %s", loginRequest.PublicKey)

	var challengeResponse pkg.ChallengeResponse
	resp, err := c.Client.doRequest("POST", "login", loginRequest, &challengeResponse)
	if err != nil {
		logger.Log.Panicf("Error requesting challenge: %s", err.Error())
	}

	switch resp.StatusCode {
	case http.StatusAccepted:
		// Nonce coming
		logger.Log.Printf("Challenge nonce received: %s", challengeResponse.Nonce)

		return []byte(challengeResponse.Nonce), nil
	case http.StatusBadRequest:
		// Internal server error (only real case I've found so far is duplicate key)
		// Not sure why this happens, for now return a specific error and log it
		logger.Log.Printf("Got %s error from server", resp.Status)
		return nil, errors.ErrUnsupported
	default:
		return nil, errors.New("Received an unsupported response status: " + resp.Status);
	}

}

func (c *LoginClient) Register(pk []byte, username string) {
	registerRequest := pkg.RegisterRequest{
		PublicKey: string(pk),
		Username: username,
	}

	logger.Log.Printf("Attempting to register with:")
	logger.Log.Printf("\tPublic Key: %s", registerRequest.PublicKey)
	logger.Log.Printf("\tUsername: %s", registerRequest.Username)

	resp, err := c.Client.doRequest("POST", "register", registerRequest, nil)
	if err != nil {
		logger.Log.Panicf("Error trying to register: %s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Panicf("Error parsing body: %s", err.Error())
	}
	logger.Log.Printf("Response Body: %s", body)
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Register ok
	default:
		logger.Log.Panicf("Got unhandled status when trying to register: %s", resp.Status)
	}
}
