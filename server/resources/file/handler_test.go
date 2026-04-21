package file

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

func injectChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func withUserID(r *http.Request, id int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), utils.CtxKeyUserID, id))
}

func newTestHandler(repo FileRepository) *Handler {
	svc := NewService(repo, nil)
	return NewHandler(svc)
}

type mockFileRepo struct {
	withTxFn                    func(tx *sql.Tx) *FileStore
	getFileFn                   func(ctx context.Context, fileID int64) (*File, error)
	createFileFn                func(ctx context.Context, file *File) error
	deleteFileFn                func(ctx context.Context, fileID int64) error
	initUploadSessionFn         func(ctx context.Context, uploadSession *UploadSession) error
	deleteUploadSessionFn       func(ctx context.Context, sessionID int64) error
	insertFileChunkFn           func(ctx context.Context, fileChunk *FileChunk) error
	deleteFileChunksFromSessFn  func(ctx context.Context, sessionID int64) error
	getChunksCountForSessionFn  func(ctx context.Context, sessionID int64) (int64, error)
	updateFileStatusAndPathFn   func(ctx context.Context, fileID int64, status FileStatus, finalPath string) error
	updateUploadSessionStatusFn func(ctx context.Context, sessionID int64, status UploadSessionStatus) error
	getUploadSessionFn          func(ctx context.Context, sessionID int64) (*UploadSession, error)
}

func (m *mockFileRepo) WithTx(tx *sql.Tx) *FileStore {
	if m.withTxFn != nil {
		return m.withTxFn(tx)
	}
	return nil
}

