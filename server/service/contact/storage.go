package contact

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetContactByID(id int) (*Contact, error) {
	// Implementation goes here
	return nil, nil
}

func (s *Store) CreateContact(c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *Store) UpdateContact(c *Contact) error {
	// Implementation goes here
	return nil
}

func (s *Store) DeleteContact(id int) error {
	// Implementation goes here
	return nil
}
