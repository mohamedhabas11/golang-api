package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordValidationConfig holds the configuration for validating passwords
type PasswordValidationConfig struct {
	MinLength  int
	Validators []PasswordValidator // Additional validators like strength check, etc.
}

// PasswordValidator is an interface for custom password validation strategies
type PasswordValidator interface {
	Validate(password string) error
}

// NewPasswordValidationConfig creates a new PasswordValidationConfig with default values
func NewPasswordValidationConfig(minLength ...int) *PasswordValidationConfig {
	if len(minLength) == 0 {
		minLength = append(minLength, 8) // Default min length
	}
	return &PasswordValidationConfig{
		MinLength: minLength[0],
	}
}

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a hashed password with the plaintext password
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
