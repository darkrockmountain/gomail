package sanitizer

import (
	"html"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSanitizerFunc_Sanitize tests the Sanitize method of SanitizerFunc.
func TestSanitizerFunc_Sanitize(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         string
		sanitizeFunc SanitizerFunc
	}{
		{
			name:  "TrimSpaces",
			input: "  some text  ",
			want:  "some text",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return strings.TrimSpace(message)
			}),
		},
		{
			name:  "ToUpper",
			input: "some text",
			want:  "SOME TEXT",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return strings.ToUpper(message)
			}),
		},
		{
			name:  "ToLower",
			input: "SOME TEXT",
			want:  "some text",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return strings.ToLower(message)
			}),
		},
		{
			name:  "EmptyString",
			input: "",
			want:  "",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return message
			}),
		},
		{
			name:  "WhitespaceOnly",
			input: "     ",
			want:  "",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return strings.TrimSpace(message)
			}),
		},
		{
			name:  "NoSanitization",
			input: "no change",
			want:  "no change",
			sanitizeFunc: SanitizerFunc(func(message string) string {
				return message
			}),
		},
		{
			name:  "SpecialCharacters",
			input: "<script>alert('xss')</script>",
			want:  "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			sanitizeFunc: SanitizerFunc(func(message string) string {
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
