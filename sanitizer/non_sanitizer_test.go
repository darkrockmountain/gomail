package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonSanitizer(t *testing.T) {
	sanitizer := NonSanitizer()

	t.Run("return input as is", func(t *testing.T) {
		input := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		expected := input
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("plain text", func(t *testing.T) {
		input := "plain text"
		expected := "plain text"
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})
}
