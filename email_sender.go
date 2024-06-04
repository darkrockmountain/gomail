// Package gomail provides a unified interface for sending emails using various providers.
//
// # Overview
//
// The gomail project allows you to send emails using different email providers
// such as Gmail, SendGrid, AWS SES, and others. It abstracts the provider-specific
// details and provides a simple API for sending emails.
//
// This project is organized into several packages:
//
// - providers: Contains implementations for various email providers.
// - credentials: Contains implementations for managing email credentials.
// - examples: Contains example applications demonstrating how to use the library.
// - docs: Contains documentation for configuring different email providers.
//
// # Usage
//
// To use the library, you need to import the desired provider package and create an
// instance of the email sender for your desired provider, then call the SendEmail function.
//
// Example:
//
//	package main
//
//	import (
//	    "github.com/darkrockmountain/gomail/providers/sendgrid"
//	)
//
//	func main() {
//	    sender := sendgrid.NewSendGridEmailSender("your-api-key")
//	    err := sender.SendEmail(gomail.NewEmailMessage([]string{"recipient@example.com"},"Subject","Email body"))
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// This library supports various email providers and can be extended to include more.
//
// # Supported Providers
//
// - Gmail
// - SendGrid
// - AWS SES
// - Mailgun
// - Mandrill
// - Postmark
// - Microsoft365
// - SparkPost
// - SMTP
//
// For more details, see the documentation for each provider in the providers package.
package gomail

import (
	"encoding/base64"
	"html"
	"mime"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// EmailSender interface defines the method to send an email.
// Implement this interface to create different email sending strategies.
type EmailSender interface {
	// SendEmail sends an email with the given message.
	// Parameters:
	// - message: An EmailMessage struct containing the details of the email to be sent.
	// Returns:
	// - error: An error if sending the email fails, otherwise nil.
	SendEmail(message EmailMessage) error
}

// EmailMessage contains the fields for sending an email.
// Use this struct to specify the sender, recipient, subject, and content of the email,
// as well as any attachments.
type EmailMessage struct {
	From        string       `json:"from"`        // Sender email address.
	To          []string     `json:"to"`          // Recipient email addresses.
	CC          []string     `json:"cc"`          // CC recipients email addresses.
	BCC         []string     `json:"bcc"`         // BCC recipients email addresses.
	ReplyTo     string       `json:"replyTo"`     // Reply-To email address.
	Subject     string       `json:"subject"`     // Email subject.
	Text        string       `json:"text"`        // Plain text content of the email.
	HTML        string       `json:"html"`        // HTML content of the email (optional).
	Attachments []Attachment `json:"attachments"` // Attachments to be included in the email (optional).
}

// NewEmailMessage creates a new EmailMessage with the required fields.
// If the body contains HTML tags, it sets the HTML field; otherwise, it sets the Text field.
// Parameters:
// - from: The sender email address.
// - to: A slice of recipient email addresses.
// - subject: The email subject.
// - body: The content of the email, which can be plain text or HTML.
//
// Returns:
//   - *EmailMessage: A pointer to the newly created EmailMessage struct.
func NewEmailMessage(from string, to []string, subject string, body string) *EmailMessage {
	email := &EmailMessage{
		From:    from,
		To:      to,
		Subject: subject,
	}

	if isHTML(body) {
		email.HTML = body
	} else {
		email.Text = body
	}

	return email
}

// NewFullEmailMessage creates a new EmailMessage with all fields.
// Parameters:
// - from: The sender email address.
// - to: A slice of recipient email addresses.
// - subject: The email subject.
// - cc: A slice of CC recipient email addresses (optional).
// - bcc: A slice of BCC recipient email addresses (optional).
// - replyTo: The reply-to email address (optional).
// - textBody: The plain text content of the email.
// - htmlBody: The HTML content of the email (optional).
// - attachments: A slice of attachments (optional).
//
// Returns:
//   - *EmailMessage: A pointer to the newly created EmailMessage struct.
func NewFullEmailMessage(from string, to []string, subject string, cc []string, bcc []string, replyTo string, textBody string, htmlBody string, attachments []Attachment) *EmailMessage {
	return &EmailMessage{
		From:        from,
		To:          to,
		CC:          cc,
		BCC:         bcc,
		ReplyTo:     replyTo,
		Subject:     subject,
		Text:        textBody,
		HTML:        htmlBody,
		Attachments: attachments,
	}
}

// SetFrom sets the sender email address.
// Parameters:
// - from: The sender email address.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetFrom(from string) *EmailMessage {
	e.From = from
	return e
}

