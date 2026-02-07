package pkg

type ChallengeRequest struct {
	PublicKey string `json:"public_key"`
}

type ChallengeResponse struct {
	Nonce string `json:"nonce"`
}

type LoginRequest struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
}

type PatchUserRequest struct {
	Username string `json:"username"`
}

type PostContactRequest struct {
	ID        int64  `json:"id"`
	PublicKey string `json:"public_key"`
	Nickname  string `json:"nickname"`
}

type PatchContactRequest struct {
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
}
