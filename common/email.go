package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/darkrockmountain/gomail/sanitizer"
)

const DefaultMaxAttachmentSize = 25 * 1024 * 1024 // 25 MB

// EmailMessage contains the fields for sending an email.
// Use this struct to specify the sender, recipient, subject, and content of the email,
// as well as any attachments. This struct also supports custom sanitizers for text
// and HTML content to ensure that email content is safe and sanitized according to
// specific requirements.
//
// Example:
//
// To marshal an EmailMessage to JSON:
//
//	email := common.NewFullEmailMessage(
//	    "sender@example.com",
//	    []string{"recipient@example.com"},
//	    "Subject",
//	    []string{"cc@example.com"},
//	    []string{"bcc@example.com"},
//	    "replyto@example.com",
//	    "This is the email content.",
//	    "<p>This is the email content.</p>",
//	    []common.Attachment{*common.NewAttachment("attachment1.txt", []byte("file content"))})
//
//	jsonData, err := json.Marshal(email)
//	if err != nil {
//	    fmt.Println("Error marshaling to JSON:", err)
//	    return
//	}
//	fmt.Println("JSON output:", string(jsonData))
//
// To unmarshal an EmailMessage from JSON:
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
//	    "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}] // base64 encoded "file content"
//	}`
//	var email common.EmailMessage
//	err := json.Unmarshal([]byte(jsonData), &email)
//	if err != nil {
//	    fmt.Println("Error unmarshaling from JSON:", err)
//	    return
//	}
//	fmt.Printf("Unmarshaled EmailMessage: %+v\n", email)
type EmailMessage struct {
	from              string              // Sender email address.
	to                []string            // Recipient email addresses.
	cc                []string            // CC recipients email addresses.
	bcc               []string            // BCC recipients email addresses.
	replyTo           string              // Reply-To email address.
	subject           string              // Email subject.
	text              string              // Plain text content of the email.
	html              string              // HTML content of the email (optional).
	attachments       []Attachment        // Attachments to be included in the email (optional).
	maxAttachmentSize int                 // Maximum size for attachments.
	textSanitizer     sanitizer.Sanitizer // Sanitizer for plain text content.
	htmlSanitizer     sanitizer.Sanitizer // Sanitizer for HTML content.
}

// NewEmailMessage creates a new EmailMessage with the required fields.
// If the body contains HTML tags, it sets the HTML field; otherwise, it sets the Text field.
//
// Parameters:
//   - from: The sender email address.
//   - to: A slice of recipient email addresses.
//   - subject: The email subject.
//   - body: The content of the email, which can be plain text or HTML.
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

	if IsHTML(body) {
		email.html = body
	} else {
		email.text = body
	}

	return email
}

// NewFullEmailMessage creates a new EmailMessage with all fields.
//
// Parameters:
//   - from: The sender email address.
//   - to: A slice of recipient email addresses.
//   - subject: The email subject.
//   - cc: A slice of CC recipient email addresses (optional).
//   - bcc: A slice of BCC recipient email addresses (optional).
//   - replyTo: The reply-to email address (optional).
//   - textBody: The plain text content of the email.
//   - htmlBody: The HTML content of the email (optional).
//   - attachments: A slice of attachments (optional).
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
//
// Parameters:
//   - from: The sender email address.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetFrom(from string) *EmailMessage {
	e.from = from
	return e
}

// SetSubject sets the email subject.
//
// Parameters:
//   - subject: The email subject.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetSubject(subject string) *EmailMessage {
	e.subject = subject
	return e
}

// SetTo sets the recipient email addresses.
//
// Parameters:
//   - to: A slice of recipient email addresses.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetTo(to []string) *EmailMessage {
	e.to = to
	return e
}

// SetCC sets the CC recipients email addresses.
//
// Parameters:
//   - cc: A slice of CC recipient email addresses.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetCC(cc []string) *EmailMessage {
	e.cc = cc
	return e
}

// SetBCC sets the BCC recipients email addresses.
//
// Parameters:
//   - bcc: A slice of BCC recipient email addresses.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetBCC(bcc []string) *EmailMessage {
	e.bcc = bcc
	return e
}

