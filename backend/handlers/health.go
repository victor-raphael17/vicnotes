package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck handles the health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "vicnotes-backend",
	})
}
