package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

var jwtKey = []byte("change-this-secret") // sync with handlers - or inject

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		sub := claims["sub"].(string)
		ctx := context.WithValue(r.Context(), UserIDKey, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