// SetReplyTo sets the reply-to email address.
//
// Parameters:
//   - replyTo: The reply-to email address.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetReplyTo(replyTo string) *EmailMessage {
	e.replyTo = replyTo
	return e
}

// SetText sets the plain text content of the email.
//
// Parameters:
//   - text: The plain text content of the email.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetText(text string) *EmailMessage {
	e.text = text
	return e
}

// SetHTML sets the HTML content of the email.
//
// Parameters:
//   - html: The HTML content of the email.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetHTML(html string) *EmailMessage {
	e.html = html
	return e
}

// SetAttachments sets the attachments for the email.
//
// Parameters:
//   - attachments: A slice of Attachment structs to be included in the email.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetAttachments(attachments []Attachment) *EmailMessage {
	e.attachments = attachments
	return e
}

// AddToRecipient adds a recipient email address to the To field.
//
// Parameters:
//   - recipient: A recipient email address to be added.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddToRecipient(recipient string) *EmailMessage {
	e.to = append(e.to, recipient)
	return e
}

// AddCCRecipient adds a recipient email address to the CC field.
//
// Parameters:
//   - recipient: A recipient email address to be added.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddCCRecipient(recipient string) *EmailMessage {
	e.cc = append(e.cc, recipient)
	return e
}

// AddBCCRecipient adds a recipient email address to the BCC field.
//
// Parameters:
//   - recipient: A recipient email address to be added.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) AddBCCRecipient(recipient string) *EmailMessage {
	e.bcc = append(e.bcc, recipient)
	return e
}

// AddAttachment adds an attachment to the email.
//
// Parameters:
//   - attachment: An Attachment struct representing the file to be attached.
//
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
// It uses the custom text sanitizer if set, otherwise the default text sanitizer to escape special characters and trim whitespace.
// If the EmailMessage is nil, it returns an empty string.
//
// Returns:
//   - string: The sanitized email subject.
func (e *EmailMessage) GetSubject() string {
	if e == nil {
		return ""
	}
	if e.textSanitizer != nil {
		return e.textSanitizer.Sanitize(e.subject)
	}
	return sanitizer.DefaultTextSanitizer().Sanitize(e.subject)
}

// GetText returns the sanitized plain text content of the email.
// It uses the custom text sanitizer if set, otherwise the default sanitizer.
//
// Returns:
//   - string: The sanitized plain text content of the email.
func (e *EmailMessage) GetText() string {
	if e == nil {
		return ""
	}
	if e.textSanitizer != nil {
		return e.textSanitizer.Sanitize(e.text)
	}
	return sanitizer.DefaultTextSanitizer().Sanitize(e.text)
}

// GetHTML returns the sanitized HTML content of the email.
// It uses the custom html sanitizer if set, otherwise the default sanitizer.
//
// Returns:
//   - string: The sanitized HTML content of the email.
func (e *EmailMessage) GetHTML() string {
	if e == nil {
		return ""
	}
	if e.htmlSanitizer != nil {
		return e.htmlSanitizer.Sanitize(e.html)
	}
	return sanitizer.DefaultHtmlSanitizer().Sanitize(e.html)
}

// SetMaxAttachmentSize sets the maximum attachment size.
//
// Parameters:
//   - size: The maximum size for attachments in bytes. If set to a value less than 0, all attachments are allowed.
//
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

// SetCustomTextSanitizer sets a custom sanitizer for text content.
// WARNING: Using a custom text sanitizer may introduce security risks
// if the sanitizer does not properly handle potentially dangerous content.
// Ensure that the custom sanitizer is thoroughly tested and used with caution.
//
// Parameters:
//   - s: The custom Sanitizer implementation for text content.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetCustomTextSanitizer(s sanitizer.Sanitizer) *EmailMessage {
	e.textSanitizer = s
	return e
}

