package providers

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

func TestBuildMimeMessage(t *testing.T) {
	tests := []struct {
		message  *gomail.EmailMessage
		contains []string
	}{
		{
			gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email."),
			[]string{"From: sender@example.com", "To: recipient@example.com", "Subject: Test Email", "This is a test email."},
		},
		{
			gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "<p>This is a test email.</p>"),
			[]string{"From: sender@example.com", "To: recipient@example.com", "Subject: Test Email", "Content-Type: text/html", "<p>This is a test email.</p>"},
		},
		{
			gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetAttachments([]gomail.Attachment{*gomail.NewAttachment("test.txt", []byte("This is a test attachment."))}),
			[]string{"From: sender@example.com", "To: recipient@example.com", "Cc: cc@example.com", "Subject: Test Email", "This is a test email.", "Content-Disposition: attachment; filename=\"test.txt\"", base64.StdEncoding.EncodeToString([]byte("This is a test attachment."))},
		},
		{
			gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetReplyTo("reply-to@example.com"),
			[]string{"From: sender@example.com", "To: recipient@example.com", "Cc: cc@example.com", "Subject: Test Email", "This is a test email.", "Reply-To: reply-to@example.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.message.GetSubject(), func(t *testing.T) {
			result, err := BuildMimeMessage(test.message)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, substring := range test.contains {
				if !bytes.Contains(result, []byte(substring)) {
					t.Fatalf("expected result to contain '%s'", substring)
				}
			}
		})
	}
}

func TestStrPtr(t *testing.T) {

	str := "String to test for pointer"
	ptrStr := StrPtr(str)
	assert.Equal(t, ptrStr, &str)
	assert.EqualValues(t, ptrStr, &str)
}
