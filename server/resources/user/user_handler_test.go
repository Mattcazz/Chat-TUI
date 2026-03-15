package user

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

// ---- Helpers ----

func injectChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func withUserID(r *http.Request, id int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), utils.CtxKeyUserID, id))
}

func newTestHandler(userRepo UserRepository, contactRepo ContactRepository, challengeRepo ChallengeRepository) *Handler {
	svc := newTestService(userRepo, contactRepo, challengeRepo)
	return NewHandler(svc)
}

// ---- registerUser ----

func TestHandler_RegisterUser_OK(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows
		},
		createUserFn: func(_ context.Context, u *User) (*User, error) {
			u.ID = 1
			return u, nil
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	body, _ := json.Marshal(pkg.RegisterRequest{Username: "alice", PublicKey: "pk-alice"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.registerUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_RegisterUser_MissingFields(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	// Username is empty
	body, _ := json.Marshal(pkg.RegisterRequest{PublicKey: "pk-alice"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.registerUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_RegisterUser_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.registerUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_RegisterUser_ServiceError(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: 1}, nil // already exists
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	body, _ := json.Marshal(pkg.RegisterRequest{Username: "alice", PublicKey: "pk-alice"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.registerUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- login (challenge phase) ----

func TestHandler_Login_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("bad")))
	w := httptest.NewRecorder()

	h.login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_Login_MissingPublicKey(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	body, _ := json.Marshal(pkg.LoginRequest{})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_Login_ChallengePhase_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return nil, sql.ErrNoRows
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	// No signature => challenge phase
	body, _ := json.Marshal(pkg.LoginRequest{PublicKey: "pk-unknown"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.login(w, req)

	// Handler returns TemporaryRedirect when user doesn't exist
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusTemporaryRedirect)
	}
}

func TestHandler_Login_ChallengePhase_OK(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByPublicKeyFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: 1, PublicKey: "pk-alice"}, nil
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	body, _ := json.Marshal(pkg.LoginRequest{PublicKey: "pk-alice"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp pkg.ChallengeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode challenge response: %v", err)
	}
	if resp.Nonce == "" {
		t.Error("expected non-empty nonce in challenge response")
	}
}

// ---- deleteUser ----

func TestHandler_DeleteUser_OK(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.deleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_DeleteUser_ServiceError(t *testing.T) {
	userRepo := &mockUserRepo{
		deleteUserFn: func(_ context.Context, _ int64) error {
			return errors.New("delete failed")
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.deleteUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- getContacts ----

func TestHandler_GetContacts_OK(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return []*pkg.ContactDetails{{ID: 1, Username: "bob"}}, nil
		},
	}
	h := newTestHandler(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.getContacts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_GetContacts_Error(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return nil, errors.New("db error")
		},
	}
	h := newTestHandler(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.getContacts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- postContactRequest ----

func TestHandler_PostContactRequest_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/contacts/requests", bytes.NewReader([]byte("bad")))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.postContactRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_PostContactRequest_SelfRequest(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserByIDFn: func(_ context.Context, id int64) (*User, error) {
			return &User{ID: id, PublicKey: "pk-alice"}, nil
		},
	}
	h := newTestHandler(userRepo, &mockContactRepo{}, &mockChallengeRepo{})

	body, _ := json.Marshal(pkg.PostContactRequest{PublicKey: "pk-alice"})
	req := httptest.NewRequest(http.MethodPost, "/contacts/requests", bytes.NewReader(body))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.postContactRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- getContactRequests ----

func TestHandler_GetContactRequests_OK(t *testing.T) {
	contactRepo := &mockContactRepo{
		getContactRequestsByUserIDFn: func(_ context.Context, _ int64) ([]*pkg.ContactDetails, error) {
			return []*pkg.ContactDetails{{ID: 5, Username: "eve"}}, nil
		},
	}
	h := newTestHandler(&mockUserRepo{}, contactRepo, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodGet, "/contacts/requests", nil)
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.getContactRequests(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

// ---- blockContact / unblockContact ----

func TestHandler_BlockContact_OK(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/contacts/99/block", nil)
	req = withUserID(req, 1)
	req = injectChiParam(req, "contact_id", "99")
	w := httptest.NewRecorder()

	h.blockContact(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_BlockContact_InvalidID(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/contacts/abc/block", nil)
	req = withUserID(req, 1)
	req = injectChiParam(req, "contact_id", "abc")
	w := httptest.NewRecorder()

	h.blockContact(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_UnblockContact_OK(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/contacts/99/unblock", nil)
	req = withUserID(req, 1)
	req = injectChiParam(req, "contact_id", "99")
	w := httptest.NewRecorder()

	h.unblockContact(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_UnblockContact_InvalidID(t *testing.T) {
	h := newTestHandler(&mockUserRepo{}, &mockContactRepo{}, &mockChallengeRepo{})

	req := httptest.NewRequest(http.MethodPost, "/contacts/abc/unblock", nil)
	req = withUserID(req, 1)
	req = injectChiParam(req, "contact_id", "abc")
	w := httptest.NewRecorder()

	h.unblockContact(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}
