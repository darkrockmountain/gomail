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

// ---- EmailMessage and related functions ----
type EmailMessage = common.EmailMessage

var NewEmailMessage = common.NewEmailMessage
var BuildMimeMessage = common.BuildMimeMessage
var NewFullEmailMessage = common.NewFullEmailMessage

// ---- Attachment and related functions ----
type Attachment = common.Attachment

var NewAttachment = common.NewAttachment
var NewAttachmentFromFile = common.NewAttachmentFromFile

// ---- Validation functions ----
var ValidateEmail = common.ValidateEmail
var ValidateEmailSlice = common.ValidateEmailSlice

// ---- Utility functions ----
var GetMimeType = common.GetMimeType
