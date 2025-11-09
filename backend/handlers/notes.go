package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"vicnotes/backend/models"
	"vicnotes/backend/utils"
)

// CreateNote handles note creation
func CreateNote(db *sql.DB, cache *utils.SimpleCache, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)

		var req models.CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Failed to parse request body",
			})
			return
		}

		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "validation_error",
				Message: "Title is required",
			})
			return
		}

		var noteID int
		err := cb.Call(func() error {
			return db.QueryRow(
				"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
				userID, req.Title, req.Content,
			).Scan(&noteID)
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to create note",
			})
			return
		}

		note := models.Note{
			ID:      noteID,
			UserID:  userID,
			Title:   req.Title,
			Content: req.Content,
		}

		// Invalidate cache for user's notes
		cache.Delete(fmt.Sprintf("notes:%d", userID))

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

// ListNotes handles listing user's notes
func ListNotes(db *sql.DB, cache *utils.SimpleCache, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)
		cacheKey := fmt.Sprintf("notes:%d", userID)

		// Try to get from cache first
		if cached, ok := cache.Get(cacheKey); ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(cached)
			return
		}

		var notes []models.Note
		err := cb.Call(func() error {
			rows, err := db.Query(
				"SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
				userID,
			)
			if err != nil {
				return err
			}
			defer rows.Close()

			notes = []models.Note{}
			for rows.Next() {
				var note models.Note
				if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
					return err
				}
				notes = append(notes, note)
			}
			return nil
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to fetch notes",
			})
			return
		}

		// Cache the result for 5 minutes
		cache.Set(cacheKey, notes, 5*time.Minute)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	}
}

// GetNote handles fetching a single note
func GetNote(db *sql.DB, cache *utils.SimpleCache, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)
		vars := mux.Vars(r)
		noteID, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Invalid note ID",
			})
			return
		}

		cacheKey := fmt.Sprintf("note:%d", noteID)

		// Try to get from cache first
		if cached, ok := cache.Get(cacheKey); ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(cached)
			return
		}

		var note models.Note
		err = cb.Call(func() error {
			return db.QueryRow(
				"SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
				noteID, userID,
			).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		})

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "not_found",
				Message: "Note not found",
			})
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to fetch note",
			})
			return
		}

		// Cache the result for 5 minutes
		cache.Set(cacheKey, note, 5*time.Minute)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	}
}

// UpdateNote handles note updates
func UpdateNote(db *sql.DB, cache *utils.SimpleCache, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)
		vars := mux.Vars(r)
		noteID, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Invalid note ID",
			})
			return
		}

		var req models.UpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Failed to parse request body",
			})
			return
		}

		// Verify ownership
		var existingUserID int
		err = cb.Call(func() error {
			return db.QueryRow("SELECT user_id FROM notes WHERE id = $1", noteID).Scan(&existingUserID)
		})
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "not_found",
				Message: "Note not found",
			})
			return
		}

		if existingUserID != userID {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "forbidden",
				Message: "You don't have permission to update this note",
			})
			return
		}

		err = cb.Call(func() error {
			_, err := db.Exec(
				"UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
				req.Title, req.Content, noteID,
			)
			return err
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to update note",
			})
			return
		}

		// Invalidate cache
		cache.Delete(fmt.Sprintf("note:%d", noteID))
		cache.Delete(fmt.Sprintf("notes:%d", userID))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Note updated successfully"})
	}
}

// DeleteNote handles note deletion
func DeleteNote(db *sql.DB, cache *utils.SimpleCache, cb *utils.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)
		vars := mux.Vars(r)
		noteID, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "invalid_request",
				Message: "Invalid note ID",
			})
			return
		}

		// Verify ownership
		var existingUserID int
		err = cb.Call(func() error {
			return db.QueryRow("SELECT user_id FROM notes WHERE id = $1", noteID).Scan(&existingUserID)
		})
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "not_found",
				Message: "Note not found",
			})
			return
		}

		if existingUserID != userID {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "forbidden",
				Message: "You don't have permission to delete this note",
			})
			return
		}

		err = cb.Call(func() error {
			_, err := db.Exec("DELETE FROM notes WHERE id = $1", noteID)
			return err
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to delete note",
			})
			return
		}

		// Invalidate cache
		cache.Delete(fmt.Sprintf("note:%d", noteID))
		cache.Delete(fmt.Sprintf("notes:%d", userID))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Note deleted successfully"})
	}
}