func (m *mockFileRepo) GetFile(ctx context.Context, fileID int64) (*File, error) {
	if m.getFileFn != nil {
		return m.getFileFn(ctx, fileID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFileRepo) CreateFile(ctx context.Context, file *File) error {
	if m.createFileFn != nil {
		return m.createFileFn(ctx, file)
	}
	return nil
}

func (m *mockFileRepo) DeleteFile(ctx context.Context, fileID int64) error {
	if m.deleteFileFn != nil {
		return m.deleteFileFn(ctx, fileID)
	}
	return nil
}

func (m *mockFileRepo) InitUploadSession(ctx context.Context, uploadSession *UploadSession) error {
	if m.initUploadSessionFn != nil {
		return m.initUploadSessionFn(ctx, uploadSession)
	}
	return nil
}

func (m *mockFileRepo) DeleteUploadSession(ctx context.Context, sessionID int64) error {
	if m.deleteUploadSessionFn != nil {
		return m.deleteUploadSessionFn(ctx, sessionID)
	}
	return nil
}

func (m *mockFileRepo) InsertFileChunk(ctx context.Context, fileChunk *FileChunk) error {
	if m.insertFileChunkFn != nil {
		return m.insertFileChunkFn(ctx, fileChunk)
	}
	return nil
}

func (m *mockFileRepo) DeleteFileChunksFromUploadSession(ctx context.Context, sessionID int64) error {
	if m.deleteFileChunksFromSessFn != nil {
		return m.deleteFileChunksFromSessFn(ctx, sessionID)
	}
	return nil
}

func (m *mockFileRepo) GetChunksCountForSession(ctx context.Context, sessionID int64) (int64, error) {
	if m.getChunksCountForSessionFn != nil {
		return m.getChunksCountForSessionFn(ctx, sessionID)
	}
	return 0, nil
}

func (m *mockFileRepo) UpdateFileStatusAndPath(ctx context.Context, fileID int64, status FileStatus, finalPath string) error {
	if m.updateFileStatusAndPathFn != nil {
		return m.updateFileStatusAndPathFn(ctx, fileID, status, finalPath)
	}
	return nil
}

func (m *mockFileRepo) UpdateUploadSessionStatus(ctx context.Context, sessionID int64, status UploadSessionStatus) error {
	if m.updateUploadSessionStatusFn != nil {
		return m.updateUploadSessionStatusFn(ctx, sessionID, status)
	}
	return nil
}

func (m *mockFileRepo) GetUploadSession(ctx context.Context, sessionID int64) (*UploadSession, error) {
	if m.getUploadSessionFn != nil {
		return m.getUploadSessionFn(ctx, sessionID)
	}
	return nil, errors.New("not implemented")
}

func TestHandler_FileInit_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockFileRepo{})

	req := httptest.NewRequest(http.MethodPost, "/files/init", bytes.NewReader([]byte("bad-json")))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.fileInit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_FileInit_UploaderMismatch(t *testing.T) {
	h := newTestHandler(&mockFileRepo{})

	body, _ := json.Marshal(pkg.InitFileUploadRequest{UploaderID: 2})
	req := httptest.NewRequest(http.MethodPost, "/files/init", bytes.NewReader(body))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.fileInit(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestHandler_UploadChunk_OK(t *testing.T) {
	called := false
	repo := &mockFileRepo{
		insertFileChunkFn: func(_ context.Context, chunk *FileChunk) error {
			called = true
			if chunk.SessionID != 12 {
				t.Fatalf("session ID: got %d, want %d", chunk.SessionID, 12)
			}
			if chunk.Index != 0 {
				t.Fatalf("chunk index: got %d, want %d", chunk.Index, 0)
			}
			return nil
		},
	}
	h := newTestHandler(repo)

	sessionDir := filepath.Join(string(TmpUploadsPath), "session-12")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	body, _ := json.Marshal(pkg.UploadFileChunkRequest{
		ChunkIndex: 0,
		ChunkData:  []byte("test"),
		Checksum:   "abc",
		Size:       4,
	})
	req := httptest.NewRequest(http.MethodPost, "/files/upload/12/chunks", bytes.NewReader(body))
	req = injectChiParam(req, "session_id", "12")
	w := httptest.NewRecorder()

	h.uploadChunk(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusAccepted)
	}
	if !called {
		t.Error("expected InsertFileChunk to be called")
	}
}

func TestHandler_UploadChunk_InvalidSessionID(t *testing.T) {
	h := newTestHandler(&mockFileRepo{})

	req := httptest.NewRequest(http.MethodPost, "/files/upload/abc/chunks", bytes.NewReader([]byte("{}")))
	req = injectChiParam(req, "session_id", "abc")
	w := httptest.NewRecorder()

	h.uploadChunk(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_UploadChunk_ServiceError(t *testing.T) {
	repo := &mockFileRepo{
		insertFileChunkFn: func(_ context.Context, _ *FileChunk) error {
			return errors.New("insert failed")
		},
	}
	h := newTestHandler(repo)

	sessionDir := filepath.Join(string(TmpUploadsPath), "session-12")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	body, _ := json.Marshal(pkg.UploadFileChunkRequest{ChunkIndex: 1, ChunkData: []byte("a"), Size: 1})
	req := httptest.NewRequest(http.MethodPost, "/files/upload/12/chunks", bytes.NewReader(body))
	req = injectChiParam(req, "session_id", "12")
	w := httptest.NewRecorder()

	h.uploadChunk(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestHandler_AssembleFile_InvalidSessionID(t *testing.T) {
	h := newTestHandler(&mockFileRepo{})

	req := httptest.NewRequest(http.MethodPost, "/files/upload/abc/assemble", nil)
	req = injectChiParam(req, "session_id", "abc")
	w := httptest.NewRecorder()

	h.assembleFile(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_AssembleFile_ServiceError(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return nil, errors.New("repo error")
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/files/upload/7/assemble", nil)
	req = injectChiParam(req, "session_id", "7")
	w := httptest.NewRecorder()

	h.assembleFile(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestHandler_DownloadFile_InvalidID(t *testing.T) {
	h := newTestHandler(&mockFileRepo{})

	req := httptest.NewRequest(http.MethodGet, "/files/download/abc", nil)
	req = injectChiParam(req, "file_id", "abc")
	w := httptest.NewRecorder()

	h.downloadFile(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_DownloadFile_GetFileError(t *testing.T) {
	repo := &mockFileRepo{
		getFileFn: func(_ context.Context, _ int64) (*File, error) {
			return nil, errors.New("db error")
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/files/download/5", nil)
	req = injectChiParam(req, "file_id", "5")
	w := httptest.NewRecorder()

	h.downloadFile(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestHandler_DownloadFile_OK(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	content := []byte("hello-download")
	if err := os.WriteFile(filePath, content, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	repo := &mockFileRepo{
		getFileFn: func(_ context.Context, fileID int64) (*File, error) {
			return &File{
				ID:             fileID,
				FileName:       "file.txt",
				Extension:      ".txt",
				ConversationID: 1,
				UploaderID:     1,
				Size:           int64(len(content)),
				Status:         FileStatusReady,
				Checksum:       "",
				StoragePath:    filePath,
				CreatedAt:      time.Now(),
			}, nil
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/files/download/5", nil)
	req = injectChiParam(req, "file_id", "5")
	w := httptest.NewRecorder()

	h.downloadFile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Body.String(); got != string(content) {
		t.Errorf("body: got %q, want %q", got, string(content))
	}
}
