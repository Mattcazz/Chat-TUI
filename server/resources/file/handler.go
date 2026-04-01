package file

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/files", func(r chi.Router) {
		r.Post("/init", middleware.JWTAuth(h.fileInit))
		r.Post("/upload/{session_id}/chunks", middleware.JWTAuth(h.uploadChunk))
		r.Post("/upload/{session_id}/assemble", middleware.JWTAuth(h.assembleFile))
		r.Get("/upload/{session_id}/status", middleware.JWTAuth(h.statusCheck))
		r.Get("/download/{file_id}", middleware.JWTAuth(h.downloadFile))
		r.Delete("/upload/{session_id}", middleware.JWTAuth(h.cancelUpload))
	})
}

func (h *Handler) fileInit(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(utils.CtxKeyUserID)
	log.Printf("Handler.fileInit: User ID %d initiating file upload", senderID)

	fileInitReq := &pkg.InitFileUploadRequest{}

	if err := json.NewDecoder(r.Body).Decode(fileInitReq); err != nil {
		log.Printf("Handler.fileInit: Failed to decode request body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.fileInit: File upload request - filename: %s, size: %d bytes, chunks: %d, conversation ID: %d",
		fileInitReq.FileName, fileInitReq.TotalSize, fileInitReq.TotalChunks, fileInitReq.ConversationID)

	if senderID != fileInitReq.UploaderID {
		log.Printf("Handler.fileInit: Uploader ID mismatch - authenticated user: %d, request uploader: %d",
			senderID, fileInitReq.UploaderID)
		utils.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("uploader ID does not match authenticated user"))
		return
	}

	log.Printf("Handler.fileInit: Calling service to initialize upload")
	resp, err := h.service.InitFileUpload(r.Context(), fileInitReq)
	if err != nil {
		log.Printf("Handler.fileInit: Service failed to initialize upload: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Handler.fileInit: Successfully initialized upload - session ID: %d, file ID: %d",
		resp.SessionID, resp.FileID)
	utils.WriteJSON(w, http.StatusAccepted, resp)
}

func (h *Handler) uploadChunk(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	log.Printf("Handler.uploadChunk: Received chunk upload request for session ID %s", sessionIDStr)

	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		log.Printf("Handler.uploadChunk: Invalid session ID format: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid session ID format"))
		return
	}

	var uploadChunkReq *pkg.UploadFileChunkRequest

	if err := json.NewDecoder(r.Body).Decode(&uploadChunkReq); err != nil {
		log.Printf("Handler.uploadChunk: Failed to decode request body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.uploadChunk: Uploading chunk %d for session ID %d, size: %d bytes",
		uploadChunkReq.ChunkIndex, sessionID, len(uploadChunkReq.ChunkData))

	err = h.service.UploadFileChunk(r.Context(), int64(sessionID), uploadChunkReq)
	if err != nil {
		log.Printf("Handler.uploadChunk: Service failed to upload chunk: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Handler.uploadChunk: Successfully uploaded chunk %d for session ID %d",
		uploadChunkReq.ChunkIndex, sessionID)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) assembleFile(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	log.Printf("Handler.assembleFile: Received file assembly request for session ID %s", sessionIDStr)

	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		log.Printf("Handler.assembleFile: Invalid session ID format: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid session ID format"))
		return
	}

	log.Printf("Handler.assembleFile: Calling service to finalize upload for session ID %d", sessionID)
	if err := h.service.FinalizeFileUpload(r.Context(), int64(sessionID)); err != nil {
		log.Printf("Handler.assembleFile: Service failed to finalize upload: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Handler.assembleFile: Successfully assembled file for session ID %d", sessionID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) statusCheck(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	log.Printf("Handler.statusCheck: Status check requested for session ID %s", sessionIDStr)
	log.Printf("Handler.statusCheck: Endpoint not yet implemented")
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
	fileIDStr := chi.URLParam(r, "file_id")
	log.Printf("Handler.downloadFile: Download requested for file ID %s", fileIDStr)
	log.Printf("Handler.downloadFile: Endpoint not yet implemented")
}

func (h *Handler) cancelUpload(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	log.Printf("Handler.cancelUpload: Cancel upload requested for session ID %s", sessionIDStr)
	log.Printf("Handler.cancelUpload: Endpoint not yet implemented")
}
