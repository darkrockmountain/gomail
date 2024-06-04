# Obtaining SMTP Credentials

To send emails using the SMTP implementation, you need the following credentials from your SMTP server:

- SMTP Host (e.g., smtp.example.com)
- SMTP Port (e.g., 587)
- Username
- Password

## Steps

### 1. Contact your SMTP service provider to obtain the SMTP server details including host, port, username, and password.

### 2. Update the SMTP server settings in the `smtp_email_sender.go` file:
```go
host := "your-smtp-host"
port := 587 // or the port provided by your SMTP service
user := "your-username"
password := "your-password"
```

### 3. Run the example function to test the email sending:
```go
// SMTP server settings
host := "smtp.example.com"
port := 587
user := "your-username"
password := "your-password"
authMethod := AUTH_PLAIN

// Create SmtpEmailSender with the SMTP server settings
emailSender, _ := gmail.NewSmtpEmailSender(
    host,
    port,
    user,
    password,
    authMethod,
)

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message
message := gomail.EmailMessage{
    From:        "your-email_or_alias@example.com",
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

Make sure to handle your credentials securely and avoid hardcoding them directly in your source code for production use.
