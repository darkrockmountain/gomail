package common

import (
	"html"
	"mime"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// sanitizeInput sanitizes a string to prevent HTML and script injection.
func sanitizeInput(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

// sanitizeHtmlInput sanitizes a string to prevent HTML and script injection but conserving the safe html tags.
func sanitizeHtmlInput(input string) string {
	return bluemonday.UGCPolicy().Sanitize(input)
}

// IsHTML checks if a string contains HTML tags.
// Parameters:
// - str: The string to check.
//
// Returns:
//   - bool: True if the string contains HTML tags, otherwise false.
func IsHTML(str string) bool {
	htmlRegex := regexp.MustCompile(`(?i)<\/?[a-z][\s\S]*>`)
	return htmlRegex.MatchString(str)
}

// GetMimeType returns the MIME type based on the file extension.
// This function takes a filename, extracts its extension, and returns the corresponding MIME type.
//
// Parameters:
// - filename: A string containing the name of the file whose MIME type is to be determined.
//
// Returns:
// - string: The MIME type corresponding to the file extension.
//
// Example:
//
//	mimeType := GetMimeType("document.pdf")
//	fmt.Println(mimeType) // Output: application/pdf
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	return mime.TypeByExtension(ext)
}

// StrPtr takes a string value and returns a pointer to that string.
// This function is useful when you need to work with string pointers, such as in
// scenarios where you need to pass a string by reference or handle optional string fields.
//
// Parameters:
//   - str (string): The input string value that you want to convert to a pointer.
//
// Returns:
//   - *string: A pointer to the input string value.
//
// Example usage:
//
//	name := "John Doe"
//	namePtr := StrPtr(name)
//	fmt.Println(namePtr)  // Output: memory address of the string
//	fmt.Println(*namePtr) // Output: "John Doe"
//
// Detailed explanation:
// The StrPtr function creates a pointer to the given string `str`.
// This can be particularly useful in the following scenarios:
//  1. Passing strings by reference to functions, which can help avoid copying large strings.
//  2. Working with data structures that use pointers to represent optional fields or nullable strings.
//  3. Interfacing with APIs or libraries that require or return string pointers.
//
// By using this function, you can easily obtain a pointer to a string and utilize it in contexts
// where pointers are needed, thus enhancing flexibility and efficiency in your Go programs.
func StrPtr(str string) *string {
	return &str
}
