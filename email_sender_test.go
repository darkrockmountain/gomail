package gomail

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
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

func TestEmailMessageGetters(t *testing.T) {
	message := EmailMessage{
		from:    "sender@example.com",
		to:      []string{"recipient1@example.com", "recipient2@example.com"},
		cc:      []string{"cc1@example.com", "cc2@example.com"},
		bcc:     []string{"bcc1@example.com", "bcc2@example.com"},
		replyTo: "replyto@example.com",
		subject: "Test Subject",
		text:    "Test Text",
		html:    "<h1>Test HTML</h1>",
		attachments: []Attachment{
			{filename: "test.txt", content: []byte("test content")},
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
			{filename: "test.txt", content: []byte("test content")},
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
		html: `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`,
	}

	t.Run("remove potential XSS attack", func(t *testing.T) {
		expected := `<div>XSS</div>`
		result := message.GetHTML()
		assert.Equal(t, expected, result)
	})

	message2 := EmailMessage{
		html: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	}

	t.Run("on methods not allowed", func(t *testing.T) {
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := message2.GetHTML()
		assert.Equal(t, expected, result)
	})

	message3 := EmailMessage{
		html: `<p href="http://www.google.com">Google</p>`,
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

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, body, email.text)
		assert.Equal(t, "", email.html)
	})

	t.Run("create HTML email", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		subject := "Subject"
		body := "<p>Email body</p>"
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, "", email.text)
		assert.Equal(t, body, email.html)
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
			{filename: "test.txt", content: []byte("test content")},
		}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, cc, email.cc)
		assert.Equal(t, bcc, email.bcc)
		assert.Equal(t, replyTo, email.replyTo)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, text, email.text)
		assert.Equal(t, html, email.html)
		assert.Equal(t, attachments, email.attachments)
	})
}

func TestEmailMessageSetters(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetFrom", func(t *testing.T) {
		expected := "sender@example.com"
		email.SetFrom(expected)
		assert.Equal(t, expected, email.from)
	})

	t.Run("SetSubject", func(t *testing.T) {
		expected := "Subject"
		email.SetSubject(expected)
		assert.Equal(t, expected, email.subject)
	})

	t.Run("SetTo", func(t *testing.T) {
		expected := []string{"recipient@example.com"}
		email.SetTo(expected)
		assert.Equal(t, expected, email.to)
	})

	t.Run("SetCC", func(t *testing.T) {
		expected := []string{"cc@example.com"}
		email.SetCC(expected)
		assert.Equal(t, expected, email.cc)
	})

	t.Run("SetBCC", func(t *testing.T) {
		expected := []string{"bcc@example.com"}
		email.SetBCC(expected)
		assert.Equal(t, expected, email.bcc)
	})

	t.Run("SetReplyTo", func(t *testing.T) {
		expected := "replyto@example.com"
		email.SetReplyTo(expected)
		assert.Equal(t, expected, email.replyTo)
	})

	t.Run("SetText", func(t *testing.T) {
		expected := "Text body"
		email.SetText(expected)
		assert.Equal(t, expected, email.text)
	})

	t.Run("SetHTML", func(t *testing.T) {
		expected := "<p>HTML body</p>"
		email.SetHTML(expected)
		assert.Equal(t, expected, email.html)
	})

	t.Run("SetAttachments", func(t *testing.T) {
		attachment := Attachment{filename: "test.txt", content: []byte("test content")}
		email.SetAttachments([]Attachment{attachment})
		assert.Contains(t, email.attachments, attachment)
		assert.EqualValues(t, email.attachments, []Attachment{attachment})
	})

	t.Run("AddAttachment", func(t *testing.T) {
		attachment := Attachment{filename: "test.txt", content: []byte("test content")}
		email.AddAttachment(attachment)
		assert.Contains(t, email.attachments, attachment)
	})

	t.Run("AddToRecipient", func(t *testing.T) {
		recipient := "newrecipient@example.com"
		email.AddToRecipient(recipient)
		assert.Contains(t, email.to, recipient)
	})

	t.Run("AddCCRecipient", func(t *testing.T) {
		recipient := "newcc@example.com"
		email.AddCCRecipient(recipient)
		assert.Contains(t, email.cc, recipient)
	})

	t.Run("AddBCCRecipient", func(t *testing.T) {
		recipient := "newbcc@example.com"
		email.AddBCCRecipient(recipient)
		assert.Contains(t, email.bcc, recipient)
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
		attachment := Attachment{filename: "test.txt", content: []byte("test content")}
		email := NewFullEmailMessage(from, nil, subject, nil, nil, replyTo, text, html, nil)

		email.AddToRecipient(to)
		email.AddCCRecipient(cc)
		email.AddBCCRecipient(bcc)
		email.AddAttachment(attachment)

		assert.Equal(t, from, email.from)
		assert.Equal(t, []string{to}, email.to)
		assert.Equal(t, []string{cc}, email.cc)
		assert.Equal(t, []string{bcc}, email.bcc)
		assert.Equal(t, replyTo, email.replyTo)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, text, email.text)
		assert.Equal(t, html, email.html)
		assert.Equal(t, []Attachment{attachment}, email.attachments)
	})
}

