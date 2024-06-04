package ses

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/darkrockmountain/gomail"
)

// sesSender is an interface that defines the method for sending emails using the Amazon SES API.
// Implement this interface to abstract the SES client behavior, making it easier to mock in tests.
type sesSender interface {
	// SendEmail sends an email using the SES API.
	// Parameters:
	//   - input: A pointer to the ses.SendEmailInput struct containing the email details.
	// Returns:
	//   - *ses.SendEmailOutput: The output from the SES service, containing the message ID.
	//   - error: An error if the email sending fails, otherwise nil.
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

// sESEmailSender defines a struct for sending emails using the Amazon SES API.
// This struct contains the SendEmail method and the SES client for authentication.
type sESEmailSender struct {
	sesClient sesSender
	sender    string
}

// NewSESEmailSenderWithCredentials creates a new instance of SESEmailSender.
// It initializes the SES email sender with the provided AWS region, sender email address, and AWS credentials.
//
// Parameters:
//   - region: The AWS region to be used for sending emails.
//   - sender: The email address to be used as the sender.
//   - accessKeyID: The AWS access key ID.
//   - secretAccessKey: The AWS secret access key.
//
// Returns:
//   - *SESEmailSender: A pointer to the initialized SESEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSESEmailSenderWithCredentials(region, sender, accessKeyID, secretAccessKey string) (*sESEmailSender, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	sesClient := ses.New(sess)
	return &sESEmailSender{
		sesClient: sesClient,
		sender:    sender,
	}, nil
}

// NewSESEmailSender creates a new instance of SESEmailSender.
// It initializes the SES email sender with the provided AWS region and sender email address.
//
// Parameters:
//   - region: The AWS region to be used for sending emails.
//   - sender: The email address to be used as the sender.
//
// Returns:
//   - *SESEmailSender: A pointer to the initialized SESEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSESEmailSender(region, sender string) (*sESEmailSender, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &sESEmailSender{
		sesClient: ses.New(sess),
		sender:    sender,
	}, nil
}

// SendEmail sends an email using the Amazon SES API.
// It constructs the email message from the given EmailMessage and sends it using the SES API.
//
// Parameters:
//   - message: A pointer to an EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *sESEmailSender) SendEmail(message *gomail.EmailMessage) error {
	input := &ses.SendEmailInput{
		Source: aws.String(s.sender),
		Destination: &ses.Destination{
			ToAddresses:  aws.StringSlice(message.GetTo()),
			CcAddresses:  aws.StringSlice(message.GetCC()),
			BccAddresses: aws.StringSlice(message.GetBCC()),
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(message.GetSubject()),
			},
			Body: &ses.Body{},
		},
	}

	if message.GetText() != "" {
		input.Message.Body.Text = &ses.Content{
			Data: aws.String(message.GetText()),
		}
	}

	if message.GetHTML() != "" {
		input.Message.Body.Html = &ses.Content{
			Data: aws.String(message.GetHTML()),
		}
	}

	if message.GetReplyTo() != "" {
		replyTo := aws.StringSlice([]string{message.GetReplyTo()})
		input.ReplyToAddresses = replyTo
	}

	// attachments are not supported by AWS SES
	// if len(message.Attachments) > 0 {
	// 	return fmt.Errorf("attachments are not supported in this implementation")
	// }

	_, err := s.sesClient.SendEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
