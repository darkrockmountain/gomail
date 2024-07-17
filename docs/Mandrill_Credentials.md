# Obtaining Mandrill API Credentials

To send emails using the Mandrill API, you need to set up API credentials and obtain an API key.

## Steps

### 1. Sign up for Mandrill:
- Go to [Mandrill](https://mandrillapp.com/).
- Sign up for an account if you don't already have one.

### 2. Obtain an API Key:
- Navigate to "Settings" > "SMTP & API Info".
- Click "New API Key" to create a new API key.
- Note down the generated API key.

### 3. Update the Mandrill settings in the `mandrill_email_sender.go` file:
```go
apiKey := "your-mandrill-api-key"
```

### 4. Run the example function to test the email sending:
```go

// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
apiKey := "your-mandrill-api-key"

// Create MandrillEmailSender with the API key
emailSender, err := NewMandrillEmailSender(apiKey)
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
Ensure your API key is stored securely and refreshed appropriately.