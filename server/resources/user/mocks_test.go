package user

import (
	"context"
	"database/sql"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

// ---- Mock UserRepository ----

type mockUserRepo struct {
	createUserFn         func(ctx context.Context, u *User) (*User, error)
	getUserByIDFn        func(ctx context.Context, id int64) (*User, error)
	getUserByPublicKeyFn func(ctx context.Context, publicKey string) (*User, error)
	updateUserFn         func(ctx context.Context, u *User) error
	deleteUserFn         func(ctx context.Context, id int64) error
}

func (m *mockUserRepo) WithTx(_ *sql.Tx) *UserStore { return nil }

func (m *mockUserRepo) CreateUser(ctx context.Context, u *User) (*User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, u)
	}
	u.ID = 1
	return u, nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id int64) (*User, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, id)
	}
	return &User{ID: id, Username: "alice", PublicKey: "pk-alice"}, nil
}

func (m *mockUserRepo) GetUserByPublicKey(ctx context.Context, publicKey string) (*User, error) {
	if m.getUserByPublicKeyFn != nil {
		return m.getUserByPublicKeyFn(ctx, publicKey)
	}
	return &User{ID: 1, Username: "alice", PublicKey: publicKey}, nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, u *User) error {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, u)
	}
	return nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id int64) error {
	if m.deleteUserFn != nil {
		return m.deleteUserFn(ctx, id)
	}
	return nil
}

// ---- Mock ContactRepository ----

type mockContactRepo struct {
	getContactsByUserIDFn        func(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error)
	getContactRequestsByUserIDFn func(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error)
	getContactByPairFn           func(ctx context.Context, userID1, userID2 int64) (*Contact, error)
	createContactFn              func(ctx context.Context, c *Contact) error
	updateContactFn              func(ctx context.Context, c *Contact) error
	deleteContactFn              func(ctx context.Context, id int64) error
}

func (m *mockContactRepo) WithTx(_ *sql.Tx) *ContactStore { return nil }

func (m *mockContactRepo) GetContactsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	if m.getContactsByUserIDFn != nil {
		return m.getContactsByUserIDFn(ctx, userID)
	}
	return []*pkg.ContactDetails{}, nil
}

func (m *mockContactRepo) GetContactRequestsByUserID(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	if m.getContactRequestsByUserIDFn != nil {
		return m.getContactRequestsByUserIDFn(ctx, userID)
	}
	return []*pkg.ContactDetails{}, nil
}

func (m *mockContactRepo) GetContactByPair(ctx context.Context, userID1, userID2 int64) (*Contact, error) {
	if m.getContactByPairFn != nil {
		return m.getContactByPairFn(ctx, userID1, userID2)
	}
	return nil, sql.ErrNoRows
}

func (m *mockContactRepo) CreateContact(ctx context.Context, c *Contact) error {
	if m.createContactFn != nil {
		return m.createContactFn(ctx, c)
	}
	return nil
}

func (m *mockContactRepo) UpdateContact(ctx context.Context, c *Contact) error {
	if m.updateContactFn != nil {
		return m.updateContactFn(ctx, c)
	}
	return nil
}

func (m *mockContactRepo) DeleteContact(ctx context.Context, id int64) error {
	if m.deleteContactFn != nil {
		return m.deleteContactFn(ctx, id)
	}
	return nil
}

// ---- Mock ChallengeRepository ----

type mockChallengeRepo struct {
	createChallengeFn func(ctx context.Context, challenge *Challenge) error
	getChallengeFn    func(ctx context.Context, id int64) (*Challenge, error)
	deleteChallengeFn func(ctx context.Context, id int64, nonce string) error
}

func (m *mockChallengeRepo) WithTx(_ *sql.Tx) *ChallengeStore { return nil }

func (m *mockChallengeRepo) CreateChallenge(ctx context.Context, challenge *Challenge) error {
	if m.createChallengeFn != nil {
		return m.createChallengeFn(ctx, challenge)
	}
	return nil
}

func (m *mockChallengeRepo) GetChallenge(ctx context.Context, id int64) (*Challenge, error) {
	if m.getChallengeFn != nil {
		return m.getChallengeFn(ctx, id)
	}
	return &Challenge{UserID: id, Nonce: "test-nonce"}, nil
}

func (m *mockChallengeRepo) DeleteChallenge(ctx context.Context, id int64, nonce string) error {
	if m.deleteChallengeFn != nil {
		return m.deleteChallengeFn(ctx, id, nonce)
	}
	return nil
}
