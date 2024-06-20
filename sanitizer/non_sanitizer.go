package sanitizer

// nonSanitizer is an implementation of the Sanitizer interface that performs no sanitization.
// It returns the input text unchanged. This can be useful in scenarios where no sanitization is desired.
type nonSanitizer struct{}

// Sanitize returns the input text without any modifications.
// This method complies with the Sanitizer interface.
//
// Parameters:
// - text: The content to be "sanitized".
//
// Returns:
// - string: The input text, unchanged.
func (s *nonSanitizer) Sanitize(text string) string {
	return text
}

// nonSanitizerInstance is the singleton instance of nonSanitizer.
var nonSanitizerInstance = &nonSanitizer{}

// NonSanitizer returns the singleton instance of nonSanitizer.
//
// nonSanitizer is an implementation of the Sanitizer interface that performs no sanitization.
// It returns the input text unchanged. This can be useful in scenarios where no sanitization is desired.
//
// WARNING: Using nonSanitizer means that the input content will not be sanitized at all.
// Ensure that this is acceptable for your use case, as it may introduce security risks, such as
// injection attacks, if user input is not properly handled.
//
// Example usage:
//
//	ns := sanitizer.NonSanitizer()
//	unsanitized := ns.Sanitize("<script>alert('xss')</script>")
//	// unsanitized will be "<script>alert('xss')</script>"
func NonSanitizer() Sanitizer {
	return nonSanitizerInstance
}
