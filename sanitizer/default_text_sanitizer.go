package sanitizer

import (
	"html"
	"strings"
)

// defaultTextSanitizer provides a basic implementation of the Sanitizer interface for plain text content.
type defaultTextSanitizer struct{}

// Sanitize sanitizes plain text content by escaping special characters and trimming whitespace.
// This method complies with the Sanitizer interface.
//
// Parameters:
//   - text: The plain text content to be sanitized.
//
// Returns:
//   - string: The sanitized text content.
func (s *defaultTextSanitizer) Sanitize(text string) string {
	return html.EscapeString(strings.TrimSpace(text))
}

// defaultTextSanitizerInstance is the singleton instance of defaultTextSanitizer.
var defaultTextSanitizerInstance = &defaultTextSanitizer{}

// DefaultTextSanitizer returns the singleton instance of defaultTextSanitizer.
//
// defaultTextSanitizer provides a basic implementation of the Sanitizer interface for plain text content.
// It sanitizes plain text content by escaping special characters and trimming whitespace.
//
// Example:
//
//	ds := sanitizer.DefaultTextSanitizer()
//	sanitized := ds.Sanitize(" <script>alert('xss')</script> ")
//	// sanitized will be "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
func DefaultTextSanitizer() Sanitizer {
	return defaultTextSanitizerInstance
}
