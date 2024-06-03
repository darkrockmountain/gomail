package providers

import (
	"fmt"

	"github.com/darkrockmountain/gomail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// sendGridEmailSender defines a struct for sending emails using the SendGrid API.
// This struct contains the SendEmail method and an API key for authentication.
type sendGridEmailSender struct {
	client *sendgrid.Client
}

// NewSendGridEmailSender creates a new instance of sendGridEmailSender.
// It initializes the SendGrid email sender with the provided API key.
//
// Parameters:
//   - apiKey: The API key to be used for authentication.
//
// Returns:
//   - *sendGridEmailSender: A pointer to the initialized sendGridEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSendGridEmailSender(apiKey string) (*sendGridEmailSender, error) {
	return &sendGridEmailSender{
		client: sendgrid.NewSendClient(apiKey),
	}, nil
}

// SendEmail sends an email using the SendGrid API.
// It constructs the email message from the given EmailMessage and sends it using the SendGrid API.
//
// Parameters:
//   - message: An EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *sendGridEmailSender) SendEmail(message gomail.EmailMessage) error {
	from := mail.NewEmail("", message.GetFrom())
	toRecipients := []*mail.Email{}

	for _, to := range message.GetTo() {
		toRecipients = append(toRecipients, mail.NewEmail("", to))
	}
	v3Mail := mail.NewV3Mail()
	v3Mail.SetFrom(from)
	v3Mail.Subject = message.GetSubject()

	// Create personalization for To recipients
	personalization := mail.NewPersonalization()
	for _, to := range toRecipients {
		personalization.AddTos(to)
	}

	// Add BCC recipients
	for _, bcc := range message.GetBCC() {
		personalization.AddBCCs(mail.NewEmail("", bcc))
	}

	// Add Reply-To if specified
	if message.GetReplyTo() != "" {
		replyTo := mail.NewEmail("", message.GetReplyTo())
		v3Mail.SetReplyTo(replyTo)
	}

	v3Mail.AddPersonalizations(personalization)

	// Add plain text content
	if message.GetText() != "" {
		v3Mail.AddContent(mail.NewContent("text/plain", message.GetText()))
	}

	// Add HTML content
	if message.GetHTML() != "" {
		v3Mail.AddContent(mail.NewContent("text/html", message.GetHTML()))
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		a := mail.NewAttachment()
		a.SetContent(attachment.GetBase64StringContent())
		a.SetType(gomail.GetMimeType(attachment.GetFilename()))
		a.SetFilename(attachment.GetFilename())
		a.SetDisposition("attachment")
		v3Mail.AddAttachment(a)
	}

	response, err := s.client.Send(v3Mail)
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
