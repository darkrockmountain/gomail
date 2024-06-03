package providers

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/darkrockmountain/gomail"
)

// buildMimeMessage constructs the MIME message for the email, including text, HTML, and attachments.
// This function builds a multipart MIME message based on the provided email message. It supports plain text,
// HTML content, and multiple attachments.
//
// Parameters:
// - message: An EmailMessage struct containing the details of the email to be sent.
//
// Returns:
// - []byte: A byte slice containing the complete MIME message.
// - error: An error if constructing the MIME message fails, otherwise nil.
//
// Example:
//
//	message := EmailMessage{
//	    From:    "sender@example.com",
//	    To:      ["recipient@example.com"],
//	    Subject: "Test Email",
//	    Text:    "This is a test email.",
//	    HTML:    "<p>This is a test email.</p>",
//	    Attachments: []Attachment{
//	        {
//	            Filename: "test.txt",
//	            Content:  []byte("This is a test attachment."),
//	        },
//	    },
//	}
//	mimeMessage, err := buildMimeMessage(message)
//	if err != nil {
//	    log.Fatalf("Failed to build MIME message: %v", err)
//	}
//	fmt.Println(string(mimeMessage))
func buildMimeMessage(message gomail.EmailMessage) ([]byte, error) {
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
			mimeType := gomail.GetMimeType(fileName)
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

// strPtr takes a string value and returns a pointer to that string.
// This function is useful when you need to work with string pointers, such as in
// scenarios where you need to pass a string by reference or handle optional string fields.
//
// Parameters:
//   - str (string): The input string value that you want to convert to a pointer.
//
// Returns:
//   - *string: A pointer to the input string value.
//
// Example usage:
//
//	name := "John Doe"
//	namePtr := strPtr(name)
//	fmt.Println(namePtr)  // Output: memory address of the string
//	fmt.Println(*namePtr) // Output: "John Doe"
//
// Detailed explanation:
// The strPtr function creates a pointer to the given string `str`.
// This can be particularly useful in the following scenarios:
//  1. Passing strings by reference to functions, which can help avoid copying large strings.
//  2. Working with data structures that use pointers to represent optional fields or nullable strings.
//  3. Interfacing with APIs or libraries that require or return string pointers.
//
// By using this function, you can easily obtain a pointer to a string and utilize it in contexts
// where pointers are needed, thus enhancing flexibility and efficiency in your Go programs.
func strPtr(str string) *string {
	return &str
}
