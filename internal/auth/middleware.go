package auth

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "username"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, err := ValidateJWT(r)
		if err != nil {
			log.Printf("Unauthorized request: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(userContextKey).(string)
	return username, ok
}
