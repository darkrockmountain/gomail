package mandrill

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

// TestEmailSenderImplementation checks if mandrillEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ gomail.EmailSender = (*mandrillEmailSender)(nil)
}

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

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
		SetCC([]string{"cc@example.com"}).
		SetBCC([]string{"bcc@example.com"}).
		SetReplyTo("replyto@example.com").
		SetHTML("<p>This is a test email.</p>").
		SetBCC([]string{"bcc@example.com"}).
		AddAttachment(*gomail.NewAttachment("test.txt", []byte("This is a test attachment.")))

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

	message := gomail.NewEmailMessage(
		string(make([]byte, 1<<20)), // Intentionally large string to cause error
		[]string{"recipient@example.com"},
		"Test Email",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestMandrillEmailSender_SendEmailNewRequestError(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	emailSender.url = "no a url"
	emailSender.requestMethod = "no a request method"

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create HTTP request")
}

func TestMandrillEmailSender_SendEmailWithSendError(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

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

func TestMandrillEmailSender_SendEmailWithSendTimeOutError(t *testing.T) {
	emailSender, err := NewMandrillEmailSender("test-api-key")
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	// Mock server to simulate a server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * clientTimeOut) // Wait for 2 times the clientTimeout
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

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

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

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{},
		"",
		"",
	)

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
