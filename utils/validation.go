package utils

import (
	"net/mail"
)

// ValidateEmail checks email format is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	// if no error, email is valid
	return err == nil
}

// ValidatePassword check if password matches security criteria
func ValidatePassword(password string) bool {
	// TODO: Adjust this for minimal password policy
	return len(password) >= 8
}
