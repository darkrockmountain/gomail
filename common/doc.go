// Package common provides utility functions and structures used across the gomail project.
//
// # Overview
//
// The common package includes essential utilities such as email validation,
// sanitization, MIME type determination, and structures for handling email
// messages and attachments. These utilities are crucial for ensuring the
// integrity and security of email handling within the gomail project.
//
// # Components
//
//   - EmailMessage: Struct for constructing and manipulating email messages.
//   - Attachment: Struct for managing email attachments, including file handling and base64 encoding.
//   - Validation: Functions for validating email addresses and slices of email addresses.
//   - Sanitization: Functions for sanitizing input to prevent injection attacks.
//
// # Usage
//
// The functions and structs in this package are used internally by other packages
// in the gomail project but can also be used directly if needed.
//
// Example:
//
//	package main
//
//	import (
//	    "github.com/darkrockmountain/gomail/common"
//	)
//
//	func main() {
//	    email := common.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Subject", "Email body")
//	    fmt.Println(email.GetSubject())
//	}
//
// For detailed information on each function and struct, refer to their respective
// documentation within this package.
package common