func TestEmailHTMLBodySanitizers(t *testing.T) {
	message := EmailMessage{
		html: `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`,
	}

	t.Run("remove potential XSS attack", func(t *testing.T) {
		expected := `<div>XSS</div>`
		result := message.GetHTML()
		assert.Equal(t, expected, result)
	})

	message2 := EmailMessage{
		html: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	}

	t.Run("on methods not allowed", func(t *testing.T) {
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := message2.GetHTML()
		assert.Equal(t, expected, result)
	})

	message3 := EmailMessage{
		html: `<p href="http://www.google.com">Google</p>`,
	}
	t.Run("<p> can't have href", func(t *testing.T) {
		expected := `<p>Google</p>`
		result := message3.GetHTML()
		assert.Equal(t, expected, result)
	})
}

func TestMarshalJSONCustom(t *testing.T) {
	t.Run("Marshal EmailMessage with attachments", func(t *testing.T) {
		email := &EmailMessage{
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			cc:      []string{"cc@example.com"},
			bcc:     []string{"bcc@example.com"},
			replyTo: "replyto@example.com",
			subject: "Subject",
			text:    "This is the email content.",
			html:    "<p>This is the email content.</p>",
			attachments: []Attachment{
				{filename: "attachment1.txt", content: []byte("content1")},
			},
			maxAttachmentSize: 1024,
		}
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)
		expected := `{"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>","attachments":[{"filename":"attachment1.txt","content":"Y29udGVudDE="}]}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("Marshal EmailMessage without attachments", func(t *testing.T) {
		email := &EmailMessage{
			from:              "sender@example.com",
			to:                []string{"recipient@example.com"},
			cc:                []string{"cc@example.com"},
			bcc:               []string{"bcc@example.com"},
			replyTo:           "replyto@example.com",
			subject:           "Subject",
			text:              "This is the email content.",
			html:              "<p>This is the email content.</p>",
			maxAttachmentSize: 1024,
		}
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)
		expected := `{"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>"}`
		assert.JSONEq(t, expected, string(jsonData))
	})
}

func TestUnmarshalJSONCustom(t *testing.T) {
	t.Run("Unmarshal EmailMessage with attachments", func(t *testing.T) {
		jsonData := `{
			"from": "sender@example.com",
			"to": ["recipient@example.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>",
			"attachments": [{"filename": "attachment1.txt", "content": "Y29udGVudDE="}]
		}`
		var email EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "sender@example.com", email.from)
		assert.Equal(t, []string{"recipient@example.com"}, email.to)
		assert.Equal(t, []string{"cc@example.com"}, email.cc)
		assert.Equal(t, []string{"bcc@example.com"}, email.bcc)
		assert.Equal(t, "replyto@example.com", email.replyTo)
		assert.Equal(t, "Subject", email.subject)
		assert.Equal(t, "This is the email content.", email.text)
		assert.Equal(t, "<p>This is the email content.</p>", email.html)
		expectedAttachment := Attachment{filename: "attachment1.txt", content: []byte("content1")}
		assert.Equal(t, []Attachment{expectedAttachment}, email.attachments)
	})

	t.Run("Unmarshal EmailMessage without attachments", func(t *testing.T) {
		jsonData := `{
			"from": "sender@example.com",
			"to": ["recipient@example.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>"
		}`
		var email EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "sender@example.com", email.from)
		assert.Equal(t, []string{"recipient@example.com"}, email.to)
		assert.Equal(t, []string{"cc@example.com"}, email.cc)
		assert.Equal(t, []string{"bcc@example.com"}, email.bcc)
		assert.Equal(t, "replyto@example.com", email.replyTo)
		assert.Equal(t, "Subject", email.subject)
		assert.Equal(t, "This is the email content.", email.text)
		assert.Equal(t, "<p>This is the email content.</p>", email.html)
		assert.Nil(t, email.attachments)
	})
}

func TestGetBase64Content(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   []byte
	}{
		{&Attachment{filename: "test.txt", content: []byte("hello")}, []byte(base64.StdEncoding.EncodeToString([]byte("hello")))},
		{&Attachment{filename: "test.txt", content: []byte("")}, []byte{}},
		{&Attachment{filename: "empty.txt", content: nil}, []byte{}},
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
		{&Attachment{filename: "test.txt", content: []byte("hello")}, base64.StdEncoding.EncodeToString([]byte("hello"))},
		{&Attachment{filename: "test.txt", content: []byte("")}, ""},
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
		{&Attachment{filename: "test.txt", content: []byte("hello")}, []byte("hello")},
		{&Attachment{filename: "test.txt", content: []byte("")}, []byte{}},
		{&Attachment{filename: "empty.txt", content: nil}, []byte{}},
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

func TestExtractFilename(t *testing.T) {
	t.Run("extract filename from valid path", func(t *testing.T) {
		filePath := "/path/to/file/document.pdf"
		expected := "document.pdf"
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})

	t.Run("extract filename from empty path", func(t *testing.T) {
		filePath := ""
		expected := ""
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})

	t.Run("extract filename from path with trailing slash", func(t *testing.T) {
		filePath := "/path/to/directory/"
		expected := ""
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})
}

func TestNewAttachmentFromFile(t *testing.T) {

	t.Run("files in test data", func(t *testing.T) {
		testFiles := []struct {
			filePath        string
			expectedName    string
			expectedContent string
		}{
			{
				filepath.Join("testdata", "testfile.txt"),
				"testfile.txt",
				`DarkRockMountain
https://darkrockmountain.com/

we make it possible

DarkRockMountain, your trusted partner for developing, implementing, scaling, and maintaining your solutions from concept to production.`,
			}, {
				filepath.Join("testdata", "testfile.md"),
				"testfile.md",
				`# DarkRockMountain
**[darkrockmountain.com](https://darkrockmountain.com/)**

### we make it possible

DarkRockMountain, your trusted partner for developing, implementing, scaling, and maintaining your solutions from concept to production.`,
			},
		}

		for _, testFile := range testFiles {
			attachment, err := NewAttachmentFromFile(testFile.filePath)
			if err != nil {
				t.Fatalf("NewAttachmentFromFile() error = %v, want nil", err)
			}

			assert.Equal(t, attachment.filename, testFile.expectedName)

			content, err := os.ReadFile(testFile.filePath)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			assert.Equal(t, string(attachment.content), string(content))

		}

	})

	t.Run("file does not exist", func(t *testing.T) {
		filePath := "nonexistentfile.txt"
		attachment, err := NewAttachmentFromFile(filePath)
		assert.NotNil(t, err)
		assert.Nil(t, attachment)
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
		attachments: []Attachment{
			{filename: "small.txt", content: []byte("small content")},
			{filename: "large.txt", content: make([]byte, 30*1024*1024)}, // 30 MB
		},
		maxAttachmentSize: 25 * 1024 * 1024, // 25 MB
	}

	t.Run("GetAttachments with size limit", func(t *testing.T) {
		expected := []Attachment{
			{filename: "small.txt", content: []byte("small content")},
		}
		result := email.GetAttachments()
		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments with no size limit", func(t *testing.T) {
		email.SetMaxAttachmentSize(-1)
		expected := email.attachments
		result := email.GetAttachments()
		assert.Equal(t, expected, result)
	})
}

func TestSetFilename(t *testing.T) {
	attachment := &Attachment{}
	t.Run("SetFilename", func(t *testing.T) {
		expected := "newfile.txt"
		attachment.SetFilename(expected)
		assert.Equal(t, expected, attachment.filename)
	})
}

func TestSetContent(t *testing.T) {
	attachment := &Attachment{}
	t.Run("SetContent", func(t *testing.T) {
		expected := []byte("new content")
		attachment.SetContent(expected)
		assert.Equal(t, expected, attachment.content)
	})
}

func TestSanitizeInput(t *testing.T) {
	t.Run("sanitize input with HTML", func(t *testing.T) {
		input := "<div>Test</div>"
		expected := "&lt;div&gt;Test&lt;/div&gt;"
		result := SanitizeInput(input)
		assert.Equal(t, expected, result)
	})

	t.Run("sanitize input with spaces", func(t *testing.T) {
		input := "  Test  "
		expected := "Test"
		result := SanitizeInput(input)
		assert.Equal(t, expected, result)
	})
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

func TestMarshalJSONEdgeCases(t *testing.T) {
	t.Run("nil EmailMessage", func(t *testing.T) {
		var email *EmailMessage
		result, err := json.Marshal(email)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})

	t.Run("nil Attachment", func(t *testing.T) {
		var attachment *Attachment
		result, err := json.Marshal(attachment)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})
}

func TestUnmarshalJSONEdgeCases(t *testing.T) {
	t.Run("empty JSON EmailMessage", func(t *testing.T) {
		jsonData := `{}`
		var email EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "", email.from)
		assert.Equal(t, "", email.subject)
		assert.Equal(t, "", email.text)
		assert.Equal(t, "", email.html)
	})

	t.Run("invalid JSON EmailMessage", func(t *testing.T) {
		// Example of JSON data that will cause the unmarshaling to fail
		// 'to' should be an array, not a string
		invalidJSONData := `{
        "from": "sender@example.com",
        "to": "invalid_recipient@example.com",
        "cc": ["cc@example.com"],
        "bcc": ["bcc@example.com"],
        "replyTo": "replyto@example.com",
        "subject": "Subject",
        "text": "This is the email content.",
        "html": "<p>This is the email content.</p>",
        "attachments": [{"filename": "attachment1.txt", "content": "Y29udGVudDE="}]
    }`

		var email EmailMessage
		err := json.Unmarshal([]byte(invalidJSONData), &email)
		assert.Error(t, err)
	})

	t.Run("invalid JSON Attachment", func(t *testing.T) {
		// Example of JSON data that will cause the unmarshaling to fail
		// 'filename' should be an string, not an integer
		jsonData := `{"filename": 123456789, "content": "invalid_base64"}`
		var attachment Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.Error(t, err)
	})

	t.Run("empty JSON Attachment", func(t *testing.T) {
		jsonData := `{}`
		var attachment Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.Nil(t, err)
		assert.Equal(t, "", attachment.filename)
		assert.Equal(t, []byte{}, attachment.content)
	})

	t.Run("invalid base64 content Attachment", func(t *testing.T) {
		jsonData := `{"filename": "file.txt", "content": "invalid_base64"}`
		var attachment Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.NotNil(t, err)
	})

}

func TestIsHTMLEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		input := ""
		result := isHTML(input)
		assert.False(t, result)
	})

	t.Run("string without HTML tags", func(t *testing.T) {
		input := "Just a plain text"
		result := isHTML(input)
		assert.False(t, result)
	})

	t.Run("string with incomplete HTML tag", func(t *testing.T) {
		input := "<div>Test"
		result := isHTML(input)
		assert.True(t, result)
	})
}

func TestNewEmailMessageEdgeCases(t *testing.T) {
	t.Run("create email with empty body", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		subject := "Subject"
		body := ""
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, body, email.text)
		assert.Equal(t, "", email.html)
	})

	t.Run("create email with only spaces in body", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		subject := "Subject"
		body := "     "
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, "     ", email.text)
		assert.Equal(t, "", email.html)
	})
}

func TestNewFullEmailMessageEdgeCases(t *testing.T) {
	t.Run("create full email message with no attachments", func(t *testing.T) {
		from := "sender@example.com"
		to := []string{"recipient@example.com"}
		cc := []string{"cc@example.com"}
		bcc := []string{"bcc@example.com"}
		replyTo := "replyto@example.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachments := []Attachment{}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, cc, email.cc)
		assert.Equal(t, bcc, email.bcc)
		assert.Equal(t, replyTo, email.replyTo)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, text, email.text)
		assert.Equal(t, html, email.html)
		assert.Equal(t, attachments, email.attachments)
	})

	t.Run("create full email message with empty fields", func(t *testing.T) {
		from := ""
		to := []string{}
		cc := []string{}
		bcc := []string{}
		replyTo := ""
		subject := ""
		text := ""
		html := ""
		attachments := []Attachment{}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.from)
		assert.Equal(t, to, email.to)
		assert.Equal(t, cc, email.cc)
		assert.Equal(t, bcc, email.bcc)
		assert.Equal(t, replyTo, email.replyTo)
		assert.Equal(t, subject, email.subject)
		assert.Equal(t, text, email.text)
		assert.Equal(t, html, email.html)
		assert.Equal(t, attachments, email.attachments)
	})
}

func TestAddToRecipientEdgeCases(t *testing.T) {
	t.Run("Add multiple recipients", func(t *testing.T) {
		email := &EmailMessage{}
		recipients := []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"}
		for _, recipient := range recipients {
			email.AddToRecipient(recipient)
		}
		assert.Equal(t, recipients, email.to)
	})

	t.Run("Add recipient to nil EmailMessage", func(t *testing.T) {
		var email *EmailMessage
		assert.Panics(t, func() { email.AddToRecipient("recipient@example.com") })
	})
}

func TestSetCCEdgeCases(t *testing.T) {
	t.Run("SetCC with empty slice", func(t *testing.T) {
		email := &EmailMessage{}
		expected := []string{}
		email.SetCC(expected)
		assert.Equal(t, expected, email.cc)
	})
}

func TestSetBCCEdgeCases(t *testing.T) {
	t.Run("SetBCC with empty slice", func(t *testing.T) {
		email := &EmailMessage{}
		expected := []string{}
		email.SetBCC(expected)
		assert.Equal(t, expected, email.bcc)
	})
}

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
		result := isHTML(test.input)
		if result != test.expected {
			t.Errorf("isHTML(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestAttachmentGetters(t *testing.T) {
	t.Run("GetFilename", func(t *testing.T) {
		attachment := Attachment{filename: "test.txt"}
		assert.Equal(t, "test.txt", attachment.GetFilename())
		assert.Equal(t, "nil_attachment", (*Attachment)(nil).GetFilename())
	})

	t.Run("GetBase64Content", func(t *testing.T) {
		attachment := Attachment{filename: "test.txt", content: []byte("hello")}
		expected := []byte(base64.StdEncoding.EncodeToString([]byte("hello")))
		assert.Equal(t, expected, attachment.GetBase64Content())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetBase64Content())
	})

	t.Run("GetRawContent", func(t *testing.T) {
		attachment := Attachment{filename: "test.txt", content: []byte("hello")}
		expected := []byte("hello")
		assert.Equal(t, expected, attachment.GetRawContent())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetRawContent())
	})
}

func TestNewAttachment(t *testing.T) {
	filename := "testfile.txt"
	content := []byte("This is a test file content.")
	attachment := NewAttachment(filename, content)

	if attachment.filename != filename {
		t.Errorf("NewAttachment() = %v; want %v", attachment.filename, filename)
	}

	if string(attachment.content) != string(content) {
		t.Errorf("NewAttachment() content = %v; want %v", string(attachment.content), string(content))
	}
}

func TestAttachmentEdgeCases(t *testing.T) {
	t.Run("GetBase64Content with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		assert.Equal(t, []byte{}, attachment.GetBase64Content())
	})

	t.Run("GetBase64StringContent with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		assert.Equal(t, "", attachment.GetBase64StringContent())
	})

	t.Run("SetContent with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		attachment.SetContent(nil)
		assert.Nil(t, attachment.content)
	})

	t.Run("SetFilename with empty string", func(t *testing.T) {
		attachment := &Attachment{}
		attachment.SetFilename("")
		assert.Equal(t, "", attachment.filename)
	})
}
