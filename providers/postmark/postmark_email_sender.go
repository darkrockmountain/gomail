package postmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/darkrockmountain/gomail"
	"github.com/pkg/errors"
)

const (
	postMarkRequestMethod = "POST"
	postMarkRequestURL    = "https://api.postmarkapp.com/email"
	clientTimeOut         = time.Millisecond * 100
)

// postmarkEmailSender defines a struct for sending emails using the Postmark API.
// This struct contains the SendEmail method and the necessary authentication token.
type postmarkEmailSender struct {
	serverToken   string
	requestMethod string
	url           string
}

// NewPostmarkEmailSender creates a new instance of postmarkEmailSender.
// It initializes the Postmark email sender with the provided server token.
//
// Parameters:
//   - serverToken: The server token to be used for authentication.
//
// Returns:
//   - *postmarkEmailSender: A pointer to the initialized postmarkEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewPostmarkEmailSender(serverToken string) (*postmarkEmailSender, error) {
	return &postmarkEmailSender{serverToken: serverToken, requestMethod: postMarkRequestMethod, url: postMarkRequestURL}, nil
}

// SendEmail sends an email using the Postmark API.
// It constructs the email message from the given EmailMessage and sends it using the Postmark API.
//
// Parameters:
//   - message: A pointer to an EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *postmarkEmailSender) SendEmail(message *gomail.EmailMessage) error {
	emailStruct := struct {
		From        string               `json:"From"`
		To          string               `json:"To"`
		CC          string               `json:"Cc,omitempty"`
		Subject     string               `json:"Subject"`
		TextBody    string               `json:"TextBody,omitempty"`
		HtmlBody    string               `json:"HtmlBody,omitempty"`
		ReplyTo     string               `json:"ReplyTo,omitempty"`
		Bcc         string               `json:"Bcc,omitempty"`
		Attachments []postmarkAttachment `json:"Attachments,omitempty"`
	}{
		From:     message.GetFrom(),
		To:       strings.Join(message.GetTo(), ","),
		CC:       strings.Join(message.GetCC(), ","),
		Subject:  message.GetSubject(),
		TextBody: message.GetText(),
		HtmlBody: message.GetHTML(),
		ReplyTo:  message.GetReplyTo(),
		Bcc:      strings.Join(message.GetBCC(), ","),
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		emailStruct.Attachments = append(emailStruct.Attachments, postmarkAttachment{
			Name:        attachment.GetFilename(),
			Content:     attachment.GetBase64StringContent(),
			ContentType: gomail.GetMimeType(attachment.GetFilename()),
		})
	}

	jsonData, err := json.Marshal(emailStruct)
	if err != nil {
		return errors.Wrap(err, "failed to marshal email data")
	}

	req, err := http.NewRequest(s.requestMethod, s.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.serverToken)

	client := &http.Client{Timeout: clientTimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send email via Postmark API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email via Postmark API, status code: %d", resp.StatusCode)
	}

	return nil
}

// postmarkAttachment represents an attachment for a Postmark email.
type postmarkAttachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}
