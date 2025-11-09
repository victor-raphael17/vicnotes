package config

import (
	"os"
)

// GetPort returns the server port from environment or default
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// GetDatabaseURL returns the database connection string
func GetDatabaseURL() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + db + "?sslmode=disable"
}

// GetJWTSecret returns the JWT secret key
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return secret
}
