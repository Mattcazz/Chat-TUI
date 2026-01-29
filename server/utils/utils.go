package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

type headerValues string

const PublicKey headerValues = "public_key"
const FingerPirnt headerValues = "finger_print"

type ctxKeyUserID string

const CtxKeyUserID ctxKeyUserID = "userID"

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func CtxWithUser(ctx context.Context, user int) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, user)
}

func UserFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(CtxKeyUserID).(int)
	return id, ok
}
