package sparkpost

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

type mockSparkPostClient struct {
	client *sp.Client
	server *httptest.Server
}

func newMockSparkPostClient(handler http.HandlerFunc) *mockSparkPostClient {
	server := httptest.NewTLSServer(handler)

	// Create a custom HTTP transport with TLS verification disabled
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Set the custom transport to the default HTTP client
	http.DefaultTransport = transport

	cfg := &sp.Config{
		BaseUrl: server.URL,
		ApiKey:  "test-api-key",
	}

	client := &sp.Client{}
	client.Init(cfg)

	return &mockSparkPostClient{client: client, server: server}
}

func (m *mockSparkPostClient) Close() {
	m.server.Close()
}

func TestNewSparkPostEmailSender(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender, err := NewSparkPostEmailSender("test-api-key")
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestSparkPostEmailSender_SendEmail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

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

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSparkPostEmailSender_SendEmailWithError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
	}

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email with SparkPost")
}

func TestSparkPostEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{},
		Subject: "",
		Text:    "",
	}

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSparkPostEmailSender_SendEmailWithAttachments(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
		Attachments: []gomail.Attachment{
			{
				Filename: "test.txt",
				Content:  []byte("This is a test attachment."),
			},
		},
	}

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSparkPostEmailSender_SendEmailWithReplyTo(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
		ReplyTo: "replyto@example.com",
	}

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSparkPostEmailSender_SendEmailWithCC(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
		CC:      []string{"cc@example.com"},
	}

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSparkPostEmailSender_SendEmailWithBCC(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]interface{}{
			"results": map[string]interface{}{
				"id": "1234567890",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	mockClient := newMockSparkPostClient(handler)
	defer mockClient.Close()

	emailSender := &sparkPostEmailSender{client: mockClient.client}

	message := gomail.EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Text:    "This is a test email.",
		BCC:     []string{"bcc@example.com"},
	}

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}
