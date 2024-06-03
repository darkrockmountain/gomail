package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

func TestNewMandrillEmailSender(t *testing.T) {
	apiKey := "test-api-key"
	emailSender, err := NewMandrillEmailSender(apiKey)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
	assert.Equal(t, apiKey, emailSender.apiKey)
}

func TestMandrillEmailSender_SendEmail(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
		HTML:    "<p>This is a test email.</p>",
		CC:      []string{"cc@example.com"},
		BCC:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
		Attachments: []gomail.Attachment{
			{
				Filename: "test.txt",
				Content:  []byte("This is a test attachment."),
			},
		},
	}

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "sent"}`))
	}))
	defer ts.Close()

	emailSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestMandrillEmailSender_SendEmailWithError(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.EmailMessage{
		From:    string(make([]byte, 1<<20)), // Intentionally large string to cause error
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
	}

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestMandrillEmailSender_SendEmailWithSendError(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
	}

	// Mock server to simulate a server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	emailSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email via Mandrill API")
}

func TestMandrillEmailSender_SendEmailWithNon200StatusCode(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
	}

	// Mock server to simulate a non-200 status code response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer ts.Close()

	emailSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status code: 400")
}

func TestMandrillEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{},
		Subject: "",
		Text:    "",
	}

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "sent"}`))
	}))
	defer ts.Close()

	emailSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}