// SetCustomHtmlSanitizer sets a custom sanitizer for HTML content.
// WARNING: Using a custom HTML sanitizer may introduce security risks
// if the sanitizer does not properly handle potentially dangerous content.
// Ensure that the custom sanitizer is thoroughly tested and used with caution.
//
// Parameters:
//   - s: The custom Sanitizer implementation for HTML content.
//
// Returns:
//   - *EmailMessage: The EmailMessage struct pointer.
func (e *EmailMessage) SetCustomHtmlSanitizer(s sanitizer.Sanitizer) *EmailMessage {
	e.htmlSanitizer = s
	return e
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
//	email := common.NewFullEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Subject", []string{"cc@example.com"}, []string{"bcc@example.com"}, "replyto@example.com", "This is the email content.", "<p>This is the email content.</p>", []common.Attachment{*common.NewAttachment("attachment1.txt", []byte("file content"))})
//
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
//	    "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}] // base64 encoded "file content"
//	}`
//	var email common.EmailMessage
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

	//adding the default max attachment size by default
	e.maxAttachmentSize = DefaultMaxAttachmentSize

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

// BuildMimeMessage constructs the MIME message for the email, including text, HTML, and attachments.
// This function builds a multipart MIME message based on the provided email message. It supports plain text,
// HTML content, and multiple attachments.
//
// Parameters:
//   - message: A pointer to an EmailMessage struct containing the details of the email to be sent.
//
// Returns:
//   - []byte: A byte slice containing the complete MIME message.
//   - error: An error if constructing the MIME message fails, otherwise nil.
//
// Example:
//
//	message := gomail.NewEmailMessage(
//		"sender@example.com",
//		[]string["recipient@example.com"],
//		"Test Email",
//		"This is a test email.",)
//		.SetHtml("<p>This is a test email.</p>").AddAttachment(Attachment{
//		Filename: "test.txt",
//		Content:  []byte("This is a test attachment."),
//	})
//	mimeMessage, err := BuildMimeMessage(message)
//	if err != nil {
//	    log.Fatalf("Failed to build MIME message: %v", err)
//	}
//	fmt.Println(string(mimeMessage))
func BuildMimeMessage(message *EmailMessage) ([]byte, error) {
	var msg bytes.Buffer

	// Determine boundaries
	mixedBoundary := fmt.Sprintf("mixed-boundary-%d", time.Now().UnixNano())
	altBoundary := fmt.Sprintf("alt-boundary-%d", time.Now().UnixNano())

	// Basic headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", message.GetFrom()))

	// Add To recipients
	toRecipients := message.GetTo()
	if len(toRecipients) > 0 {
		msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(toRecipients, ",")))
	}

	ccRecipients := message.GetCC()

	if len(ccRecipients) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(ccRecipients, ",")))
	}

	if message.GetReplyTo() != "" {
		msg.WriteString(fmt.Sprintf("Reply-To: %s\r\n", message.GetReplyTo()))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", message.GetSubject()))

	msg.WriteString("MIME-Version: 1.0\r\n")

	// Use multipart/mixed if there are attachments, otherwise multipart/alternative
	attachments := message.GetAttachments()
	if len(attachments) > 0 {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", mixedBoundary))
		msg.WriteString("\r\n")
		// Start multipart/alternative
		msg.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", altBoundary))
		msg.WriteString("\r\n")
	} else {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", altBoundary))
		msg.WriteString("\r\n")
	}

	// Plain text part
	textMessage := message.GetText()
	if textMessage != "" {
		msg.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(textMessage)
		msg.WriteString("\r\n")
	}

	// HTML part
	htmlMessage := message.GetHTML()
	if htmlMessage != "" {
		msg.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(htmlMessage)
		msg.WriteString("\r\n")
	}

	// End multipart/alternative
	msg.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))

	// Attachments
	if len(attachments) > 0 {
		for _, attachment := range attachments {
			fileName := attachment.GetFilename()
			mimeType := GetMimeType(fileName)
			msg.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
			msg.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
			msg.WriteString("Content-Transfer-Encoding: base64\r\n")
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
			msg.WriteString("\r\n")
			msg.Write(attachment.GetBase64Content())
			msg.WriteString("\r\n")
		}

		// End multipart/mixed
		msg.WriteString(fmt.Sprintf("--%s--\r\n", mixedBoundary))
	}

	return msg.Bytes(), nil
}
