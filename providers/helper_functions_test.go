package providers

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/darkrockmountain/gomail"
)

func TestBuildMimeMessage(t *testing.T) {
	tests := []struct {
		message  gomail.EmailMessage
		contains []string
	}{
		{
			gomail.EmailMessage{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				Text:    "This is a test email.",
			},
			[]string{"From: sender@example.com", "To: recipient@example.com", "Subject: Test Email", "This is a test email."},
		},
		{
			gomail.EmailMessage{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				HTML:    "<p>This is a test email.</p>",
			},
			[]string{"From: sender@example.com", "To: recipient@example.com", "Subject: Test Email", "Content-Type: text/html", "<p>This is a test email.</p>"},
		},
		{
			gomail.EmailMessage{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				CC:      []string{"cc@example.com"},
				BCC:     []string{"bcc@example.com"},
				Subject: "Test Email",
				Text:    "This is a test email.",
				Attachments: []gomail.Attachment{
					{Filename: "test.txt", Content: []byte("This is a test attachment.")},
				},
			},
			[]string{"From: sender@example.com", "To: recipient@example.com", "Cc: cc@example.com", "Subject: Test Email", "This is a test email.", "Content-Disposition: attachment; filename=\"test.txt\"", base64.StdEncoding.EncodeToString([]byte("This is a test attachment."))},
		},
	}

	for _, test := range tests {
		t.Run(test.message.Subject, func(t *testing.T) {
			result, err := buildMimeMessage(test.message)
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
