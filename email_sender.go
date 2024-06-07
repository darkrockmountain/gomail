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
	"encoding/json"
	"html"
	"io/ioutil"
	"mime"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const DefaultMaxAttachmentSize = 25 * 1024 * 1024 // 25 MB

// EmailSender interface defines the method to send an email.
// Implement this interface to create different email sending strategies.
type EmailSender interface {
	// SendEmail sends an email with the given message.
	// Parameters:
	// - message: A pointer to an EmailMessage struct containing the details of the email to be sent.
	// Returns:
	// - error: An error if sending the email fails, otherwise nil.
	SendEmail(message *EmailMessage) error
}

// EmailMessage contains the fields for sending an email.
// Use this struct to specify the sender, recipient, subject, and content of the email,
// as well as any attachments.
type EmailMessage struct {
	from              string       // Sender email address.
	to                []string     // Recipient email addresses.
	cc                []string     // CC recipients email addresses.
	bcc               []string     // BCC recipients email addresses.
	replyTo           string       // Reply-To email address.
	subject           string       // Email subject.
	text              string       // Plain text content of the email.
	html              string       // HTML content of the email (optional).
	attachments       []Attachment // Attachments to be included in the email (optional).
	maxAttachmentSize int          // Maximum size for attachments.

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
		from:              from,
		to:                to,
		subject:           subject,
		maxAttachmentSize: DefaultMaxAttachmentSize,
	}

	if isHTML(body) {
		email.html = body
	} else {
		email.text = body
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
		from:              from,
		to:                to,
		cc:                cc,
		bcc:               bcc,
		replyTo:           replyTo,
		subject:           subject,
		text:              textBody,
		html:              htmlBody,
		attachments:       attachments,
		maxAttachmentSize: DefaultMaxAttachmentSize,
	}
}

// SetFrom sets the sender email address.
// Parameters:
// - from: The sender email address.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetFrom(from string) *EmailMessage {
	e.from = from
	return e
}

// SetSubject sets the email subject.
// Parameters:
// - subject: The email subject.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetSubject(subject string) *EmailMessage {
	e.subject = subject
	return e
}

// SetTo sets the recipient email addresses.
// Parameters:
// - to: A slice of recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetTo(to []string) *EmailMessage {
	e.to = to
	return e
}

// SetCC sets the CC recipients email addresses.
// Parameters:
// - cc: A slice of CC recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetCC(cc []string) *EmailMessage {
	e.cc = cc
	return e
}

// SetBCC sets the BCC recipients email addresses.
// Parameters:
// - bcc: A slice of BCC recipient email addresses.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetBCC(bcc []string) *EmailMessage {
	e.bcc = bcc
	return e
}

// SetReplyTo sets the reply-to email address.
// Parameters:
// - replyTo: The reply-to email address.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetReplyTo(replyTo string) *EmailMessage {
	e.replyTo = replyTo
	return e
}

// SetText sets the plain text content of the email.
// Parameters:
// - text: The plain text content of the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetText(text string) *EmailMessage {
	e.text = text
	return e
}

// SetHTML sets the HTML content of the email.
// Parameters:
// - html: The HTML content of the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetHTML(html string) *EmailMessage {
	e.html = html
	return e
}

// SetAttachments sets the attachments for the email.
// Parameters:
// - attachments: A slice of Attachment structs to be included in the email.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetAttachments(attachments []Attachment) *EmailMessage {
	e.attachments = attachments
	return e
}

// AddToRecipient adds a recipient email address to the To field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddToRecipient(recipient string) *EmailMessage {
	e.to = append(e.to, recipient)
	return e
}

// AddCCRecipient adds a recipient email address to the CC field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddCCRecipient(recipient string) *EmailMessage {
	e.cc = append(e.cc, recipient)
	return e
}

// AddBCCRecipient adds a recipient email address to the BCC field.
// Parameters:
// - recipient: A recipient email address to be added.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddBCCRecipient(recipient string) *EmailMessage {
	e.bcc = append(e.bcc, recipient)
	return e
}

// AddAttachment adds an attachment to the email.
// Parameters:
// - attachment: An Attachment struct representing the file to be attached.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddAttachment(attachment Attachment) *EmailMessage {
	e.attachments = append(e.attachments, attachment)
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
	return ValidateEmail(e.from)
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
	return ValidateEmailSlice(e.to)
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
	return ValidateEmailSlice(e.cc)
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
	return ValidateEmailSlice(e.bcc)
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
	return ValidateEmail(e.replyTo)
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
	return SanitizeInput(e.subject)
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
	return SanitizeInput(e.text)
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
	return bluemonday.UGCPolicy().Sanitize(e.html)
}