// SetSubject sets the email subject.
// Parameters:
// - subject: The email subject.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetSubject(subject string) *EmailMessage {
	e.Subject = subject
	return e
}

// SetTo sets the recipient email addresses.
// Parameters:
// - to: A slice of recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetTo(to []string) *EmailMessage {
	e.To = to
	return e
}

// SetCC sets the CC recipients email addresses.
// Parameters:
// - cc: A slice of CC recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetCC(cc []string) *EmailMessage {
	e.CC = cc
	return e
}

// SetBCC sets the BCC recipients email addresses.
// Parameters:
// - bcc: A slice of BCC recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetBCC(bcc []string) *EmailMessage {
	e.BCC = bcc
	return e
}

// SetReplyTo sets the reply-to email address.
// Parameters:
// - replyTo: The reply-to email address.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetReplyTo(replyTo string) *EmailMessage {
	e.ReplyTo = replyTo
	return e
}

// SetText sets the plain text content of the email.
// Parameters:
// - text: The plain text content of the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetText(text string) *EmailMessage {
	e.Text = text
	return e
}

// SetHTML sets the HTML content of the email.
// Parameters:
// - html: The HTML content of the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetHTML(html string) *EmailMessage {
	e.HTML = html
	return e
}

// SetAttachments sets the attachments for the email.
// Parameters:
// - attachments: A slice of Attachment structs to be included in the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetAttachments(attachments []Attachment) *EmailMessage {
	e.Attachments = attachments
	return e
}

// AddToRecipient adds a recipient email address to the To field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddToRecipient(recipient string) *EmailMessage {
	e.To = append(e.To, recipient)
	return e
}

// AddCCRecipient adds a recipient email address to the CC field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddCCRecipient(recipient string) *EmailMessage {
	e.CC = append(e.CC, recipient)
	return e
}

// AddBCCRecipient adds a recipient email address to the BCC field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddBCCRecipient(recipient string) *EmailMessage {
	e.BCC = append(e.BCC, recipient)
	return e
}

// AddAttachment adds an attachment to the email.
// Parameters:
// - attachment: An Attachment struct representing the file to be attached.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddAttachment(attachment Attachment) *EmailMessage {
	e.Attachments = append(e.Attachments, attachment)
	return e
}

// GetFrom returns the trimmed and validated sender email address.
// Returns an empty string if the email is invalid.
//
// Returns:
//   - string: The validated sender email address.
func (e *EmailMessage) GetFrom() string {
	if e == nil {
		return ""
	}
	return ValidateEmail(e.From)
}

// GetTo returns a slice of trimmed and validated recipient email addresses.
// Excludes any invalid email addresses.
//
// Returns:
//   - []string: The validated recipient email addresses.
func (e *EmailMessage) GetTo() []string {
	if e == nil {
		return []string{}
	}
	return ValidateEmailSlice(e.To)
}

// GetCC returns a slice of trimmed and validated CC recipient email addresses.
// Excludes any invalid email addresses.
//
// Returns:
//   - []string: The validated CC recipient email addresses.
func (e *EmailMessage) GetCC() []string {
	if e == nil {
		return []string{}
	}
	return ValidateEmailSlice(e.CC)
}

// GetBCC returns a slice of trimmed and validated BCC recipient email addresses.
// Excludes any invalid email addresses.
//
// Returns:
//   - []string: The validated BCC recipient email addresses.
func (e *EmailMessage) GetBCC() []string {
	if e == nil {
		return []string{}
	}
	return ValidateEmailSlice(e.BCC)
}

// GetReplyTo returns the trimmed and validated reply-to email address.
// Returns an empty string if the email is invalid.
//
// Returns:
//   - string: The validated reply-to email address.
func (e *EmailMessage) GetReplyTo() string {
	if e == nil {
		return ""
	}
	return ValidateEmail(e.ReplyTo)
}

