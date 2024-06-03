# Obtaining Microsoft Graph API Credentials for Application Authentication

To send emails using the Microsoft Graph API with application authentication, you need to set up Azure AD credentials and obtain an access token.

## Steps

### 1. Register an application in Azure AD:
- Go to [Azure Portal](https://portal.azure.com/).
- Navigate to "Azure Active Directory" > "App registrations".
- Click "New registration" and create a new application.

### 2. Configure API permissions:
- In the application registration, go to "API permissions".
- Add permissions for "Microsoft Graph" > "Application permissions" > "Mail.Send".

### 3. Set up client credentials:
- Go to "Certificates & secrets".
- Create a new client secret and note it down.

### 4. Find your Client ID and Tenant ID:
- In the application registration, go to the "Overview" section.
- Here you will find the "Application (client) ID" and the "Directory (tenant) ID".
- Note these down for use in your application.

### 5. Update the Microsoft Graph API settings in the `microsoft365_email_sender.go` file:
```go
clientId := "your-client-id"
clientSecret := "your-client-secret"
tenantId := "your-tenant-id"
```
### 6. Run the example function to test the email sending:
```go
// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
clientId := "your-client-id"
tenantId := "your-tenant-id"
clientSecret := "your-client-secret"

user := "user@domain.com"

// Create MSGraphEmailSender with the client credentials
emailSender, err := NewMSGraphEmailSenderAppAuth(clientId, tenantId, clientSecret, user)
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

Ensure your client secrets and tokens are stored securely and refreshed appropriately.