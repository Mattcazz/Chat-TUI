package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

// newTestService builds a Service with no TxManager.
// Tests that go through s.tx (ContactRequest happy path) need a real DB
// and belong in integration tests.
func newTestService(userRepo UserRepository, contactRepo ContactRepository, challengeRepo ChallengeRepository) *Service {
	return &Service{
		userRepo:      userRepo,
		contactRepo:   contactRepo,
		challengeRepo: challengeRepo,
	}
}

// ---- CreateUser ----

func TestService_CreateUser_OK(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows // user does not exist yet
		},
		createUserFn: func(_ context.Context, u *User) (*User, error) {
			u.ID = 42
			return u, nil
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	resp, err := svc.CreateUser(context.Background(), "pk-alice", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Username != "alice" {
		t.Errorf("username: got %q, want %q", resp.Username, "alice")
	}
	if resp.ID != 42 {
		t.Errorf("ID: got %d, want 42", resp.ID)
	}
}

func TestService_CreateUser_AlreadyExists(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: 1, Username: "alice"}, nil // user already registered
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	_, err := svc.CreateUser(context.Background(), "pk-alice", "alice")
	if err == nil {
		t.Fatal("expected error for duplicate user, got nil")
	}
}

func TestService_CreateUser_RepoError(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows
		},
		createUserFn: func(_ context.Context, _ *User) (*User, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	_, err := svc.CreateUser(context.Background(), "pk-alice", "alice")
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

// ---- DeleteUser ----

func TestService_DeleteUser_OK(t *testing.T) {
	svc := newTestService(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})
	if err := svc.DeleteUser(context.Background(), 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestService_DeleteUser_Error(t *testing.T) {
	userRepo := &mockUserRepo{
		deleteUserFn: func(_ context.Context, _ int64) error {
			return errors.New("delete failed")
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	if err := svc.DeleteUser(context.Background(), 1); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---- GenerateChallenge ----

func TestService_GenerateChallenge_OK(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: 5, PublicKey: "pk-bob"}, nil
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	nonce, err := svc.GenerateChallenge(context.Background(), "pk-bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonce == "" {
		t.Error("expected a non-empty nonce")
	}
}

func TestService_GenerateChallenge_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	_, err := svc.GenerateChallenge(context.Background(), "unknown-pk")
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
	if !IsUserDoesNotExistError(err) {
		t.Errorf("expected UserDoesNotExistError, got %T: %v", err, err)
	}
}

func TestService_GenerateChallenge_ChallengeRepoError(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: 5}, nil
		},
	}
	challengeRepo := &mockChallengeRepo{
		createChallengeFn: func(_ context.Context, _ *Challenge) error {
			return errors.New("insert failed")
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, challengeRepo)

	_, err := svc.GenerateChallenge(context.Background(), "pk-bob")
	if err == nil {
		t.Fatal("expected error from challenge repo, got nil")
	}
}

// ---- GetContacts ----

func TestService_GetContacts_ReturnsContacts(t *testing.T) {
	expected := []*pkg.ContactDetails{
		{ID: 1, Username: "bob"},
		{ID: 2, Username: "charlie"},
	}
	contactRepo := &mockContactRepo{
		getContactsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return expected, nil
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	contacts, err := svc.GetContacts(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 2 {
		t.Errorf("got %d contacts, want 2", len(contacts))
	}
}

func TestService_GetContacts_NilBecomesEmptySlice(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return nil, nil // DB returned zero rows
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	contacts, err := svc.GetContacts(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contacts == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(contacts) != 0 {
		t.Errorf("expected 0 contacts, got %d", len(contacts))
	}
}

func TestService_GetContacts_Error(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	_, err := svc.GetContacts(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---- GetContactRequests ----

func TestService_GetContactRequests_ReturnsRequests(t *testing.T) {
	expected := []*pkg.ContactDetails{{ID: 10, Username: "dave"}}
	contactRepo := &mockContactRepo{
		getContactRequestsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return expected, nil
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	reqs, err := svc.GetContactRequests(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reqs) != 1 {
		t.Errorf("got %d requests, want 1", len(reqs))
	}
}

func TestService_GetContactRequests_NilBecomesEmptySlice(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactRequestsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return nil, nil
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	reqs, err := svc.GetContactRequests(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reqs == nil {
		t.Error("expected empty slice, got nil")
	}
}

func TestService_GetContactRequests_Error(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactRequestsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	_, err := svc.GetContactRequests(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---- BlockContact / UnblockContact ----

func TestService_BlockContact_SetsStatusAndID(t *testing.T) {
	var captured *Contact
	contactRepo := &mockContactRepo{
		updateContactFn: func(_ context.Context, c *Contact) error {
			captured = c
			return nil
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	if err := svc.BlockContact(context.Background(), 1, 99); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Status != StatusBlocked {
		t.Errorf("status: got %q, want %q", captured.Status, StatusBlocked)
	}
	if captured.ID != 99 {
		t.Errorf("contact ID: got %d, want 99", captured.ID)
	}
}

func TestService_BlockContact_Error(t *testing.T) {
	contactRepo := &mockContactRepo{
		updateContactFn: func(_ context.Context, _ *Contact) error {
			return errors.New("update failed")
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	if err := svc.BlockContact(context.Background(), 1, 99); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_UnblockContact_SetsStatusAccepted(t *testing.T) {
	var captured *Contact
	contactRepo := &mockContactRepo{
		updateContactFn: func(_ context.Context, c *Contact) error {
			captured = c
			return nil
		},
	}
	svc := newTestService(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	if err := svc.UnblockContact(context.Background(), 1, 99); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Status != StatusAccept {
		t.Errorf("status: got %q, want %q", captured.Status, StatusAccept)
	}
}

// ---- ContactRequest (pre-tx error paths only) ----
// The happy path and mutual-accept flow require a real TxManager and DB.
// Those belong in integration tests.

func TestService_ContactRequest_FromUserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByIDFn: func(_ context.Context, _ int64) (*User, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	err := svc.ContactRequest(context.Background(), 1, "pk-bob", "")
	if err == nil {
		t.Fatal("expected error for missing from-user, got nil")
	}
}

func TestService_ContactRequest_SelfRequest(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByIDFn: func(_ context.Context, id int64) (*User, error) {
			return &User{ID: id, PublicKey: "pk-alice"}, nil
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	err := svc.ContactRequest(context.Background(), 1, "pk-alice", "")
	if err == nil {
		t.Fatal("expected error for self-request, got nil")
	}
}

func TestService_ContactRequest_ToUserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByIDFn: func(_ context.Context, id int64) (*User, error) {
			return &User{ID: id, PublicKey: "pk-alice"}, nil
		},
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := newTestService(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	err := svc.ContactRequest(context.Background(), 1, "pk-unknown", "")
	if err == nil {
		t.Fatal("expected error for missing to-user, got nil")
	}
}
