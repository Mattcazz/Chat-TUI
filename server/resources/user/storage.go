package user

import "database/sql"

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetUserByID(id int) (*User, error) {
	// Implementation goes here
	return nil, nil
}

func (s *UserStore) CreateUser(c *User) error {
	// Implementation goes here
	return nil
}

func (s *UserStore) UpdateUser(c *User) error {
	// Implementation goes here
	return nil
}

func (s *UserStore) DeleteUser(id int) error {
	// Implementation goes here
	return nil
}

type ContactStore struct {
	db *sql.DB
}

func NewContactStore(db *sql.DB) *ContactStore {
	return &ContactStore{
		db: db,
	}
}

func (s *ContactStore) GetContactByID(id int) (*Contact, error) {
	// Implementation goes here
	return nil, nil
}

func (s *ContactStore) CreateContact(c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *ContactStore) UpdateContact(c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *ContactStore) DeleteContact(id int) error {
	// Implementation goes here
	return nil
}
