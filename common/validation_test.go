package common_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/darkrockmountain/gomail/common"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{"test@example.com", "test@example.com"},
		{"test@domain_name.com", "test@domain_name.com"},
		{"test@domain-name.com", "test@domain-name.com"},
		{"test@subdomain.example.com", "test@subdomain.example.com"},
		{"test_name@subdomain.example.com", "test_name@subdomain.example.com"},
		{"test.name@subdomain.example.com", "test.name@subdomain.example.com"},
		{"test-name@subdomain.example.com", "test-name@subdomain.example.com"},
		{"  test@example.com  ", "test@example.com"},
		{"invalid-email", ""},
		{"test@.com", ""},
		{"@example.com", ""},
		{"test@com", ""},
		{"test@com.", ""},
		{"test@sub.example.com", "test@sub.example.com"},
		{"test+alias@example.com", "test+alias@example.com"},
		{"test.email@example.com", "test.email@example.com"},
		{"test-email@example.com", "test-email@example.com"},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			result := common.ValidateEmail(test.email)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestValidateEmailSlice(t *testing.T) {
	tests := []struct {
		emails   []string
		expected []string
	}{
		{[]string{"test@example.com"}, []string{"test@example.com"}},
		{[]string{"test@example.com", "invalid-email"}, []string{"test@example.com"}},
		{[]string{" test@example.com ", "test2@example.com"}, []string{"test@example.com", "test2@example.com"}},
		{[]string{"invalid-email", "@example.com"}, []string{}},
		{[]string{"test@example.com", "test2@sub.example.com"}, []string{"test@example.com", "test2@sub.example.com"}},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.emails, ","), func(t *testing.T) {
			result := common.ValidateEmailSlice(test.emails)
			assert.Equal(t, test.expected, result)
		})
	}
}

func ExampleValidateEmail() {
	email := "test@example.com"
	result := common.ValidateEmail(email)
	fmt.Println(result)
	// Output: test@example.com
}

func ExampleValidateEmail_not() {
	email := "test@com"
	result := common.ValidateEmail(email)
	fmt.Println(result)
	// Output:
}

func ExampleValidateEmail_trim() {
	email := "  test@example.com  "
	result := common.ValidateEmail(email)
	fmt.Println(result)
	// Output: test@example.com
}

func ExampleValidateEmailSlice() {
	emails := []string{"test@example.com", "test@domain_name.com"}
	result := common.ValidateEmailSlice(emails)
	fmt.Println(result)
	// Output: [test@example.com test@domain_name.com]
}
func ExampleValidateEmailSlice_partial() {
	emails := []string{"test@example.com", "test@com"}
	result := common.ValidateEmailSlice(emails)
	fmt.Println(result)
	// Output: [test@example.com]
}
