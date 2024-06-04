# Obtaining SendGrid API Credentials

To send emails using the SendGrid API, you need to set up API credentials and obtain an API key.

## Steps

### 1. Sign up for SendGrid:
- Go to [SendGrid](https://sendgrid.com/).
- Sign up for an account if you don't already have one.

### 2. Create an API Key:
- Navigate to "Settings" > "API Keys".
- Click "Create API Key".
- Assign a name to the API key and give it the required permissions.
- Click "Create & View" and note down the generated API key.

### 3. Update the SendGrid settings in the `sendgrid_email_sender.go` file:
```go
apiKey := "your-sendgrid-api-key"
```
### 4. Run the example function to test the email sending:
```go
// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
apiKey := "your-sendgrid-api-key"

// Create SendGridEmailSender with the API key
emailSender, err := NewSendGridEmailSender(apiKey)
if err != nil {
    log.Fatalf("Failed to create email sender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message
message := *gomail.NewEmailMessage("your-email_or_alias@example.com",[]string{"recipient@example.com"}, "Test Email with attachment", "This is the plain text part of the email.").
		SetHTML("<p>This is the <b>HTML</b> part of the <i>email</i>.</p>").AddAttachments(gomail.Attachment{
			Filename: "attachment.jpg",  Content: attachmentContent,
		})

// Send email
if err := emailSender.SendEmail(message); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")
```
Ensure your API key is stored securely and refreshed appropriately.