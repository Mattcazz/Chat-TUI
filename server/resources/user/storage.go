package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/db"
)

type UserStore struct {
	db db.DBTX
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) WithTx(db *sql.Tx) *UserStore {
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
	query := `SELECT id, username, public_key FROM users WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	return scanUser(row)
}

func (s *UserStore) GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	// Implementation goes here
	query := `SELECT id, username, public_key FROM users WHERE public_key = $1`

	row := s.db.QueryRowContext(ctx, query, publicKey)

	return scanUser(row)
}

func (s *UserStore) UpdateUser(ctx context.Context, c *User) error {
	query := `UPDATE users SET username = $1, public_key = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, c.Username, c.PublicKey, c.ID)

	return err
}

func (s *UserStore) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)

	return err
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
	db db.DBTX
}

func NewContactStore(db *sql.DB) *ContactStore {
	return &ContactStore{
		db: db,
	}
}

func (s *ContactStore) WithTx(db *sql.Tx) *ContactStore {
	return &ContactStore{
		db: db,
	}
}

func (s *ContactStore) GetContactByPair(ctx context.Context, userID1, userID2 int64) (*Contact, error) {
	query := `SELECT id, from_user_id, to_user_id, nickname, created_at FROM contacts WHERE from_user_id = $1 AND to_user_id = $2`

	row := s.db.QueryRowContext(ctx, query, userID1, userID2)

	contact, err := scanContact(row)

	return contact, err
}

func (s *ContactStore) CreateContact(ctx context.Context, c *Contact) error {
	query := `INSERT INTO contacts (from_user_id, to_user_id, nickname, status, created_at) VALUES ($1, $2, $3, $4::contact_status, $5)`

	_, err := s.db.ExecContext(ctx, query, c.FromUserID, c.ToUserID, c.Nickname, c.Status, c.CreatedAt)

	return err
}

func (s *ContactStore) UpdateContact(ctx context.Context, c *Contact) error {
	query := `UPDATE contacts SET nickname = $1, status = $2::contact_status, updated_at = $3 WHERE id = $4`

	log.Printf("Updating contact %d: nickname=%s, status=%s", c.ID, c.Nickname, c.Status)

	_, err := s.db.ExecContext(ctx, query, c.Nickname, c.Status, c.UpdatedAt, c.ID)

	return err
}

func (s *ContactStore) DeleteContact(ctx context.Context, id int64) error {
	query := `DELETE FROM contacts WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)

	return err
}

func (s *ContactStore) GetContactsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	query := `
		SELECT c.id, c.nickname, u.public_key, c.created_at
		FROM contacts c
		JOIN users u ON c.to_user_id = u.id
		WHERE c.from_user_id = $1 AND c.status = 'accepted'		
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contacts []*pkg.ContactDetails

	for rows.Next() {
		contact := new(pkg.ContactDetails)
		if err := rows.Scan(&contact.ID, &contact.Username, &contact.PublicKey, &contact.CreatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (s *ContactStore) GetContactRequestsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	query := `
		SELECT c.id, u.username, u.public_key, c.created_at
		FROM contacts c
		JOIN users u ON c.from_user_id = u.id
		WHERE c.to_user_id = $1 AND c.status = 'pending'		
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contacts []*pkg.ContactDetails

	for rows.Next() {
		contact := new(pkg.ContactDetails)
		if err := rows.Scan(&contact.ID, &contact.Username, &contact.PublicKey, &contact.CreatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func scanContact(row *sql.Row) (*Contact, error) {
	contact := new(Contact)

	err := row.Scan(
		&contact.ID,
		&contact.FromUserID,
		&contact.ToUserID,
		&contact.Nickname,
		&contact.CreatedAt)

	return contact, err
}

type ChallengeStore struct {
	db db.DBTX
}

func NewChallengeStore(db *sql.DB) *ChallengeStore {
	return &ChallengeStore{
		db: db,
	}
}

func (s *ChallengeStore) WithTx(db *sql.Tx) *ChallengeStore {
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

	var challenge Challenge

	err := row.Scan(&challenge.UserID, &challenge.Nonce, &challenge.ExpiresAt)

	return &challenge, err
}

func (s *ChallengeStore) DeleteChallenge(ctx context.Context, id int64, nonce string) error {
	query := `DELETE FROM auth_challenges WHERE user_id = $1 AND nonce = $2`

	_, err := s.db.ExecContext(ctx, query, id, nonce)

	return err
}
