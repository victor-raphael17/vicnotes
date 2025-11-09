package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vicnotes/backend/config"
)

// Claims represents JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
}

// GenerateToken generates a JWT token
func GenerateToken(userID int, email string) (string, error) {
	secret := config.GetJWTSecret()

	// Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Payload
	payload := Claims{
		UserID: userID,
		Email:  email,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Signature
	message := headerB64 + "." + payloadB64
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := message + "." + signature
	return token, nil
}

// VerifyToken verifies and parses a JWT token
func VerifyToken(tokenString string) (*Claims, error) {
	secret := config.GetJWTSecret()

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Verify signature
	message := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if parts[2] != expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiration
	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
