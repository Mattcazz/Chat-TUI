package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(secret []byte, UserId int64) (string, error) {
	claims := userClaims{
		UserID: int64(UserId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func JWTAuth(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		tokenString = r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*userClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.CtxKeyUserID, claims.UserID)

		next(w, r.WithContext(ctx))

	}
}
