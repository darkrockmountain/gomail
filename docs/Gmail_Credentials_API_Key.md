# Obtaining Gmail API Credentials for API Key

To send emails using the Gmail API with an API key, you need to set up an API key in the Google Cloud Console.

## Steps

### 1. Create a project in Google Cloud Console:

- Go to [Google Cloud Console](https://console.cloud.google.com/).
- Create a new project.

### 2. Enable the Gmail API:

- In the Google Cloud Console, navigate to "APIs & Services" > "Library".
- Search for "Gmail API" and enable it for your project.

### 3. Create an API Key:

- Go to "APIs & Services" > "Credentials".
- Click "Create credentials" and select "API Key".
- Copy the generated API key and store it securely.

### 4. Run the example function to test the email sending:


```go
// API key (in production, retrieve from a secure file or secret manager)
apiKey := "your-api-key"

// User for whom the email is being sent
user := "me"

// Initialize the GmailEmailSender
emailSender, err := gmail.NewGmailEmailSenderAPIKey(apiKey, user)
if err != nil {
    log.Fatalf("Failed to create GmailEmailSender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message using an alias email address
message := gomail.NewEmailMessage("your-email_or_alias@example.com",[]string{"recipient@example.com"}, "Test Email with attachment", "This is the plain text part of the email.").
		SetHTML("<p>This is the <b>HTML</b> part of the <i>email</i>.</p>").AddAttachments(gomail.Attachment{
			Filename: "attachment.jpg",  Content: attachmentContent,
		})



// Send email
if err := emailSender.SendEmail(message); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")
```

