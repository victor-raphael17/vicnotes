package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"vicnotes/backend/models"
	"vicnotes/backend/utils"
)

// Register handles user registration
func Register(db *sql.DB, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Failed to parse request body",
			})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "validation_error",
				Message: "Email and password are required",
			})
			return
		}

		// Hash password
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to process password",
			})
			return
		}

		// Insert user with circuit breaker
		var userID int
		err = cb.Call(func() error {
			return db.QueryRow(
				"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
				req.Email, passwordHash,
			).Scan(&userID)
		})

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "user_exists",
				Message: "User with this email already exists",
			})
			return
		}

		// Generate token
		token, err := utils.GenerateToken(userID, req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to generate token",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.AuthResponse{
			Token: token,
			User: models.User{
				ID:    userID,
				Email: req.Email,
			},
		})
	}
}

// Login handles user login
func Login(db *sql.DB, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Failed to parse request body",
			})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "validation_error",
				Message: "Email and password are required",
			})
			return
		}

		// Get user with circuit breaker
		var user models.User
		var passwordHash string
		err := cb.Call(func() error {
			return db.QueryRow(
				"SELECT id, email, password_hash FROM users WHERE email = $1",
				req.Email,
			).Scan(&user.ID, &user.Email, &passwordHash)
		})

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to query user",
			})
			return
		}

		// Verify password
		if !utils.VerifyPassword(passwordHash, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
			return
		}

		// Generate token
		token, err := utils.GenerateToken(user.ID, user.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to generate token",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.AuthResponse{
			Token: token,
			User:  user,
		})
	}
}
