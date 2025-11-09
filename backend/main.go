package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"vicnotes/backend/config"
	"vicnotes/backend/database"
	"vicnotes/backend/handlers"
	"vicnotes/backend/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Auth routes
	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.Register(db)).Methods("POST")
	authRouter.HandleFunc("/login", handlers.Login(db)).Methods("POST")

	// Protected routes
	notesRouter := router.PathPrefix("/api/v1/notes").Subrouter()
	notesRouter.Use(middleware.AuthMiddleware)
	notesRouter.HandleFunc("", handlers.CreateNote(db)).Methods("POST")
	notesRouter.HandleFunc("", handlers.ListNotes(db)).Methods("GET")
	notesRouter.HandleFunc("/{id}", handlers.GetNote(db)).Methods("GET")
	notesRouter.HandleFunc("/{id}", handlers.UpdateNote(db)).Methods("PUT")
	notesRouter.HandleFunc("/{id}", handlers.DeleteNote(db)).Methods("DELETE")

	// Get port from config
	port := config.GetPort()
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting VicNotes backend server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
