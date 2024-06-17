package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailHTMLBodySanitzers(t *testing.T) {

	html1 := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`

	t.Run("remove potential XSS attack", func(t *testing.T) {
		expected := `<div>XSS</div>`
		result := sanitizeHtmlInput(html1)
		assert.Equal(t, expected, result)
	})

	html2 := `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`

	t.Run("on methods not allowed", func(t *testing.T) {
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := sanitizeHtmlInput(html2)
		assert.Equal(t, expected, result)
	})

	html3 := `<p href="http://www.google.com">Google</p>`

	t.Run("<p> can't have href", func(t *testing.T) {
		expected := `<p>Google</p>`
		result := sanitizeHtmlInput(html3)
		assert.Equal(t, expected, result)
	})

}

func TestEmailHTMLBodySanitizers(t *testing.T) {
	html1 := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`

	t.Run("remove potential XSS attack", func(t *testing.T) {
		expected := `<div>XSS</div>`
		result := sanitizeHtmlInput(html1)
		assert.Equal(t, expected, result)
	})

	html2 := `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`

	t.Run("on methods not allowed", func(t *testing.T) {
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := sanitizeHtmlInput(html2)
		assert.Equal(t, expected, result)
	})

	html3 := `<p href="http://www.google.com">Google</p>`

	t.Run("<p> can't have href", func(t *testing.T) {
		expected := `<p>Google</p>`
		result := sanitizeHtmlInput(html3)
		assert.Equal(t, expected, result)
	})
}

func TestStrPtr(t *testing.T) {

	str := "String to test for pointer"
	ptrStr := StrPtr(str)
	assert.Equal(t, ptrStr, &str)
	assert.EqualValues(t, ptrStr, &str)
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"document.pdf", "application/pdf"},
		{"image.png", "image/png"},
		{"archive.zip", "application/zip"},
		{"unknownfile.unknown", ""},
		{"text.txt", "text/plain; charset=utf-8"},
		{"no_extension", ""},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := GetMimeType(test.filename)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetMimeTypeEdgeCases(t *testing.T) {
	t.Run("unknown extension", func(t *testing.T) {
		filename := "file.unknownext"
		expected := ""
		result := GetMimeType(filename)
		assert.Equal(t, expected, result)
	})

	t.Run("empty filename", func(t *testing.T) {
		filename := ""
		expected := ""
		result := GetMimeType(filename)
		assert.Equal(t, expected, result)
	})
}
