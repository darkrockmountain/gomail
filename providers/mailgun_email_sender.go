package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/darkrockmountain/gomail"
	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSubSet defines a subset of the mailgun.Mailgun interface API.
type MailgunSubSet interface {
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
	Send(ctx context.Context, m *mailgun.Message) (string, string, error)
}

// mailgunEmailSender defines a struct for sending emails using the Mailgun API.
// This struct contains the SendEmail method and the Mailgun client for authentication.
type mailgunEmailSender struct {
	mailgunClient MailgunSubSet
}

// NewMailgunEmailSender creates a new instance of mailgunEmailSender.
// It initializes the Mailgun email sender with the provided domain and API key.
//
// Parameters:
//   - domain: The Mailgun domain to be used for sending emails.
//   - apiKey: The API key to be used for authentication.
//
// Returns:
//   - *mailgunEmailSender: A pointer to the initialized mailgunEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMailgunEmailSender(domain, apiKey string) (*mailgunEmailSender, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	return &mailgunEmailSender{
		mailgunClient: mg,
	}, nil
}

// SendEmail sends an email using the Mailgun API.
// It constructs the email message from the given EmailMessage and sends it using the Mailgun API.
//
// Parameters:
//   - message: An EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *mailgunEmailSender) SendEmail(message gomail.EmailMessage) error {
	mg := s.mailgunClient

	mGMessage := mg.NewMessage(
		message.GetFrom(),
		message.GetSubject(),
		message.GetText(),
		message.GetTo()...,
	)

	html := message.GetHTML()
	if html != "" {
		mGMessage.SetHtml(html)
	}

	toCcs := message.GetCC()
	if len(toCcs) > 0 {
		for _, bcc := range toCcs {
			mGMessage.AddCC(bcc)
		}
	}

	toBccs := message.GetBCC()
	if len(toBccs) > 0 {
		for _, bcc := range toBccs {
			mGMessage.AddBCC(bcc)
		}
	}

	replyTo := message.GetReplyTo()
	if replyTo != "" {
		mGMessage.SetReplyTo(replyTo)
	}

	for _, attachment := range message.GetAttachments() {
		mGMessage.AddBufferAttachment(attachment.GetFilename(), attachment.GetRawContent())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, mGMessage)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
