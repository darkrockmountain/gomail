// Package sanitizer provides interfaces and default implementations for sanitizing email content.
//
// # Overview
//
// The sanitizer package defines an interface for content sanitization and provides default implementations
// for sanitizing plain text and HTML content. This allows for flexible and customizable content sanitization
// in the gomail project.
//
// # Usage
//
// To use the sanitizer package, you can either use the provided default implementations or create your
// own custom implementations of the Sanitizer interface and set them in the EmailMessage struct.
//
// Example:
//
//	import (
//	    "github.com/darkrockmountain/gomail/sanitizer"
//	    "github.com/darkrockmountain/gomail/common"
//	)
//
//	func main() {
//	    email := common.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Subject", "<p>HTML content</p>")
//	    customSanitizer := &sanitizer.CustomSanitizer{}
//	    email.SetCustomTextSanitizer(customSanitizer)
//	    email.SetCustomHtmlSanitizer(customSanitizer)
//	}
//
// The sanitizer package is designed to be used in conjunction with the gomail project for flexible email content sanitization.
package sanitizer

// Sanitizer defines a method for sanitizing content.
// Implement this interface to provide custom sanitization logic for email content.
type Sanitizer interface {
	// Sanitize sanitizes the provided content.
	// Parameters:
	// - content: The content to be sanitized.
	// Returns:
	// - string: The sanitized content.
	Sanitize(content string) string
}
