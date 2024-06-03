# Obtaining Mailgun API Credentials

To send emails using the Mailgun API, you need to set up API credentials and obtain an API key and domain.

## Steps

### 1. Sign up for Mailgun:
- Go to [Mailgun](https://www.mailgun.com/).
- Sign up for an account if you don't already have one.

### 2. Add and verify your domain:
- Navigate to "Sending" > "Domains".
- Click "Add New Domain".
- Follow the instructions to add and verify your domain.

### 3. Create an API Key:
- Navigate to "Settings" > "API Keys".
- Note down the "Private API Key".

### 4. Update the Mailgun settings in the `mailgun_email_sender.go` file:
```go
domain := "your-mailgun-domain"
apiKey := "your-mailgun-api-key"
```
### 5. Run the example function to test the email sending:
```go
// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
domain := "your-mailgun-domain"
apiKey := "your-mailgun-api-key"

// Create MailgunEmailSender with the domain and API key
emailSender, err := NewMailgunEmailSender(domain, apiKey)
if err != nil {
    log.Fatalf("Failed to create email sender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message
message := gomail.EmailMessage{
    From:        "your-email@example.com",
    To:          []string{"recipient@example.com"},
    Subject:     "Test Email with attachment",
    Text:        "This is the plain text part of the email.",
    HTML:        "<p>This is the <b>HTML</b> part of the <i>email</i>.</p>",
    Attachments: []gomail.Attachment{{Filename: "attachment.jpg", Content: attachmentContent}},
}

// Send email
if err := emailSender.SendEmail(message); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")
```

Ensure your API key and domain are stored securely and refreshed appropriately.