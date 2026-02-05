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

func (s *UserStore) CreateUser(ctx context.Context, u *User) (*User, error) {
	query := `INSERT INTO users (username, public_key) VALUES ($1, $2) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, u.Username, u.PublicKey)

	err := row.Scan(&u.ID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	// Implementation goes here
	return nil, nil
}

func (s *UserStore) GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	// Implementation goes here
	query := `SELECT id, username, public_key FROM users WHERE public_key = $1`

	row := s.db.QueryRowContext(ctx, query, publicKey)

	return scanUser(row)
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
	user := new(User)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PublicKey)

	return user, err
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

func (s *ChallengeStore) CreateChallenge(ctx context.Context, challenge *Challenge) error {

	query := `INSERT INTO auth_challenges (user_id, nonce, expires_at) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, query, challenge.UserID, challenge.Nonce, challenge.ExpiresAt)

	return err
}

func (s *ChallengeStore) GetChallenge(ctx context.Context, id int64) (*Challenge, error) {
	query := `SELECT user_id, nonce, expires_at FROM auth_challenges WHERE user_id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	return scanChallenge(row)
}

func (s *ChallengeStore) DeleteChallenge(ctx context.Context, id int64, nonce string) error {
	query := `DELETE FROM auth_challenges WHERE user_id = $1 AND nonce = $2`

	_, err := s.db.ExecContext(ctx, query, id, nonce)

	return err
}

func scanChallenge(row *sql.Row) (*Challenge, error) {
	challenge := new(Challenge)

	err := row.Scan(
		&challenge.UserID,
		&challenge.Nonce,
		&challenge.ExpiresAt)

	return challenge, err
}
