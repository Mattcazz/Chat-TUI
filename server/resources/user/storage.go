package user

import (
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	// Implementation goes here
	return nil, nil
}

func (s *UserStore) GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	// Implementation goes here
	return nil, nil
}

func (s *UserStore) CreateOrLoginUser(ctx context.Context, u *User) error {
	var id int
	query := `INSERT INTO users (username, public_key)
          VALUES ($1, $2)
          ON CONFLICT (username, public_key) DO NOTHING
          RETURNING id`

	err := s.db.QueryRowContext(ctx, query, u.Username, u.PublicKey).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUser(ctx context.Context, c *User) error {
	// Implementation goes here
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id int64) error {
	// Implementation goes here
	return nil
}

func scanUser(row *sql.Row) (*User, error) {
	return &User{}, nil
}

type ContactStore struct {
	db *sql.DB
}

func NewContactStore(db *sql.DB) *ContactStore {
	return &ContactStore{
		db: db,
	}
}

func (s *ContactStore) GetContactByID(ctx context.Context, id int64) (*Contact, error) {
	// Implementation goes here
	return nil, nil
}

func (s *ContactStore) CreateContact(ctx context.Context, c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *ContactStore) UpdateContact(ctx context.Context, c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *ContactStore) DeleteContact(ctx context.Context, id int64) error {
	// Implementation goes here
	return nil
}

func (s *ContactStore) GetContactsByUserID(ctx context.Context, userID int64) ([]*Contact, error) {
	// Implementation goes here
	return nil, nil
}

type ChallengeStore struct {
	db *sql.DB
}

func NewChallengeStore(db *sql.DB) *ChallengeStore {
	return &ChallengeStore{
		db: db,
	}
}

func (s *ChallengeStore) CreateChallenge(ctx context.Context, publicKey, nonce string) error {
	return nil
}

func (s *ChallengeStore) GetNonceByPublicKey(ctx context.Context, publickey string) (string, error) {
	return "", nil
}

func (s *ChallengeStore) DeleteChallenge(ctx context.Context, publicKey, nonce string) error {
	return nil
}
