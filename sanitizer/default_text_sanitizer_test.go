package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTextSanitizer(t *testing.T) {
	sanitizer := DefaultTextSanitizer()

	t.Run("escape special characters", func(t *testing.T) {
		input := `Hello <world> & "everyone"`
		expected := `Hello &lt;world&gt; &amp; &#34;everyone&#34;`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("string with spaces", func(t *testing.T) {
		input := "   "
		expected := ""
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("string with no special characters", func(t *testing.T) {
		input := "plain text"
		expected := "plain text"
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})

	t.Run("sanitize documentation content", func(t *testing.T) {
		input := " <script>alert('xss')</script> "
		expected := `&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;`
		result := sanitizer.Sanitize(input)
		assert.Equal(t, expected, result)
	})
}
