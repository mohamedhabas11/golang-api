package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeyLength = 32

// GenerateJWT generates a JWT token with a given expiration or defaults to 24 hours
func GenerateJWT(userID uint, email string, expiration time.Duration) (string, error) {

	// Retrieve the secret key from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	// Validate the strength of the JWT secret key using ValidatePassword
	jwtSecretConfig := NewPasswordValidationConfig(minSecretKeyLength)
	if err := ValidatePassword(jwtSecret, jwtSecretConfig); err != nil {
		return "", fmt.Errorf("JWT_SECRET is too weak; must be at least %d characters long", minSecretKeyLength)
	}

	// Generate JWT token with claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiration).Unix(), // Set the expiration time
	}

	// Create a new JWT token with the claims and the signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret and return it
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return signedToken, nil
}
