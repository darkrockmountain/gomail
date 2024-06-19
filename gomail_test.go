package gomail

import (
	"os"
	"testing"

	"github.com/darkrockmountain/gomail/common"
	"github.com/stretchr/testify/assert"
)

type MockEmailSender struct{}

func (m *MockEmailSender) SendEmail(message *EmailMessage) error {
	return nil
}

func TestNewEmailMessage(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	subject := "Test Subject"
	body := "This is a test email body."

	emailMessage := NewEmailMessage(from, to, subject, body)

	assert.Equal(t, from, emailMessage.GetFrom())
	assert.Equal(t, to, emailMessage.GetTo())
	assert.Equal(t, subject, emailMessage.GetSubject())
	assert.Equal(t, body, emailMessage.GetText())
	assert.Empty(t, emailMessage.GetHTML())
}

func TestNewFullEmailMessage(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	cc := []string{"cc@example.com"}
	bcc := []string{"bcc@example.com"}
	replyTo := "replyto@example.com"
	subject := "Test Subject"
	textBody := "This is a test text body."
	htmlBody := "<p>This is a test HTML body.</p>"
	attachments := []common.Attachment{
		*NewAttachment("test.txt", []byte("test content")),
	}

	emailMessage := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, textBody, htmlBody, attachments)

	assert.Equal(t, from, emailMessage.GetFrom())
	assert.Equal(t, to, emailMessage.GetTo())
	assert.Equal(t, cc, emailMessage.GetCC())
	assert.Equal(t, bcc, emailMessage.GetBCC())
	assert.Equal(t, replyTo, emailMessage.GetReplyTo())
	assert.Equal(t, subject, emailMessage.GetSubject())
	assert.Equal(t, textBody, emailMessage.GetText())
	assert.Equal(t, htmlBody, emailMessage.GetHTML())
	assert.Equal(t, attachments, emailMessage.GetAttachments())
}

func TestNewAttachment(t *testing.T) {
	filename := "test.txt"
	content := []byte("test content")

	attachment := NewAttachment(filename, content)

	assert.Equal(t, filename, attachment.GetFilename())
	assert.Equal(t, content, attachment.GetRawContent())
}

func TestNewAttachmentFromFile(t *testing.T) {
	filePath := "testdata/testfile.txt"
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	attachment, err := NewAttachmentFromFile(filePath)

	assert.NoError(t, err)
	assert.Equal(t, "testfile.txt", attachment.GetFilename())
	assert.Equal(t, content, attachment.GetRawContent())
}

func TestBuildMimeMessage(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	subject := "Test Subject"
	body := "This is a test email body."

	emailMessage := NewEmailMessage(from, to, subject, body)
	mimeMessage, err := BuildMimeMessage(emailMessage)

	assert.NoError(t, err)
	assert.NotEmpty(t, mimeMessage)
}

func TestValidateEmail(t *testing.T) {
	validEmail := "valid@example.com"
	invalidEmail := "invalid-email"

	assert.Equal(t, validEmail, ValidateEmail(validEmail))
	assert.Empty(t, ValidateEmail(invalidEmail))
}

func TestValidateEmailSlice(t *testing.T) {
	emails := []string{"valid@example.com", "invalid-email"}
	validEmails := ValidateEmailSlice(emails)

	assert.Len(t, validEmails, 1)
	assert.Equal(t, "valid@example.com", validEmails[0])
}

func TestGetMimeType(t *testing.T) {
	assert.Contains(t, GetMimeType("test.txt"), "text/plain")
	assert.Contains(t, GetMimeType("tarFile.tar"), "application/x-tar")
}

func TestSendEmail(t *testing.T) {
	sender := &MockEmailSender{}
	message := NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Subject", "This is a test email body.")

	err := sender.SendEmail(message)
	assert.NoError(t, err)
}
