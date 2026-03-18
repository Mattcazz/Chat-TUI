package pkg

import (
	"time"
)

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
	PublicKey string `json:"public_key"`
	Nickname  string `json:"nickname"`
}

type PatchContactRequest struct {
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
}

type ContactDetails struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}

type InboxConversationResponse struct {
	UserName  string    `json:"username"`
	ID        int64     `json:"id"`
	LastMsg   string    `json:"last_message"`
	LastMsgAt time.Time `json:"last_message_at"`
}

type UserResponse struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

type InboxResponse struct {
	User          *UserResponse               `json:"user"`
	Conversations []InboxConversationResponse `json:"conversations"`
}

type SendMsgRequest struct {
	Content string `json:"content"`
}

type CreateConversationDmRequest struct {
	ParticipantID int64 `json:"participant_id"`
}

type MsgResponse struct {
	UserName  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationResponse struct {
	ID                       int64         `json:"id"`
	OtherParticipantNickname string        `json:"other_participant_nickname"`
	Messages                 []MsgResponse `json:"messages"`
}

// Files types

type InitFileUploadRequest struct {
	ToUserID    int64  `json:"to_user_id"`
	FileName    string `json:"file_name"`
	TotalSize   int64  `json:"total_size"`
	TotalChunks int64  `json:"total_chunks"`
}

type UploadFileChunkRequest struct{}
