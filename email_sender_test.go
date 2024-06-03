package gomail

import (
	"encoding/base64"
	"strings"
	"testing"

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
			result := ValidateEmail(test.email)
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
			result := ValidateEmailSlice(test.emails)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetBase64Content(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   []byte
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, []byte(base64.StdEncoding.EncodeToString([]byte("hello")))},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, []byte{}},
		{&Attachment{Filename: "empty.txt", Content: nil}, []byte{}},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetBase64Content()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetBase64StringContent(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   string
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, base64.StdEncoding.EncodeToString([]byte("hello"))},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, ""},
		{nil, ""},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetBase64StringContent()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetRawContent(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   []byte
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, []byte("hello")},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, []byte{}},
		{&Attachment{Filename: "empty.txt", Content: nil}, []byte{}},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetRawContent()
			assert.Equal(t, test.expected, result)
		})
	}
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

func TestEmailMessageGetters(t *testing.T) {
	message := EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient1@example.com", "recipient2@example.com"},
		CC:      []string{"cc1@example.com", "cc2@example.com"},
		BCC:     []string{"bcc1@example.com", "bcc2@example.com"},
		ReplyTo: "replyto@example.com",
		Subject: "Test Subject",
		Text:    "Test Text",
		HTML:    "<h1>Test HTML</h1>",
		Attachments: []Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		},
	}

	t.Run("GetFrom", func(t *testing.T) {
		expected := "sender@example.com"
		result := message.GetFrom()
		assert.Equal(t, expected, result)
	})

	t.Run("GetReplyTo", func(t *testing.T) {
		expected := "replyto@example.com"
		result := message.GetReplyTo()
		assert.Equal(t, expected, result)
	})

	t.Run("GetTo", func(t *testing.T) {
		expected := []string{"recipient1@example.com", "recipient2@example.com"}
		result := message.GetTo()
		assert.Equal(t, expected, result)
	})

	t.Run("GetCC", func(t *testing.T) {
		expected := []string{"cc1@example.com", "cc2@example.com"}
		result := message.GetCC()
		assert.Equal(t, expected, result)
	})

	t.Run("GetBCC", func(t *testing.T) {
		expected := []string{"bcc1@example.com", "bcc2@example.com"}
		result := message.GetBCC()
		assert.Equal(t, expected, result)
	})

	t.Run("GetSubject", func(t *testing.T) {
		expected := "Test Subject"
		result := message.GetSubject()
		assert.Equal(t, expected, result)
	})

	t.Run("GetText", func(t *testing.T) {
		expected := "Test Text"
		result := message.GetText()
		assert.Equal(t, expected, result)
	})

	t.Run("GetHTML", func(t *testing.T) {
		expected := "<h1>Test HTML</h1>"
		result := message.GetHTML()
		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments", func(t *testing.T) {
		expected := []Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		}
		result := message.GetAttachments()
		assert.Equal(t, expected, result)
	})
}

func TestNilEmailMessageGetters(t *testing.T) {
	var message *EmailMessage

	t.Run("GetFrom", func(t *testing.T) {
		result := message.GetFrom()
		assert.Equal(t, "", result)
	})

	t.Run("GetReplyTo", func(t *testing.T) {
		result := message.GetReplyTo()
		assert.Equal(t, "", result)
	})

	t.Run("GetTo", func(t *testing.T) {
		result := message.GetTo()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetCC", func(t *testing.T) {
		result := message.GetCC()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetBCC", func(t *testing.T) {
		result := message.GetBCC()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetSubject", func(t *testing.T) {
		result := message.GetSubject()
		assert.Equal(t, "", result)
	})

	t.Run("GetText", func(t *testing.T) {
		result := message.GetText()
		assert.Equal(t, "", result)
	})

	t.Run("GetHTML", func(t *testing.T) {
		result := message.GetHTML()
		assert.Equal(t, "", result)
	})

	t.Run("GetAttachments", func(t *testing.T) {
		result := message.GetAttachments()
		assert.Equal(t, []Attachment{}, result)
	})
}

func TestEmailHTMLBodySanitzers(t *testing.T) {
	message := EmailMessage{
		HTML: `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`,
	}

	t.Run("remove potential XSS attack", func(t *testing.T) {
		expected := `<div>XSS</div>`
		result := message.GetHTML()
		assert.Equal(t, expected, result)
	})

	message2 := EmailMessage{
		HTML: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	}

	t.Run("on methods not allowed", func(t *testing.T) {
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := message2.GetHTML()
		assert.Equal(t, expected, result)
	})

	message3 := EmailMessage{
		HTML: `<p href="http://www.google.com">Google</p>`,
	}
	t.Run("<p> can't have href", func(t *testing.T) {
		expected := `<p>Google</p>`
		result := message3.GetHTML()
		assert.Equal(t, expected, result)
	})

}
