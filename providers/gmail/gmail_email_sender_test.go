package gmail

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// MockTokenManager is a mock implementation of the GmailTokenManager interface.
type MockTokenManager struct {
	token []byte
	err   error
}

func (m *MockTokenManager) GetToken() ([]byte, error) {
	return m.token, m.err
}

// GmailMockRoundTripper implements the http.RoundTripper interface
type GmailMockRoundTripper struct {
	Response *http.Response
	Err      error
}

func readRequestBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	return body, nil
}

func (m *GmailMockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	reqBody, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(reqBody)),
		Header:     make(http.Header),
	}, nil
}

// MockGmailService is a mock implementation of the Gmail UsersMessagesService
type MockGmailService struct{}

func (m *MockGmailService) Send(userId string, message *gmail.Message) *gmail.UsersMessagesSendCall {
	srv, _ := gmail.NewService(context.Background())
	return gmail.NewUsersMessagesService(srv).Send(userId, message)
}

// buildMockGmailMessageSenderWrapper builds a mock gmailMessageSenderWrapper for testing.
func buildMockGmailMessageSenderWrapper(err error) *gmailMessageSenderWrapper {
	mockGmailService, _ := gmail.NewService(context.Background(), option.WithHTTPClient(&http.Client{
		Transport: &GmailMockRoundTripper{
			Err: err,
		},
	}))

	return &gmailMessageSenderWrapper{
		messageSender: mockGmailService.Users.Messages,
		user:          "me",
	}
}

type mockGmailTokenManager struct{}

func (m *mockGmailTokenManager) GetToken() ([]byte, error) {
	token := &oauth2.Token{
		AccessToken: "mockAccessToken",
		TokenType:   "Bearer",
	}
	return json.Marshal(token)
}

// Mock implementations for new edge cases
type mockInvalidTokenManager struct{}

func (m *mockInvalidTokenManager) GetToken() ([]byte, error) {
	return nil, errors.New("invalid token")
}

func TestSendEmailWithMockService(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(nil)

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailWithSendMessageError(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(errors.New("send message error"))
	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithMockServiceError(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(errors.New("mock service error"))

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithBCC(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(nil)

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.").
		SetBCC([]string{"bcc@example.com"})

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)

}

func TestNewGmailEmailSenderOauth2(t *testing.T) {
	configJson := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockGmailTokenManager{}
	user := "me"

	emailSender, err := NewGmailEmailSenderOauth2(context.Background(), configJson, tokenManager, user)
	assert.NoError(t, err)
	assert.Equal(t, user, emailSender.user)
	assert.NotNil(t, emailSender.messageSender)
}

func TestGmailEmailSenderOauth2GetGmailMessageSender(t *testing.T) {
	configJson := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockGmailTokenManager{}
	user := "me"

	emailSender, err := NewGmailEmailSenderOauth2(context.Background(), configJson, tokenManager, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderOauth2InvalidToken(t *testing.T) {
	configJson := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockInvalidTokenManager{}
	user := "me"

	_, err := NewGmailEmailSenderOauth2(context.Background(), configJson, tokenManager, user)
	assert.Error(t, err)
}

func TestNewGmailEmailSenderServiceAccount(t *testing.T) {
	jsonCredentials := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewGmailEmailSenderServiceAccount(context.Background(), jsonCredentials, user)
	assert.NoError(t, err)
	assert.Equal(t, user, emailSender.user)
	assert.NotNil(t, emailSender.messageSender)
}

func TestGmailEmailSenderServiceAccountGetGmailMessageSender(t *testing.T) {
	jsonCredentials := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewGmailEmailSenderServiceAccount(context.Background(), jsonCredentials, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderAPIKey(t *testing.T) {
	apiKey := "mock_api_key"
	user := "me"

	emailSender, err := NewGmailEmailSenderAPIKey(context.Background(), apiKey, user)
	assert.NoError(t, err)
	assert.Equal(t, user, emailSender.user)
	assert.NotNil(t, emailSender.messageSender)
}

func TestGmailEmailSenderAPIKeyGetGmailMessageSender(t *testing.T) {
	apiKey := "mock_api_key"
	user := "me"

	emailSender, err := NewGmailEmailSenderAPIKey(context.Background(), apiKey, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderJWT(t *testing.T) {
	configJson := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewGmailEmailSenderJWT(context.Background(), configJson, user)
	assert.NoError(t, err)
	assert.Equal(t, user, emailSender.user)
	assert.NotNil(t, emailSender.messageSender)
}

func TestGmailEmailSenderJWTGetGmailMessageSender(t *testing.T) {
	configJson := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewGmailEmailSenderJWT(context.Background(), configJson, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderJWTInvalidJson(t *testing.T) {
	configJson := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewGmailEmailSenderJWT(context.Background(), configJson, user)
	assert.Error(t, err)
}

func TestNewGmailEmailSenderJWTAccessInvalidJson(t *testing.T) {
	jsonCredentials := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewGmailEmailSenderJWTAccess(context.Background(), jsonCredentials, user)
	assert.Error(t, err)
}

func TestSendEmailWithNilGmailService(t *testing.T) {
	emailSender := &gmailMessageSenderWrapper{}

	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error no UsersMessagesService initiated")
}

func TestNewGmailEmailSenderServiceAccountInvalidJson(t *testing.T) {
	jsonCredentials := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewGmailEmailSenderServiceAccount(context.Background(), jsonCredentials, user)
	assert.Error(t, err)
}

// func TestNewGmailEmailSenderJWTAccess(t *testing.T) {
// 	jsonCredentials := []byte(`{
// 		"type": "service_account",
// 		"project_id": "mock_project_id",
// 		"private_key_id": "mock_private_key_id",
// 		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
// 		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
// 		"client_id": "mock_client_id",
// 		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
// 		"token_uri": "https://oauth2.googleapis.com/token",
// 		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
// 		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
// 	}`)
// 	user := "me"

// 	emailSender, err := NewGmailEmailSenderJWTAccess(context.Background(), jsonCredentials, user)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, emailSender)
// 	assert.Equal(t, user, emailSender.user)
// 	assert.NotNil(t, emailSender.messageSender)
// }
