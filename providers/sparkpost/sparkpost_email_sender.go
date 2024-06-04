package sparkpost

import (
	"fmt"
	"strings"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/darkrockmountain/gomail"
	"github.com/pkg/errors"
)

// sparkPostEmailSender defines a struct for sending emails using the SparkPost API.
// This struct contains the SendEmail method and the SparkPost client for authentication.
type sparkPostEmailSender struct {
	client *sp.Client
}

// NewSparkPostEmailSender creates a new instance of sparkPostEmailSender.
// It initializes the SparkPost email sender with the provided API key.
//
// Parameters:
//   - apiKey: The API key to be used for authentication.
//
// Returns:
//   - *sparkPostEmailSender: A pointer to the initialized sparkPostEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSparkPostEmailSender(apiKey string) (*sparkPostEmailSender, error) {
	cfg := &sp.Config{
		ApiKey: apiKey,
	}

	client := &sp.Client{}
	err := client.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SparkPost client: %v", err)
	}

	return &sparkPostEmailSender{client: client}, nil
}

// SendEmail sends an email using the SparkPost API.
// It constructs the email message from the given EmailMessage and sends it using the SparkPost API.
//
// Parameters:
//   - message: A pointer to an EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *sparkPostEmailSender) SendEmail(message *gomail.EmailMessage) error {
	content := sp.Content{
		From:    message.GetFrom(),
		Subject: message.GetSubject(),
		Text:    message.GetText(),
		HTML:    message.GetHTML(),
	}

	// Add Reply-To if specified
	if message.GetReplyTo() != "" {
		content.ReplyTo = message.GetReplyTo()
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		content.Attachments = append(content.Attachments, sp.Attachment{
			MIMEType: gomail.GetMimeType(attachment.GetFilename()),
			B64Data:  attachment.GetBase64StringContent(),
			Filename: attachment.GetFilename(),
		})
	}

	tx := &sp.Transmission{Content: content}

	headerTo := strings.Join(message.GetTo(), ",")
	// Add recipients
	tx.Recipients = []sp.Recipient{}

	for _, to := range message.GetTo() {
		tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
			Address: sp.Address{Email: to},
		})
	}

	// Add CC recipients
	ccs := message.GetCC()
	if len(ccs) > 0 {
		for _, cc := range ccs {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: cc, HeaderTo: headerTo},
			})
		}
		// add cc header to content
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		content.Headers["cc"] = strings.Join(ccs, ",")
	}

	// Add BCC recipients
	bccs := message.GetBCC()
	if len(bccs) > 0 {
		for _, b := range bccs {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: b, HeaderTo: headerTo},
			})
		}
	}

	_, _, err := s.client.Send(tx)
	if err != nil {
		return errors.Wrap(err, "failed to send email with SparkPost")
	}

	return nil

}
