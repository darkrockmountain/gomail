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
		maxAttachmentSize: DefaultMaxAttachmentSize,
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

func TestNewEmailMessage(t *testing.T) {
	t.Run("create plain text email", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		subject := "Subject"
		body := "Email body"
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, body, email.Text)
		assert.Equal(t, "", email.HTML)
	})

	t.Run("create HTML email", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		subject := "Subject"
		body := "<p>Email body</p>"
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, "", email.Text)
		assert.Equal(t, body, email.HTML)
	})
}

func TestNewFullEmailMessage(t *testing.T) {
	t.Run("create full email message", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		cc := []string{"cc@example.com"}
		bcc := []string{"bcc@example.com"}
		replyTo := "replyto@example.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachments := []Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, cc, email.CC)
		assert.Equal(t, bcc, email.BCC)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, attachments, email.Attachments)
	})
}

func TestEmailMessageSetters(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetFrom", func(t *testing.T) {
		expected := "sender@example.com"
		email.SetFrom(expected)
		assert.Equal(t, expected, email.From)
	})

	t.Run("SetSubject", func(t *testing.T) {
		expected := "Subject"
		email.SetSubject(expected)
		assert.Equal(t, expected, email.Subject)
	})

	t.Run("SetTo", func(t *testing.T) {
		expected := []string{"recipient@example.com"}
		email.SetTo(expected)
		assert.Equal(t, expected, email.To)
	})

	t.Run("SetCC", func(t *testing.T) {
		expected := []string{"cc@example.com"}
		email.SetCC(expected)
		assert.Equal(t, expected, email.CC)
	})

	t.Run("SetBCC", func(t *testing.T) {
		expected := []string{"bcc@example.com"}
		email.SetBCC(expected)
		assert.Equal(t, expected, email.BCC)
	})

	t.Run("SetReplyTo", func(t *testing.T) {
		expected := "replyto@example.com"
		email.SetReplyTo(expected)
		assert.Equal(t, expected, email.ReplyTo)
	})

	t.Run("SetText", func(t *testing.T) {
		expected := "Text body"
		email.SetText(expected)
		assert.Equal(t, expected, email.Text)
	})

	t.Run("SetHTML", func(t *testing.T) {
		expected := "<p>HTML body</p>"
		email.SetHTML(expected)
		assert.Equal(t, expected, email.HTML)
	})

	t.Run("SetAttachments", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email.SetAttachments([]Attachment{attachment})
		assert.Contains(t, email.Attachments, attachment)
		assert.EqualValues(t, email.Attachments, []Attachment{attachment})
	})

	t.Run("AddAttachment", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email.AddAttachment(attachment)
		assert.Contains(t, email.Attachments, attachment)
	})

	t.Run("AddToRecipient", func(t *testing.T) {
		recipient := "newrecipient@example.com"
		email.AddToRecipient(recipient)
		assert.Contains(t, email.To, recipient)
	})

	t.Run("AddCCRecipient", func(t *testing.T) {
		recipient := "newcc@example.com"
		email.AddCCRecipient(recipient)
		assert.Contains(t, email.CC, recipient)
	})

	t.Run("AddBCCRecipient", func(t *testing.T) {
		recipient := "newbcc@example.com"
		email.AddBCCRecipient(recipient)
		assert.Contains(t, email.BCC, recipient)
	})
}

func TestAddsEmailMessageToNils(t *testing.T) {
	t.Run("create full email message", func(t *testing.T) {
		from := "sender@example.com"
		to := "recipient@example.com"
		cc := "cc@example.com"
		bcc := "bcc@example.com"
		replyTo := "replyto@example.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email := NewFullEmailMessage(from, nil, subject, nil, nil, replyTo, text, html, nil)

		email.AddToRecipient(to)
		email.AddCCRecipient(cc)
		email.AddBCCRecipient(bcc)
		email.AddAttachment(attachment)

		assert.Equal(t, from, email.From)
		assert.Equal(t, []string{to}, email.To)
		assert.Equal(t, []string{cc}, email.CC)
		assert.Equal(t, []string{bcc}, email.BCC)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, []Attachment{attachment}, email.Attachments)
	})
}

func TestAttachmentGetters(t *testing.T) {
	t.Run("GetFilename", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt"}
		assert.Equal(t, "test.txt", attachment.GetFilename())
		assert.Equal(t, "nil_attachment", (*Attachment)(nil).GetFilename())
	})

	t.Run("GetBase64Content", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("hello")}
		expected := []byte(base64.StdEncoding.EncodeToString([]byte("hello")))
		assert.Equal(t, expected, attachment.GetBase64Content())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetBase64Content())
	})

	t.Run("GetRawContent", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("hello")}
		expected := []byte("hello")
		assert.Equal(t, expected, attachment.GetRawContent())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetRawContent())
	})
}

func TestEmailHTMLBodySanitizers(t *testing.T) {
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

func TestSetMaxAttachmentSize(t *testing.T) {
	email := &EmailMessage{}
	t.Run("SetMaxAttachmentSize", func(t *testing.T) {
		expected := 10 * 1024 * 1024 // 10 MB
		email.SetMaxAttachmentSize(expected)
		assert.Equal(t, expected, email.maxAttachmentSize)
	})
}

func TestGetAttachmentsWithMaxSize(t *testing.T) {
	email := &EmailMessage{
		Attachments: []Attachment{
			{Filename: "small.txt", Content: []byte("small content")},
			{Filename: "large.txt", Content: make([]byte, 30*1024*1024)}, // 30 MB
		},
		maxAttachmentSize: 25 * 1024 * 1024, // 25 MB
	}

	t.Run("GetAttachments with size limit", func(t *testing.T) {
		expected := []Attachment{
			{Filename: "small.txt", Content: []byte("small content")},
		}
		result := email.GetAttachments()
		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments with no size limit", func(t *testing.T) {
		email.SetMaxAttachmentSize(-1)
		expected := email.Attachments
		result := email.GetAttachments()
		assert.Equal(t, expected, result)
	})
}
