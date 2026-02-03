package user

import "context"

// temporary

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error)
	CreateOrLoginUser(ctx context.Context, u *User) error
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type ContactRepository interface {
	GetContactByID(ctx context.Context, id int64) (*Contact, error)
	GetContactsByUserID(ctx context.Context, userID int64) ([]*Contact, error)
	CreateContact(ctx context.Context, c *Contact) error
	UpdateContact(ctx context.Context, c *Contact) error
	DeleteContact(ctx context.Context, id int64) error
}

type Contact struct {
	ID       int64
	UserID   int64
	Username string
}

type ChallengeRepository interface {
	CreateChallenge(ctx context.Context, publicKey, nonce string) error
	GetNonceByPublicKey(ctx context.Context, publickey string) (string, error)
	DeleteChallenge(ctx context.Context, publicKey, nonce string) error
}
