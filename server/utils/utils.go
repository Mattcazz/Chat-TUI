package utils

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
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

func WriteJsonMsg(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"message": msg})
}

func WriteJSONError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func CtxWithUser(ctx context.Context, user int) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, user)
}

func UserFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(CtxKeyUserID).(int)
	return id, ok
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func StrictDecoder(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	return decoder.Decode(v)
}
