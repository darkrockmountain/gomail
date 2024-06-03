// mandrill_email_sender.go
package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/darkrockmountain/gomail"
	"github.com/pkg/errors"
)

const (
	mandrillRequestMethod = "POST"
	mandrillRequestURL    = "https://mandrillapp.com/api/1.0/messages/send.json"
)

// mandrillEmailSender defines a struct for sending emails using the Mandrill API.
// This struct contains the SendEmail method and the Mandrill API key for authentication.
type mandrillEmailSender struct {
	apiKey        string
	requestMethod string
	url           string
}

// NewMandrillEmailSender creates a new instance of mandrillEmailSender.
// It initializes the Mandrill email sender with the provided API key.
//
// Parameters:
//   - apiKey: The API key to be used for authentication.
//
// Returns:
//   - *mandrillEmailSender: A pointer to the initialized mandrillEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMandrillEmailSender(apiKey string) (*mandrillEmailSender, error) {
	return &mandrillEmailSender{apiKey: apiKey, requestMethod: mandrillRequestMethod, url: mandrillRequestURL}, nil
}

// SendEmail sends an email using the Mandrill API.
// It constructs the email message from the given EmailMessage and sends it using the Mandrill API.
//
// Parameters:
//   - message: An EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *mandrillEmailSender) SendEmail(message gomail.EmailMessage) error {
	maMessage := mandrillMessage{
		FromEmail: message.GetFrom(),
		To:        make([]mandrillRecipient, 0),
		Subject:   message.GetSubject(),
		Text:      message.GetText(),
		HTML:      message.GetHTML(),
	}

	// Add Reply-To if specified
	replyTo := message.GetReplyTo()
	if replyTo != "" {
		maMessage.Headers = map[string]string{
			"Reply-To": replyTo,
		}
	}

	// Add recipients
	for _, to := range message.GetTo() {
		maMessage.To = append(maMessage.To, mandrillRecipient{
			Email: to,
			Type:  "to",
		})
	}

	// Add CC recipients
	for _, cc := range message.GetCC() {
		maMessage.To = append(maMessage.To, mandrillRecipient{
			Email: cc,
			Type:  "cc",
		})
	}

	// Add BCC recipients
	for _, bcc := range message.GetBCC() {
		maMessage.To = append(maMessage.To, mandrillRecipient{
			Email: bcc,
			Type:  "bcc",
		})
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		fileName := attachment.GetFilename()
		maMessage.Attachments = append(maMessage.Attachments, mandrillAttachment{
			Type:    gomail.GetMimeType(fileName),
			Name:    fileName,
			Content: attachment.GetBase64StringContent(),
		})
	}

	// Create the payload for the Mandrill API
	payload := mandrillRequest{
		Key:     s.apiKey,
		Message: maMessage,
	}

	// Send the request to the Mandrill API
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal email data")
	}

	req, err := http.NewRequest(s.requestMethod, s.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send email via Mandrill API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email via Mandrill API, status code: %d", resp.StatusCode)
	}

	return nil
}

// mandrillRequest represents the payload for the Mandrill API.
type mandrillRequest struct {
	Key     string          `json:"key"`
	Message mandrillMessage `json:"message"`
}

// mandrillMessage represents an email message to be sent via Mandrill.
type mandrillMessage struct {
	FromEmail   string               `json:"from_email"`
	To          []mandrillRecipient  `json:"to"`
	Subject     string               `json:"subject"`
	Text        string               `json:"text,omitempty"`
	HTML        string               `json:"html,omitempty"`
	Headers     map[string]string    `json:"headers,omitempty"`
	Attachments []mandrillAttachment `json:"attachments,omitempty"`
	Metadata    map[string]string    `json:"metadata,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Subaccount  string               `json:"subaccount,omitempty"`
}

// mandrillRecipient represents a recipient of an email message.
type mandrillRecipient struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

// mandrillAttachment represents an attachment for a Mandrill email message.
type mandrillAttachment struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
