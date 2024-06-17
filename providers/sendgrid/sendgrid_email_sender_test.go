package sendgrid

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkrockmountain/gomail"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

// TestEmailSenderImplementation checks if sendGridEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ gomail.EmailSender = (*sendGridEmailSender)(nil)
}

// Mocking the SendGrid API response
func mockSendGridServer(t *testing.T, statusCode int, responseBody string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})
	return httptest.NewServer(handler)
}

func NewMockSendGridEmailSender(apiKey, url string) *sendGridEmailSender {

	requestHeaders := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"User-Agent":    "sendgrid/" + "3.14.0" + ";go",
		"Accept":        "application/json",
	}

	request := rest.Request{
		Method:  "POST",
		BaseURL: url,
		Headers: requestHeaders,
	}

	return &sendGridEmailSender{
		client: &sendgrid.Client{Request: request},
	}
}

func TestNewSendGridEmailSender(t *testing.T) {
	apiKey := "test-api-key"
	emailSender, err := NewSendGridEmailSender(apiKey)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestSendGridEmailSender_SendEmail(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
		SetCC([]string{"cc@example.com"}).
		SetBCC([]string{"bcc@example.com"}).
		SetReplyTo("replyto@example.com").
		SetHTML("<p>This is a test email.</p>").
		SetBCC([]string{"bcc@example.com"}).
		AddAttachment(*gomail.NewAttachment("test.txt", []byte("This is a test attachment.")))

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendGridEmailSender_SendEmailWithError(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusInternalServerError, `{"message": "error"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email, status code: 500")
}

func TestSendGridEmailSender_SendEmailWithNon200StatusCode(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusBadRequest, `{"message": "Bad Request"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email, status code: 400")
}

func TestSendGridEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{},
		"",
		"",
	)

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendGridEmailSender_SendEmailWithAttachments(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	attachmentContent := "This is a test attachment."
	attachmentContentBase64 := base64.StdEncoding.EncodeToString([]byte(attachmentContent))

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(*gomail.NewAttachment("test.txt", []byte(attachmentContent)))

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)

	// Verify the attachment content
	v3Mail := mail.NewV3Mail()
	attachment := mail.NewAttachment()
	attachment.SetContent(attachmentContentBase64)
	attachment.SetType("text/plain")
	attachment.SetFilename("test.txt")
	attachment.SetDisposition("attachment")
	v3Mail.AddAttachment(attachment)

	assert.Equal(t, v3Mail.Attachments[0].Content, message.GetAttachments()[0].GetBase64StringContent())
}
