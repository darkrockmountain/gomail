// Package providers provides functionality for sending emails using various providers.
//
// # Overview
//
// The providers package allows you to send emails using different email providers
// such as Gmail, SendGrid, AWS SES, and others. It abstracts the provider-specific
// details and provides a simple API for sending emails.
//
// # Usage
//
// To use the package, you need to create an instance of the email sender for your
// desired provider and then call the SendEmail function.
//
// Example:
//
//	package main
//
//	import (
//	    "github.com/darkrockmountain/gomail"
//	    "github.com/darkrockmountain/gomail/providers/sendgrid"
//	)
//
//	func main() {
//	    sender := sendgrid.NewSendGridEmailSender("your-api-key")
//	    err := sender.SendEmail(EmailMessage{To:[]string{"recipient@example.com"}, Subject:"Subject", Text:"Email body"})
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// This package supports various email providers and can be extended to include more.
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
package providers
