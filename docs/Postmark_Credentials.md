# Obtaining Postmark API Credentials

To send emails using the Postmark API, you need to set up API credentials and obtain a server token.

## Steps

### 1. Sign up for Postmark:
- Go to [Postmark](https://postmarkapp.com/).
- Sign up for an account if you don't already have one.

### 2. Create a new server:
- Navigate to "Servers".
- Click "Add Server".
- Assign a name to the server.

### 3. Obtain the Server API Token:
- In the server settings, go to the "API Tokens" tab.
- Note down the "Server API Token".

### 4. Update the Postmark settings in the `postmark_email_sender.go` file:
```go
serverToken := "your-postmark-server-token"
```
### 5. Run the example function to test the email sending:
```go
// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
serverToken := "your-postmark-server-token"

// Create PostmarkEmailSender with the server token
emailSender, err := NewPostmarkEmailSender(serverToken)
if err != nil {
    log.Fatalf("Failed to create email sender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message
message := gomail.NewEmailMessage("your-email_or_alias@example.com",[]string{"recipient@example.com"}, "Test Email with attachment", "This is the plain text part of the email.").
		SetHTML("<p>This is the <b>HTML</b> part of the <i>email</i>.</p>").AddAttachments(*gomail.NewAttachment("attachment.jpg", attachmentContent))

// Send email
if err := emailSender.SendEmail(message); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")

```
Ensure your server token is stored securely and refreshed appropriately.