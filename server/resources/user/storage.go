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
	log.Printf("UserStore.CreateUser: Inserting new user with username %s", u.Username)
	query := `INSERT INTO users (username, public_key) VALUES ($1, $2) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, u.Username, u.PublicKey)

	err := row.Scan(&u.ID)
	if err != nil {
		log.Printf("UserStore.CreateUser: Failed to create user with username %s: %v", u.Username, err)
		return nil, err
	}

	log.Printf("UserStore.CreateUser: Successfully created user with ID %d, username %s", u.ID, u.Username)
	return u, nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	log.Printf("UserStore.GetUserByID: Retrieving user with ID %d", id)
	query := `SELECT id, username, public_key FROM users WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	user, err := scanUser(row)
	if err != nil {
		log.Printf("UserStore.GetUserByID: Failed to retrieve user with ID %d: %v", id, err)
		return nil, err
	}

	log.Printf("UserStore.GetUserByID: Successfully retrieved user with ID %d, username %s", user.ID, user.Username)
	return user, nil
}

func (s *UserStore) GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	log.Printf("UserStore.GetUserByPublicKey: Retrieving user by public key")
	query := `SELECT id, username, public_key FROM users WHERE public_key = $1`

	row := s.db.QueryRowContext(ctx, query, publicKey)

	user, err := scanUser(row)
	if err != nil {
		log.Printf("UserStore.GetUserByPublicKey: Failed to retrieve user by public key: %v", err)
		return nil, err
	}

	log.Printf("UserStore.GetUserByPublicKey: Successfully retrieved user with ID %d, username %s", user.ID, user.Username)
	return user, nil
}

func (s *UserStore) UpdateUser(ctx context.Context, c *User) error {
	log.Printf("UserStore.UpdateUser: Updating user with ID %d, username %s", c.ID, c.Username)
	query := `UPDATE users SET username = $1, public_key = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, c.Username, c.PublicKey, c.ID)
	if err != nil {
		log.Printf("UserStore.UpdateUser: Failed to update user with ID %d: %v", c.ID, err)
		return err
	}

	log.Printf("UserStore.UpdateUser: Successfully updated user with ID %d", c.ID)
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id int64) error {
	log.Printf("UserStore.DeleteUser: Deleting user with ID %d", id)
	query := `DELETE FROM users WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("UserStore.DeleteUser: Failed to delete user with ID %d: %v", id, err)
		return err
	}

	log.Printf("UserStore.DeleteUser: Successfully deleted user with ID %d", id)
	return nil
}

func (s *UserStore) GetUserConversations(ctx context.Context, userID int64) ([]*pkg.InboxConversationResponse, error) {
	var conversations []*pkg.InboxConversationResponse

	return conversations, nil
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
	log.Printf("ContactStore.GetContactByPair: Retrieving contact between user ID %d and user ID %d", userID1, userID2)
	query := `SELECT id, from_user_id, to_user_id, nickname, created_at FROM contacts WHERE from_user_id = $1 AND to_user_id = $2`

	row := s.db.QueryRowContext(ctx, query, userID1, userID2)

	contact, err := scanContact(row)
	if err != nil {
		log.Printf("ContactStore.GetContactByPair: Failed to retrieve contact between user ID %d and user ID %d: %v", userID1, userID2, err)
		return nil, err
	}

	log.Printf("ContactStore.GetContactByPair: Successfully retrieved contact with ID %d", contact.ID)
	return contact, nil
}

func (s *ContactStore) CreateContact(ctx context.Context, c *Contact) error {
	log.Printf("ContactStore.CreateContact: Creating contact from user ID %d to user ID %d with nickname %s", c.FromUserID, c.ToUserID, c.Nickname)
	query := `INSERT INTO contacts (from_user_id, to_user_id, nickname, status, created_at) VALUES ($1, $2, $3, $4::contact_status, $5)`

	_, err := s.db.ExecContext(ctx, query, c.FromUserID, c.ToUserID, c.Nickname, c.Status, c.CreatedAt)
	if err != nil {
		log.Printf("ContactStore.CreateContact: Failed to create contact from user ID %d to user ID %d: %v", c.FromUserID, c.ToUserID, err)
		return err
	}

	log.Printf("ContactStore.CreateContact: Successfully created contact from user ID %d to user ID %d", c.FromUserID, c.ToUserID)
	return nil
}

func (s *ContactStore) UpdateContact(ctx context.Context, c *Contact) error {
	log.Printf("ContactStore.UpdateContact: Updating contact with ID %d: nickname=%s, status=%s", c.ID, c.Nickname, c.Status)
	query := `UPDATE contacts SET nickname = $1, status = $2::contact_status, updated_at = $3 WHERE id = $4`

	_, err := s.db.ExecContext(ctx, query, c.Nickname, c.Status, c.UpdatedAt, c.ID)
	if err != nil {
		log.Printf("ContactStore.UpdateContact: Failed to update contact with ID %d: %v", c.ID, err)
		return err
	}

	log.Printf("ContactStore.UpdateContact: Successfully updated contact with ID %d", c.ID)
	return nil
}

