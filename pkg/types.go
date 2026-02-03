package pkg

type ChallengeRequest struct {
	PublicKey string `json:"public_key"`
}

type ChallengeResponse struct {
	Nonce string `json:"nonce"`
}

type LoginRequest struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"nonce"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
