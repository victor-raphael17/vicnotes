package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"vicnotes/backend/models"
)

// CreateNote handles note creation
func CreateNote(db *sql.DB) http.HandlerFunc {
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
		err := db.QueryRow(
			"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
			userID, req.Title, req.Content,
		).Scan(&noteID)

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

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

// ListNotes handles listing user's notes
func ListNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("user_id").(int)

		rows, err := db.Query(
			"SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
			userID,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to fetch notes",
			})
			return
		}
		defer rows.Close()

		notes := []models.Note{}
		for rows.Next() {
			var note models.Note
			if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.ErrorResponse{
					Error:   "server_error",
					Message: "Failed to parse notes",
				})
				return
			}
			notes = append(notes, note)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	}
}

// GetNote handles fetching a single note
func GetNote(db *sql.DB) http.HandlerFunc {
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

		var note models.Note
		err = db.QueryRow(
			"SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)

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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	}
}

// UpdateNote handles note updates
func UpdateNote(db *sql.DB) http.HandlerFunc {
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
		err = db.QueryRow("SELECT user_id FROM notes WHERE id = $1", noteID).Scan(&existingUserID)
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

		_, err = db.Exec(
			"UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
			req.Title, req.Content, noteID,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to update note",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Note updated successfully"})
	}
}

// DeleteNote handles note deletion
func DeleteNote(db *sql.DB) http.HandlerFunc {
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
		err = db.QueryRow("SELECT user_id FROM notes WHERE id = $1", noteID).Scan(&existingUserID)
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

		_, err = db.Exec("DELETE FROM notes WHERE id = $1", noteID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "server_error",
				Message: "Failed to delete note",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Note deleted successfully"})
	}
}
