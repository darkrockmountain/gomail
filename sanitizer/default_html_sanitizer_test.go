package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHtmlSanitizer(t *testing.T) {
	sanitizer := DefaultHtmlSanitizer()

	t.Run("remove potential XSS attack", func(t *testing.T) {
		input := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		expected := `<div>XSS</div>`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("on methods not allowed", func(t *testing.T) {
		input := `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("<p> can't have href", func(t *testing.T) {
		input := `<p href="http://www.google.com">Google</p>`
		expected := `<p>Google</p>`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("allow safe HTML tags", func(t *testing.T) {
		input := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		expected := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("sanitize mixed content", func(t *testing.T) {
		input := `<div>Hello <script>alert("xss")</script> World</div>`
		expected := `<div>Hello  World</div>`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})
}
