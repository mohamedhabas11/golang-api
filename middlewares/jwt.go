package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mohamedhabas11/golang-api/utils"
)

const minSecretKeyLength = 32

// GenerteJWTSecret generates a secure random secret key and sets it as an environment variable.
func GenerteJWTSecret() (string, error) {
	// Generate a secure random byte slice of the required length
	secretBytes := make([]byte, minSecretKeyLength)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt secret key: %v", err)
	}

	// Encode the byte slice to a base64 string
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	// Set the secret key as the expected enviorment variable
	err = os.Setenv("JWT_SECRET", secret)
	if err != nil {
		return "", fmt.Errorf("failed to set JWT_SECRET environment variable: %v", err)
	}

	// Log that a new JWT secret was created
	log.Printf("%v A new JWT_SECRET was generated and set.", time.Now().Format(time.RFC3339))
	return secret, nil
}

// GenerateJWT generates a JWT token with a given expiration or defaults to 24 hours
func GenerateJWT(userID uint, email string, expiration time.Duration) (string, error) {

	// Retrieve the secret key from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Log a warning that the secret is not set
		log.Printf("%v Warning: JWT_SECRET env variable is not set. Generating a new secret. ", time.Now().Format(time.RFC3339))

		// Generate and set new secret
		var err error
		jwtSecret, err = GenerteJWTSecret()
		if err != nil {
			return "", fmt.Errorf("failed to generate a new JWT secret: %v", err)
		}
	}

	// Validate the strength of the JWT secret key using ValidatePassword
	jwtSecretConfig := utils.NewPasswordValidationConfig(minSecretKeyLength)
	if err := utils.ValidatePassword(jwtSecret, jwtSecretConfig); err != nil {
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
