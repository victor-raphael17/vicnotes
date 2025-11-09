package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"vicnotes/backend/utils"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"missing authorization header"}`))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid authorization header"}`))
			return
		}

		claims, err := utils.VerifyToken(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid token"}`))
			return
		}

		// Store user ID in context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