// GetSubject returns the sanitized email subject.
// It escapes special characters like "<" to become "&lt;"
// If the EmailMessage is nil, it returns an empty string.
//
// Returns:
//   - string: The email subject.
func (e *EmailMessage) GetSubject() string {
	if e == nil {
		return ""
	}
	return SanitizeInput(e.Subject)
}

// GetText returns the sanitized plain text content of the email.
// It escapes special characters like "<" to become "&lt;"
// If the EmailMessage is nil, it returns an empty string.
//
// Returns:
//   - string: The plain text content of the email.
func (e *EmailMessage) GetText() string {
	if e == nil {
		return ""
	}
	return SanitizeInput(e.Text)
}

// GetHTML returns the sanitized HTML with only the UGC
// content of the email.
// If the EmailMessage is nil, it returns an empty string.
//
// Returns:
//   - string: The HTML sanitized content of the email.
func (e *EmailMessage) GetHTML() string {
	if e == nil {
		return ""
	}
	return bluemonday.UGCPolicy().Sanitize(e.HTML)
}

// GetAttachments returns the attachments to be included in the email.
// If the EmailMessage is nil, it returns an empty slice.
//
// Returns:
//   - []Attachment: The attachments to be included in the email.
func (e *EmailMessage) GetAttachments() []Attachment {
	if e == nil {
		return []Attachment{}
	}
	return e.Attachments
}

// Attachment represents an email attachment with its filename and content.
// Use this struct to specify files to be attached to the email.
type Attachment struct {
	Filename string // The name of the file.
	Content  []byte // The content of the file.
}

// GetBase64StringContent returns the content of the attachment as a base64-encoded string.
// If the attachment is nil, it returns an empty string.
//
// Returns:
//   - string: The base64-encoded content of the attachment as a string.
//     Returns an empty string if the attachment is nil.
func (a *Attachment) GetBase64StringContent() string {
	if a == nil {
		return ""
	}
	return string(a.GetBase64Content())
}

// GetFilename safely returns the filename of the attachment.
// It escapes special characters like "<" to become "&lt;"
// If the attachment is nil, it returns an empty string.
// Returns:
//   - string: The Filename as a string.
//     Returns an "nil_attachment" string if the attachment is nil.
func (a *Attachment) GetFilename() string {
	if a == nil {
		return "nil_attachment"
	}
	return SanitizeInput(a.Filename)
}

// GetBase64Content returns the content of the attachment as a base64-encoded byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The base64-encoded content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetBase64Content() []byte {
	if a == nil || len(a.Content) == 0 {
		return []byte{}
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(a.Content)))
	base64.StdEncoding.Encode(buf, a.Content)
	return buf
}

// GetRawContent returns the content of the attachment as its raw byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetRawContent() []byte {
	if a == nil || len(a.Content) == 0 {
		return []byte{}
	}
	return a.Content
}

// GetMimeType returns the MIME type based on the file extension.
// This function takes a filename, extracts its extension, and returns the corresponding MIME type.
//
// Parameters:
// - filename: A string containing the name of the file whose MIME type is to be determined.
//
// Returns:
// - string: The MIME type corresponding to the file extension.
//
// Example:
//
//	mimeType := GetMimeType("document.pdf")
//	fmt.Println(mimeType) // Output: application/pdf
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	return mime.TypeByExtension(ext)
}

// regex for validating email addresses
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9._\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail trims the email and checks if it is a valid email address.
// Returns the trimmed email if valid, otherwise returns an empty string.
func ValidateEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	if !emailRegex.MatchString(trimmed) {
		return ""
	}
	return trimmed
}

// ValidateEmailSlice trims and validates each email in the slice.
// Returns a slice of trimmed valid emails, excluding any invalid emails.
func ValidateEmailSlice(emails []string) []string {
	validEmails := []string{}
	for _, email := range emails {
		if validEmail := ValidateEmail(email); validEmail != "" {
			validEmails = append(validEmails, validEmail)
		}
	}
	return validEmails
}

// SanitizeInput sanitizes a string to prevent HTML and script injection.
func SanitizeInput(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

// isHTML checks if a string contains HTML tags.
// Parameters:
// - str: The string to check.
//
// Returns:
//   - bool: True if the string contains HTML tags, otherwise false.
func isHTML(str string) bool {
	htmlRegex := regexp.MustCompile(`(?i)<\/?[a-z][\s\S]*>`)
	return htmlRegex.MatchString(str)
}
