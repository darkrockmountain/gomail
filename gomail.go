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
//   - providers: Contains implementations for various email providers.
//   - credentials: Contains implementations for managing email credentials.
//   - examples: Contains example applications demonstrating how to use the library.
//   - docs: Contains documentation for configuring different email providers.
//   - common: Contains shared utilities and types used across the project.
//   - sanitizer: Contains implementations for sanitizing email content.
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
//	    err := sender.SendEmail(gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Subject", "Email body"))
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// This library supports various email providers and can be extended to include more.
//
// # Supported Providers
//
//   - Gmail
//   - SendGrid
//   - AWS SES
//   - Mailgun
//   - Mandrill
//   - Postmark
//   - Microsoft365
//   - SparkPost
//   - SMTP
//
// For more details, see the documentation for each provider in the providers package.
package gomail

import "github.com/darkrockmountain/gomail/common"

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

// EmailMessage represents an email message with various fields such as sender, recipients, subject, and content.
type EmailMessage = common.EmailMessage

// Attachment represents an email attachment with its filename and content.
type Attachment = common.Attachment

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
	return common.NewEmailMessage(from, to, subject, body)
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
	return common.NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, textBody, htmlBody, attachments)
}

// NewAttachment creates a new Attachment instance with the specified filename and content.
// It initializes the private fields of the Attachment struct with the provided values.
//
// Parameters:
//   - filename: The name of the file to be attached.
//   - content: The content of the file as a byte slice.
//
// Returns:
//   - *Attachment: A pointer to the newly created Attachment struct.
func NewAttachment(filename string, content []byte) *Attachment {
	return common.NewAttachment(filename, content)
}

// NewAttachmentFromFile creates a new Attachment instance from the specified file path.
// It reads the content of the file and initializes the private fields of the Attachment struct.
//
// Parameters:
//   - filePath: The path to the file to be attached.
//
// Returns:
//   - *Attachment: A pointer to the newly created Attachment struct.
//   - error: An error if reading the file fails, otherwise nil.
func NewAttachmentFromFile(filePath string) (*Attachment, error) {
	return common.NewAttachmentFromFile(filePath)
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
func BuildMimeMessage(message *EmailMessage) ([]byte, error) {
	return common.BuildMimeMessage(message)
}

// ValidateEmail validates and sanitizes an email address.
//
// Parameters:
//   - email: The email address to be validated and sanitized.
//
// Returns:
//   - string: The validated and sanitized email address, or an empty string if invalid.
func ValidateEmail(email string) string {
	return common.ValidateEmail(email)
}

// ValidateEmailSlice validates and sanitizes a slice of email addresses.
//
// Parameters:
//   - emails: A slice of email addresses to be validated and sanitized.
//
// Returns:
//   - []string: A slice of validated and sanitized email addresses, excluding any invalid addresses.
func ValidateEmailSlice(emails []string) []string {
	return common.ValidateEmailSlice(emails)
}

// GetMimeType returns the MIME type based on the file extension.
// This function takes a filename, extracts its extension, and returns the corresponding MIME type.
//
// Parameters:
//   - filename: A string containing the name of the file whose MIME type is to be determined.
//
// Returns:
//   - string: The MIME type corresponding to the file extension.
func GetMimeType(filename string) string {
	return common.GetMimeType(filename)
}
