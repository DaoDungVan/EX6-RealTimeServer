package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.
			HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, _ := jwt.
			Parse(tokenStr, func(token *jwt.Token) (any, error) {
				return jwtSecret, nil
			})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.
				WithValue(r.Context(), UserIDKey, claims["user_id"])
			next(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}
