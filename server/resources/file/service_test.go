package file

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

func newTestService(repo FileRepository) *Service {
	return NewService(repo, nil)
}

func TestService_UploadFileChunk_OK(t *testing.T) {
	repo := &mockFileRepo{}
	svc := newTestService(repo)

	sessionDir := filepath.Join(string(TmpUploadsPath), "session-25")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	req := &pkg.UploadFileChunkRequest{
		ChunkIndex: 2,
		ChunkData:  []byte("abcd"),
		Checksum:   "sum",
		Size:       4,
	}

	if err := svc.UploadFileChunk(context.Background(), 25, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chunkPath := filepath.Join(sessionDir, "chunk-2.bin")
	if _, err := os.Stat(chunkPath); err != nil {
		t.Fatalf("expected chunk file to exist: %v", err)
	}
}

func TestService_UploadFileChunk_SizeMismatch(t *testing.T) {
	repo := &mockFileRepo{}
	svc := newTestService(repo)

	sessionDir := filepath.Join(string(TmpUploadsPath), "session-26")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	req := &pkg.UploadFileChunkRequest{
		ChunkIndex: 1,
		ChunkData:  []byte("abcd"),
		Checksum:   "sum",
		Size:       5,
	}

	err := svc.UploadFileChunk(context.Background(), 26, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "expected") {
		t.Fatalf("unexpected error: %v", err)
	}

	chunkPath := filepath.Join(sessionDir, "chunk-1.bin")
	if _, err := os.Stat(chunkPath); !os.IsNotExist(err) {
		t.Fatalf("expected chunk file to be removed after mismatch, got err: %v", err)
	}
}

func TestService_UploadFileChunk_InsertErrorRemovesFile(t *testing.T) {
	repo := &mockFileRepo{
		insertFileChunkFn: func(_ context.Context, _ *FileChunk) error {
			return errors.New("insert failed")
		},
	}
	svc := newTestService(repo)

	sessionDir := filepath.Join(string(TmpUploadsPath), "session-27")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	req := &pkg.UploadFileChunkRequest{ChunkIndex: 0, ChunkData: []byte("abc"), Size: 3}
	err := svc.UploadFileChunk(context.Background(), 27, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	chunkPath := filepath.Join(sessionDir, "chunk-0.bin")
	if _, err := os.Stat(chunkPath); !os.IsNotExist(err) {
		t.Fatalf("expected chunk file to be removed on repo error, got err: %v", err)
	}
}

func TestService_FinalizeFileUpload_SessionNotFound(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return nil, nil
		},
	}
	svc := newTestService(repo)

	err := svc.FinalizeFileUpload(context.Background(), 100)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FinalizeFileUpload_SessionError(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(repo)

	if err := svc.FinalizeFileUpload(context.Background(), 100); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FinalizeFileUpload_InvalidStatus(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return &UploadSession{ID: 1, FileID: 1, TotalChunks: 2, Status: FileSessionStatusCompleted}, nil
		},
	}
	svc := newTestService(repo)

	err := svc.FinalizeFileUpload(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FinalizeFileUpload_ChunkCountMismatch(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return &UploadSession{ID: 1, FileID: 1, TotalChunks: 2, Status: FileSessionStatusUploading}, nil
		},
		getChunksCountForSessionFn: func(_ context.Context, _ int64) (int64, error) {
			return 1, nil
		},
	}
	svc := newTestService(repo)

	err := svc.FinalizeFileUpload(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FinalizeFileUpload_GetChunksError(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return &UploadSession{ID: 1, FileID: 1, TotalChunks: 2, Status: FileSessionStatusUploading}, nil
		},
		getChunksCountForSessionFn: func(_ context.Context, _ int64) (int64, error) {
			return 0, errors.New("count error")
		},
	}
	svc := newTestService(repo)

	if err := svc.FinalizeFileUpload(context.Background(), 1); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FinalizeFileUpload_GetFileError(t *testing.T) {
	repo := &mockFileRepo{
		getUploadSessionFn: func(_ context.Context, _ int64) (*UploadSession, error) {
			return &UploadSession{ID: 1, FileID: 1, TotalChunks: 1, Status: FileSessionStatusUploading}, nil
		},
		getChunksCountForSessionFn: func(_ context.Context, _ int64) (int64, error) {
			return 1, nil
		},
		getFileFn: func(_ context.Context, _ int64) (*File, error) {
			return nil, errors.New("get file error")
		},
	}
	svc := newTestService(repo)

	if err := svc.FinalizeFileUpload(context.Background(), 1); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_GetFile_OK(t *testing.T) {
	expected := &File{ID: 5, FileName: "x.txt", StoragePath: "/tmp/x.txt", CreatedAt: time.Now()}
	repo := &mockFileRepo{
		getFileFn: func(_ context.Context, fileID int64) (*File, error) {
			if fileID != 5 {
				t.Fatalf("file ID: got %d, want 5", fileID)
			}
			return expected, nil
		},
	}
	svc := newTestService(repo)

	file, err := svc.GetFile(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.ID != expected.ID {
		t.Fatalf("got file ID %d, want %d", file.ID, expected.ID)
	}
}

func TestService_GetFile_Error(t *testing.T) {
	repo := &mockFileRepo{
		getFileFn: func(_ context.Context, _ int64) (*File, error) {
			return nil, errors.New("db")
		},
	}
	svc := newTestService(repo)

	if _, err := svc.GetFile(context.Background(), 5); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_DeleteSessionChunks_NoOp(t *testing.T) {
	svc := newTestService(&mockFileRepo{})
	if err := svc.DeleteSessionChunks(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssembleFile_OK(t *testing.T) {
	sessionDir := filepath.Join(string(TmpUploadsPath), "session-201")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	chunk1 := []byte("hello ")
	chunk2 := []byte("world")
	if err := os.WriteFile(filepath.Join(sessionDir, "chunk-0.bin"), chunk1, 0o600); err != nil {
		t.Fatalf("write chunk 0: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "chunk-1.bin"), chunk2, 0o600); err != nil {
		t.Fatalf("write chunk 1: %v", err)
	}

	finalDir := t.TempDir()
	finalPath := filepath.Join(finalDir, "assembled.bin")
	finalFile, err := os.Create(finalPath)
	if err != nil {
		t.Fatalf("create final file: %v", err)
	}
	t.Cleanup(func() { _ = finalFile.Close() })

	checksum, err := pkg.CalculateFileChecksum(finalPath)
	if err != nil {
		t.Fatalf("calculate checksum before write: %v", err)
	}
	_ = checksum

	expected := append(chunk1, chunk2...)
	if err := os.WriteFile(finalPath+".expected", expected, 0o600); err != nil {
		t.Fatalf("write expected file: %v", err)
	}
	expectedChecksum, err := pkg.CalculateFileChecksum(finalPath + ".expected")
	if err != nil {
		t.Fatalf("calculate expected checksum: %v", err)
	}

	file := &File{Size: int64(len(expected)), Checksum: expectedChecksum}
	session := &UploadSession{ID: 201, TotalChunks: 2}

	if err := assembleFile(finalFile, finalPath, file, session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssembleFile_SizeMismatch(t *testing.T) {
	sessionDir := filepath.Join(string(TmpUploadsPath), "session-202")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	if err := os.WriteFile(filepath.Join(sessionDir, "chunk-0.bin"), []byte("abc"), 0o600); err != nil {
		t.Fatalf("write chunk: %v", err)
	}

	finalDir := t.TempDir()
	finalPath := filepath.Join(finalDir, "assembled.bin")
	finalFile, err := os.Create(finalPath)
	if err != nil {
		t.Fatalf("create final file: %v", err)
	}
	t.Cleanup(func() { _ = finalFile.Close() })

	file := &File{Size: 5, Checksum: "ignored"}
	session := &UploadSession{ID: 202, TotalChunks: 1}

	err = assembleFile(finalFile, finalPath, file, session)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "size") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssembleFile_ChecksumMismatch(t *testing.T) {
	sessionDir := filepath.Join(string(TmpUploadsPath), "session-203")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("mkdir session dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sessionDir) })

	if err := os.WriteFile(filepath.Join(sessionDir, "chunk-0.bin"), []byte("abc"), 0o600); err != nil {
		t.Fatalf("write chunk: %v", err)
	}

	finalDir := t.TempDir()
	finalPath := filepath.Join(finalDir, "assembled.bin")
	finalFile, err := os.Create(finalPath)
	if err != nil {
		t.Fatalf("create final file: %v", err)
	}
	t.Cleanup(func() { _ = finalFile.Close() })

	file := &File{Size: 3, Checksum: "wrong"}
	session := &UploadSession{ID: 203, TotalChunks: 1}

	err = assembleFile(finalFile, finalPath, file, session)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("unexpected error: %v", err)
	}
}