// SetMaxAttachmentSize sets the maximum attachment size.
// Parameters:
// - size: The maximum size for attachments in bytes. If set to a value less than 0, all attachments are allowed.
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetMaxAttachmentSize(size int) *EmailMessage {
	e.maxAttachmentSize = size
	return e
}

// GetAttachments returns the attachments to be included in the email,
// filtering out those that exceed the maximum size.
// If the EmailMessage is nil, it returns an empty slice.
//
// Returns:
//   - []Attachment: The attachments to be included in the email.
func (e *EmailMessage) GetAttachments() []Attachment {
	if e == nil {
		return []Attachment{}
	}
	// if maxAttachmentSize return the attachments without the limit
	if e.maxAttachmentSize < 0 {
		return e.attachments
	}
	// return only the attachments withing the max size limit
	var validAttachments []Attachment
	for _, attachment := range e.attachments {
		if len(attachment.content) <= e.maxAttachmentSize {
			validAttachments = append(validAttachments, attachment)
		}
	}
	return validAttachments
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

// Attachment represents an email attachment with its filename and content.
// Use this struct to specify files to be attached to the email.
type Attachment struct {
	filename string // The name of the file.
	content  []byte // The content of the file.
}

// NewAttachment creates a new Attachment instance with the specified filename and content.
// It initializes the private fields of the Attachment struct with the provided values.
//
// Example:
//
//	content := []byte("file content")
//	attachment := NewAttachment("document.pdf", content)
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
//	fmt.Println(string(attachment.GetContent())) // Output: file content
func NewAttachment(filename string, content []byte) *Attachment {
	return &Attachment{
		filename: filename,
		content:  content,
	}
}

// NewAttachmentFromFile creates a new Attachment instance from the specified file path.
// It reads the content of the file and initializes the private fields of the Attachment struct.
//
// Example:
//
//	attachment, err := NewAttachmentFromFile("path/to/document.pdf")
//	if err != nil {
//	    fmt.Println("Error creating attachment:", err)
//	    return
//	}
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
//	fmt.Println(string(attachment.GetContent())) // Output: (file content)
func NewAttachmentFromFile(filePath string) (*Attachment, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	filename := extractFilename(filePath)
	return NewAttachment(
		filename,
		content,
	), nil
}

// extractFilename extracts the filename from the file path.
// This is a helper function to get the filename from a given file path.
func extractFilename(filePath string) string {
	// Implement this function based on your needs, for simplicity using base method
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// SetFilename sets the filename of the attachment.
// It assigns the provided filename to the private filename field.
//
// Example:
//
//	var attachment Attachment
//	attachment.SetFilename("document.pdf")
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
func (a *Attachment) SetFilename(filename string) {
	a.filename = filename
}

// GetFilename safely returns the filename of the attachment.
// It escapes special characters like "<" to become "&lt;"
// If the attachment is nil, it returns an empty string.
// Returns:
//   - string: The filename as a string.
//     Returns an "nil_attachment" string if the attachment is nil.
func (a *Attachment) GetFilename() string {
	if a == nil {
		return "nil_attachment"
	}
	return SanitizeInput(a.filename)
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

// SetContent sets the content of the attachment.
// It assigns the provided content to the private content field.
//
// Example:
//
//	var attachment Attachment
//	content := []byte("file content")
//	attachment.SetContent(content)
//	fmt.Println(string(attachment.GetContent())) // Output: file content
func (a *Attachment) SetContent(content []byte) {
	a.content = content
}

// GetBase64Content returns the content of the attachment as a base64-encoded byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The base64-encoded content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetBase64Content() []byte {
	if a == nil || len(a.content) == 0 {
		return []byte{}
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(a.content)))
	base64.StdEncoding.Encode(buf, a.content)
	return buf
}

// GetRawContent returns the content of the attachment as its raw byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetRawContent() []byte {
	if a == nil || len(a.content) == 0 {
		return []byte{}
	}
	return a.content
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

// jsonEmailMessage represents the JSON structure for an email message.
type jsonEmailMessage struct {
	From        string       `json:"from"`
	To          []string     `json:"to"`
	CC          []string     `json:"cc,omitempty"`
	BCC         []string     `json:"bcc,omitempty"`
	ReplyTo     string       `json:"replyTo,omitempty"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	HTML        string       `json:"html,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// MarshalJSON custom marshaler for EmailMessage
// This method converts the EmailMessage struct into a JSON representation.
// It creates an anonymous struct with exported fields and JSON tags,
// copies the values from the private fields, and then marshals it to JSON.
//
// Example:
//
//	email := &EmailMessage{
//	    from:    "sender@example.com",
//	    to:      []string{"recipient@example.com"},
//	    cc:      []string{"cc@example.com"},
//	    bcc:     []string{"bcc@example.com"},
//	    replyTo: "replyto@example.com",
//	    subject: "Subject",
//	    text:    "This is the email content.",
//	    html:    "<p>This is the email content.</p>",
//	    attachments: []Attachment{
//	        {filename: "attachment1.txt", content: []byte("content1")},
//	    },
//	    maxAttachmentSize: 1024,
//	}
//	jsonData, err := json.Marshal(email)
//	if err != nil {
//	    fmt.Println("Error marshaling to JSON:", err)
//	    return
//	}
//	fmt.Println("JSON output:", string(jsonData))
func (e *EmailMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonEmailMessage{
		From:        e.from,
		To:          e.to,
		CC:          e.cc,
		BCC:         e.bcc,
		ReplyTo:     e.replyTo,
		Subject:     e.subject,
		Text:        e.text,
		HTML:        e.html,
		Attachments: e.attachments,
	})
}

// UnmarshalJSON custom unmarshaler for EmailMessage
// This method converts a JSON representation into an EmailMessage struct.
// It creates an anonymous struct with exported fields and JSON tags,
// unmarshals the JSON data into this struct, and then copies the values
// to the private fields of the EmailMessage struct.
//
// Example:
//
//	jsonData := `{
//	    "from": "sender@example.com",
//	    "to": ["recipient@example.com"],
//	    "cc": ["cc@example.com"],
//	    "bcc": ["bcc@example.com"],
//	    "replyTo": "replyto@example.com",
//	    "subject": "Subject",
//	    "text": "This is the email content.",
//	    "html": "<p>This is the email content.</p>",
//	    "attachments": [{"filename": "attachment1.txt", "content": "Y29udGVudDE="}] // base64 encoded "content1"
//	}`
//	var email EmailMessage
//	err := json.Unmarshal([]byte(jsonData), &email)
//	if err != nil {
//	    fmt.Println("Error unmarshaling from JSON:", err)
//	    return
//	}
//	fmt.Printf("Unmarshaled EmailMessage: %+v\n", email)
func (e *EmailMessage) UnmarshalJSON(data []byte) error {
	aux := &jsonEmailMessage{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	e.from = aux.From
	e.to = aux.To
	e.cc = aux.CC
	e.bcc = aux.BCC
	e.replyTo = aux.ReplyTo
	e.subject = aux.Subject
	e.text = aux.Text
	e.html = aux.HTML
	e.attachments = aux.Attachments

	return nil
}

// jsonAttachment represents the JSON structure for an email attachment.
type jsonAttachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"` // Content will be base64 encoded
}

// MarshalJSON custom marshaler for Attachment
// This method converts the Attachment struct into a JSON representation.
// It creates an anonymous struct with exported fields and JSON tags,
// copies the values from the private fields, and then marshals it to JSON.
//
// Example:
//
//	attachment := Attachment{
//	    filename: "file.txt",
//	    content:  []byte("file content"),
//	}
//	jsonData, err := json.Marshal(attachment)
//	if err != nil {
//	    fmt.Println("Error marshaling to JSON:", err)
//	    return
//	}
//	fmt.Println("JSON output:", string(jsonData))
func (a Attachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAttachment{
		Filename: a.filename,
		Content:  base64.StdEncoding.EncodeToString(a.content), // Encode content to base64
	})
}

// UnmarshalJSON custom unmarshaler for Attachment
// This method converts a JSON representation into an Attachment struct.
// It creates an anonymous struct with exported fields and JSON tags,
// unmarshals the JSON data into this struct, and then copies the values
// to the private fields of the Attachment struct.
//
// Example:
//
//	jsonData := `{
//	    "filename": "file.txt",
//	    "content": "ZmlsZSBjb250ZW50" // base64 encoded "file content"
//	}`
//	var attachment Attachment
//	err := json.Unmarshal([]byte(jsonData), &attachment)
//	if err != nil {
//	    fmt.Println("Error unmarshaling from JSON:", err)
//	    return
//	}
//	fmt.Printf("Unmarshaled Attachment: %+v\n", attachment)
func (a *Attachment) UnmarshalJSON(data []byte) error {
	aux := &jsonAttachment{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	a.filename = aux.Filename
	content, err := base64.StdEncoding.DecodeString(aux.Content) // Decode content from base64
	if err != nil {
		return err
	}
	a.content = content

	return nil
}
