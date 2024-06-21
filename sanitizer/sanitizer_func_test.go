package sanitizer_test

import (
	"fmt"
	"html"
	"strings"
	"testing"

	"github.com/darkrockmountain/gomail/sanitizer"
	"github.com/stretchr/testify/assert"
)

// TestSanitizerFunc_Sanitize tests the Sanitize method of SanitizerFunc.
func TestSanitizerFunc_Sanitize(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         string
		sanitizeFunc sanitizer.SanitizerFunc
	}{
		{
			name:  "TrimSpaces",
			input: "  some text  ",
			want:  "some text",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return strings.TrimSpace(message)
			}),
		},
		{
			name:  "ToUpper",
			input: "some text",
			want:  "SOME TEXT",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return strings.ToUpper(message)
			}),
		},
		{
			name:  "ToLower",
			input: "SOME TEXT",
			want:  "some text",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return strings.ToLower(message)
			}),
		},
		{
			name:  "EmptyString",
			input: "",
			want:  "",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return message
			}),
		},
		{
			name:  "WhitespaceOnly",
			input: "     ",
			want:  "",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return strings.TrimSpace(message)
			}),
		},
		{
			name:  "NoSanitization",
			input: "no change",
			want:  "no change",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return message
			}),
		},
		{
			name:  "SpecialCharacters",
			input: "<script>alert('xss')</script>",
			want:  "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			sanitizeFunc: sanitizer.SanitizerFunc(func(message string) string {
				return html.EscapeString(message)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sanitizeFunc.Sanitize(tt.input))
		})
	}
}

func ExampleSanitizerFunc() {
	sanitizeFunc := sanitizer.SanitizerFunc(func(message string) string {
		// Implement your custom sanitizer logic
		return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(message)), " ", "_")
	})
	sanitizedMessage := sanitizeFunc.Sanitize("  some text  ")
	fmt.Println(sanitizedMessage)
	// Output: some_text
}
