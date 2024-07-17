package sanitizer

// SanitizerFunc type is an adapter that allows the use of
// ordinary functions as Sanitizers. If f is a function
// with the appropriate signature, SanitizerFunc(f) is a
// Sanitizer that calls f.
//
// Example:
//
//	sanitizeFunc := sanitizer.SanitizerFunc(func(message string) string {
//	    // Implement your custom sanitizer logic
//	    return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(message)), " ", "_")
//	})
//	sanitizedMessage := sanitizeFunc.Sanitize("  some text  ")
//	// sanitizedMessage will be "some_text"
type SanitizerFunc func(message string) string

// Sanitize calls the function f with the given message.
//
// Parameters:
//   - message: The content to be sanitized.
//
// Returns:
//   - string: The sanitized content.
func (f SanitizerFunc) Sanitize(message string) string {
	return f(message)
}
