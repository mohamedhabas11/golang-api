package utils

import (
	"errors"
	"net/mail"
	"strings"
)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePassword checks if the password meets the minimum length and passes additional validations
func ValidatePassword(password string, config *PasswordValidationConfig) error {
	if config == nil {
		config = NewPasswordValidationConfig() // Use default config
	}

	// Validate minimum length
	if len(password) < config.MinLength {
		return errors.New("password is too short")
	}

	// Apply additional validators if any
	for _, validator := range config.Validators {
		if err := validator.Validate(password); err != nil {
			return err
		}
	}
	return nil
}

// AddPasswordStrengthValidator checks if the password contains at least one uppercase letter and one number
type AddPasswordStrengthValidator struct{}

// Validate checks if the password has at least one uppercase letter and one number
func (v *AddPasswordStrengthValidator) Validate(password string) error {
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("password must contain at least one number")
	}
	return nil
}

// AddPasswordSpecialCharValidator checks if the password contains at least one special character
type AddPasswordSpecialCharValidator struct{}

// Validate checks if the password has at least one special character
func (v *AddPasswordSpecialCharValidator) Validate(password string) error {
	if !strings.ContainsAny(password, "!@#$%^&*()_+[]{}|;:,.<>?") {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
