package file

import (
	"encoding/json"
	"fmt"
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

	fileInitReq := &pkg.InitFileUploadRequest{}

	if err := json.NewDecoder(r.Body).Decode(fileInitReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
	}

	if senderID != fileInitReq.UploaderID {
		utils.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("uploader ID does not match authenticated user"))
	}

	resp, err := h.service.InitFileUpload(r.Context(), fileInitReq)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusAccepted, resp)
}

func (h *Handler) uploadChunk(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid session ID format"))
	}

	var uploadChunkReq *pkg.UploadFileChunkRequest

	if json.NewDecoder(r.Body).Decode(uploadChunkReq) != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
	}

	err = h.service.UploadFileChunk(r.Context(), int64(sessionID), uploadChunkReq)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) assembleFile(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid session ID format"))
	}

	if err := h.service.FinalizeFileUpload(r.Context(), int64(sessionID)); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) statusCheck(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) cancelUpload(w http.ResponseWriter, r *http.Request) {
}
