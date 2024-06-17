package common

import (
	"regexp"
	"strings"
)

// regex for validating email addresses
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9._\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail trims the email and checks if it is a valid email address.
// Returns the trimmed email if valid, otherwise returns an empty string.
func ValidateEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	if !emailRegex.MatchString(trimmed) {
		return ""
	}
	return trimmed
}

// ValidateEmailSlice trims and validates each email in the slice.
// Returns a slice of trimmed valid emails, excluding any invalid emails.
func ValidateEmailSlice(emails []string) []string {
	validEmails := []string{}
	for _, email := range emails {
		if validEmail := ValidateEmail(email); validEmail != "" {
			validEmails = append(validEmails, validEmail)
		}
	}
	return validEmails
}
