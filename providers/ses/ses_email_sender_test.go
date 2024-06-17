package ses

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

// TestEmailSenderImplementation checks if sESEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ gomail.EmailSender = (*sESEmailSender)(nil)
}

// Mock SES client
type mockSESClient struct {
	sesiface.SESAPI
	SendEmailOutput *ses.SendEmailOutput
	SendEmailError  error
}

func (m *mockSESClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return m.SendEmailOutput, m.SendEmailError
}

// Helper function to create a mock SES client
func createMockSESClient(output *ses.SendEmailOutput, err error) *mockSESClient {
	return &mockSESClient{
		SendEmailOutput: output,
		SendEmailError:  err,
	}
}

func TestNewSESEmailSender(t *testing.T) {
	emailSender, err := NewSESEmailSender("us-west-2", "sender@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
	assert.Equal(t, "sender@example.com", emailSender.sender)
}

func TestNewSESEmailSenderWithCredentials(t *testing.T) {
	emailSender, err := NewSESEmailSenderWithCredentials("us-west-2", "sender@example.com", "accessKeyID", "secretAccessKey")
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
	assert.Equal(t, "sender@example.com", emailSender.sender)
}

func TestNewSESEmailSenderError(t *testing.T) {
	os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "made up endpoint")
	defer os.Unsetenv("AWS_STS_REGIONAL_ENDPOINTS")
	emailSender, err := NewSESEmailSender("us-west-2", "sender@example.com")
	assert.Error(t, err)
	assert.Nil(t, emailSender)
}

func TestNewSESEmailSenderWithCredentialsError(t *testing.T) {
	os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "made up endpoint")
	defer os.Unsetenv("AWS_STS_REGIONAL_ENDPOINTS")
	emailSender, err := NewSESEmailSenderWithCredentials("us-west-2", "sender@example.com", "accessKeyID", "secretAccessKey")
	assert.Error(t, err)
	assert.Nil(t, emailSender)
}

func TestSESEmailSender_SendEmail(t *testing.T) {
	mockClient := createMockSESClient(&ses.SendEmailOutput{}, nil)

	emailSender := &sESEmailSender{
		sesClient: mockClient,
		sender:    "sender@example.com",
	}

	message := gomail.NewFullEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Email",
		[]string{"cc@example.com"},
		[]string{"bcc@example.com"},
		"replyto@example.com",
		"This is a test email.",
		"<p>This is a test email.</p>", []gomail.Attachment{},
	)

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSESEmailSender_SendEmailWithError(t *testing.T) {
	mockClient := createMockSESClient(nil, awserr.New(ses.ErrCodeMessageRejected, "mock error", errors.New("mock error")))

	emailSender := &sESEmailSender{
		sesClient: mockClient,
		sender:    "sender@example.com",
	}

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestSESEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	mockClient := createMockSESClient(&ses.SendEmailOutput{}, nil)

	emailSender := &sESEmailSender{
		sesClient: mockClient,
		sender:    "sender@example.com",
	}

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{},
		"",
		"",
	)

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSESEmailSender_SendEmailWithAttachments(t *testing.T) {
	mockClient := createMockSESClient(&ses.SendEmailOutput{}, nil)

	emailSender := &sESEmailSender{
		sesClient: mockClient,
		sender:    "sender@example.com",
	}

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(*gomail.NewAttachment("test.txt", []byte("This is a test attachment.")))

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSESEmailSender_SendEmailWithReplyTo(t *testing.T) {
	mockClient := createMockSESClient(&ses.SendEmailOutput{}, nil)

	emailSender := &sESEmailSender{
		sesClient: mockClient,
		sender:    "sender@example.com",
	}

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).SetReplyTo("replyto@example.com")

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}
