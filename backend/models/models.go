package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Note represents a note created by a user
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CreateNoteRequest represents the create note request payload
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateNoteRequest represents the update note request payload
type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
