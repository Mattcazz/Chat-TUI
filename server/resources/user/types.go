package user

import (
	"context"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type ContactRepository interface {
	GetContactsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error)
	GetContactByPair(ctx context.Context, userID1, userID2 int64) (*Contact, error)
	GetContactRequestsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error)
	CreateContact(ctx context.Context, c *Contact) error
	UpdateContact(ctx context.Context, c *Contact) error
	DeleteContact(ctx context.Context, id int64) error
}

type Contact struct {
	ID         int64         `json:"id"`
	FromUserID int64         `json:"from_user_id"`
	ToUserID   int64         `json:"to_user_id"`
	Nickname   string        `json:"nickname"`
	Status     contactStatus `json:"status"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Created_at time.Time     `json:"created_at"`
}

type ChallengeRepository interface {
	CreateChallenge(ctx context.Context, challenge *Challenge) error
	GetChallenge(ctx context.Context, id int64) (*Challenge, error)
	DeleteChallenge(ctx context.Context, id int64, nonce string) error
}

type Challenge struct {
	UserID    int64     `json:"user_id"`
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserDoesNotExistError struct {
	Message string
}

func (e *UserDoesNotExistError) Error() string {
	return e.Message
}

func NewUserDoesNotExistError() *UserDoesNotExistError {
	return &UserDoesNotExistError{
		Message: "User does not exist",
	}
}

func IsUserDoesNotExistError(err error) bool {
	_, ok := err.(*UserDoesNotExistError)
	return ok
}

type contactStatus string

const StatusAccept contactStatus = "accepted"
const StatusPending contactStatus = "pending"
const StatusBlock contactStatus = "blocked"
