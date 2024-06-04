package postmark

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

func TestNewPostmarkEmailSender(t *testing.T) {
	serverToken := "test-server-token"
	emailSender, err := NewPostmarkEmailSender(serverToken)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
	assert.Equal(t, serverToken, emailSender.serverToken)
	assert.Equal(t, postMarkRequestMethod, emailSender.requestMethod)
	assert.Equal(t, postMarkRequestURL, emailSender.url)
}

func TestPostmarkEmailSender_SendEmail(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)
	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
		SetCC([]string{"cc@example.com"}).
		SetBCC([]string{"bcc@example.com"}).
		SetReplyTo("replyto@example.com").
		SetHTML("<p>This is a test email.</p>").
		SetBCC([]string{"bcc@example.com"}).
		AddAttachment(gomail.Attachment{
			Filename: "test.txt",
			Content:  []byte("This is a test attachment."),
		})

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

func TestPostmarkEmailSender_SendEmailWithMarshalError(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)
	message := *gomail.NewEmailMessage(
		string(make([]byte, 1<<20)), // Intentionally large string to cause marshal error
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestPostmarkEmailSender_SendEmailWithRequestCreationError(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	emailSender.url = "no a url"
	emailSender.requestMethod = "no a request method"

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestPostmarkEmailSender_SendEmailWithSendError(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	// Mock server to simulate a server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * clientTimeOut) // Wait for 2 times the clientTimeout

		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	emailSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email via Postmark API")
}

func TestPostmarkEmailSender_SendEmailWithNon200StatusCode(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

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

func TestPostmarkEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	emailSender, err := NewPostmarkEmailSender("test-server-token")
	assert.NoError(t, err)

	message := *gomail.NewEmailMessage(
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
