package user

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByID(id int) (*User, error) {
	// Implementation goes here
	return nil, nil
}

func (s *Store) CreateUser(c *User) error {
	// Implementation goes here
	return nil
}

func (s *Store) UpdateUser(c *User) error {
	// Implementation goes here
	return nil
}

func (s *Store) DeleteUser(id int) error {
	// Implementation goes here
	return nil
}
