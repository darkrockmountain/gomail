package sanitizer

import "github.com/microcosm-cc/bluemonday"

// defaultHtmlSanitizer provides a basic implementation of the Sanitizer interface for HTML content.
// It uses the bluemonday library to sanitize HTML content, ensuring that only user-generated content
// (UGC) is allowed. This helps prevent injection attacks by removing potentially dangerous tags and attributes.
type defaultHtmlSanitizer struct{}

// Sanitize sanitizes HTML content by removing potentially dangerous tags and attributes.
// It uses the bluemonday UGCPolicy to allow only safe HTML for user-generated content.
//
// Parameters:
//   - htmlContent: The HTML content to be sanitized.
//
// Returns:
//   - string: The sanitized HTML content.
func (s *defaultHtmlSanitizer) Sanitize(htmlContent string) string {
	// Use bluemonday's UGCPolicy for robust HTML sanitization.
	return bluemonday.UGCPolicy().Sanitize(htmlContent)
}

// defaultHtmlSanitizerInstance is the singleton instance of defaultHtmlSanitizer.
var defaultHtmlSanitizerInstance = &defaultHtmlSanitizer{}

// DefaultHtmlSanitizer returns the singleton instance of defaultHtmlSanitizer.
//
// defaultHtmlSanitizer provides a basic implementation of the Sanitizer interface for HTML content.
// It uses the bluemonday library to sanitize HTML content, ensuring that only user-generated content
// (UGC) is allowed. This helps prevent injection attacks by removing potentially dangerous tags and attributes.
//
// Example:
//
//	hs := sanitizer.DefaultHtmlSanitizer()
//	sanitized := hs.Sanitize("<script>alert('xss')</script><b>Bold</b>")
//	// sanitized will be "<b>Bold</b>"
func DefaultHtmlSanitizer() Sanitizer {
	return defaultHtmlSanitizerInstance
}