func (s *ContactStore) DeleteContact(ctx context.Context, id int64) error {
	log.Printf("ContactStore.DeleteContact: Deleting contact with ID %d", id)
	query := `DELETE FROM contacts WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("ContactStore.DeleteContact: Failed to delete contact with ID %d: %v", id, err)
		return err
	}

	log.Printf("ContactStore.DeleteContact: Successfully deleted contact with ID %d", id)
	return nil
}

func (s *ContactStore) GetContactsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	log.Printf("ContactStore.GetContactsByUserID: Fetching accepted contacts for user ID %d", userID)
	query := `
		SELECT c.id, c.nickname, u.public_key, c.created_at
		FROM contacts c
		JOIN users u ON c.to_user_id = u.id
		WHERE c.from_user_id = $1 AND c.status = 'accepted'
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("ContactStore.GetContactsByUserID: Failed to query contacts for user ID %d: %v", userID, err)
		return nil, err
	}

	defer rows.Close()

	var contacts []*pkg.ContactDetails

	for rows.Next() {
		contact := new(pkg.ContactDetails)
		if err := rows.Scan(&contact.ID, &contact.Username, &contact.PublicKey, &contact.CreatedAt); err != nil {
			log.Printf("ContactStore.GetContactsByUserID: Failed to scan contact row for user ID %d: %v", userID, err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ContactStore.GetContactsByUserID: Error iterating contact rows for user ID %d: %v", userID, err)
		return nil, err
	}

	log.Printf("ContactStore.GetContactsByUserID: Successfully retrieved %d accepted contacts for user ID %d", len(contacts), userID)
	return contacts, nil
}

func (s *ContactStore) GetContactRequestsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	log.Printf("ContactStore.GetContactRequestsByUserID: Fetching pending contact requests for user ID %d", userID)
	query := `
		SELECT c.id, u.username, u.public_key, c.created_at
		FROM contacts c
		JOIN users u ON c.from_user_id = u.id
		WHERE c.to_user_id = $1 AND c.status = 'pending'
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("ContactStore.GetContactRequestsByUserID: Failed to query contact requests for user ID %d: %v", userID, err)
		return nil, err
	}

	defer rows.Close()

	var contacts []*pkg.ContactDetails

	for rows.Next() {
		contact := new(pkg.ContactDetails)
		if err := rows.Scan(&contact.ID, &contact.Username, &contact.PublicKey, &contact.CreatedAt); err != nil {
			log.Printf("ContactStore.GetContactRequestsByUserID: Failed to scan contact request row for user ID %d: %v", userID, err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ContactStore.GetContactRequestsByUserID: Error iterating contact request rows for user ID %d: %v", userID, err)
		return nil, err
	}

	log.Printf("ContactStore.GetContactRequestsByUserID: Successfully retrieved %d pending contact requests for user ID %d", len(contacts), userID)
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
	log.Printf("ChallengeStore.CreateChallenge: Creating challenge for user ID %d with nonce expiring at %v", challenge.UserID, challenge.ExpiresAt)
	query := `INSERT INTO auth_challenges (user_id, nonce, expires_at) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, query, challenge.UserID, challenge.Nonce, challenge.ExpiresAt)
	if err != nil {
		log.Printf("ChallengeStore.CreateChallenge: Failed to create challenge for user ID %d: %v", challenge.UserID, err)
		return err
	}

	log.Printf("ChallengeStore.CreateChallenge: Successfully created challenge for user ID %d", challenge.UserID)
	return nil
}

func (s *ChallengeStore) GetChallenge(ctx context.Context, id int64) (*Challenge, error) {
	log.Printf("ChallengeStore.GetChallenge: Retrieving challenge for user ID %d", id)
	query := `SELECT user_id, nonce, expires_at FROM auth_challenges WHERE user_id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	var challenge Challenge

	err := row.Scan(&challenge.UserID, &challenge.Nonce, &challenge.ExpiresAt)
	if err != nil {
		log.Printf("ChallengeStore.GetChallenge: Failed to retrieve challenge for user ID %d: %v", id, err)
		return nil, err
	}

	log.Printf("ChallengeStore.GetChallenge: Successfully retrieved challenge for user ID %d", id)
	return &challenge, nil
}

func (s *ChallengeStore) DeleteChallenge(ctx context.Context, id int64, nonce string) error {
	log.Printf("ChallengeStore.DeleteChallenge: Deleting challenge for user ID %d with nonce", id)
	query := `DELETE FROM auth_challenges WHERE user_id = $1 AND nonce = $2`

	_, err := s.db.ExecContext(ctx, query, id, nonce)
	if err != nil {
		log.Printf("ChallengeStore.DeleteChallenge: Failed to delete challenge for user ID %d: %v", id, err)
		return err
	}

	log.Printf("ChallengeStore.DeleteChallenge: Successfully deleted challenge for user ID %d", id)
	return nil
}
