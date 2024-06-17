package common

import (
	"testing"
)

func TestIsHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<html><body>Hello</body></html>", true},
		{"Just a plain text", false},
		{"<div>HTML content</div>", true},
		{"Plain text with <html> tag", true},
	}

	for _, test := range tests {
		result := IsHTML(test.input)
		if result != test.expected {
			t.Errorf("isHTML(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
