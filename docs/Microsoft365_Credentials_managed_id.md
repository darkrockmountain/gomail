# Obtaining Microsoft Graph API Credentials for Managed Identity Authentication

To send emails using the Microsoft Graph API with managed identity authentication, you need to set up Azure AD managed identity credentials.

## Steps

### 1. Register an application in Azure AD:
- Go to [Azure Portal](https://portal.azure.com/).
- Navigate to "Azure Active Directory" > "App registrations".
- Click "New registration" and create a new application.

### 2. Configure API permissions:
- In the application registration, go to "API permissions".
- Add permissions for "Microsoft Graph" > "Application permissions" > "Mail.Send".

### 3. Set up managed identity credentials:
- Ensure your Azure environment is configured to support managed identities.
- Assign the managed identity to your Azure resource (e.g., VM, App Service).

### 4. Update the Microsoft Graph API settings in the `microsoft365_email_sender.go` file:
```go
user := "user@domain.com"
```
### 5. Run the example function to test the email sending:
```go
user := "user@domain.com"

// Create MSGraphEmailSender with managed identity credentials
emailSender, err := NewMSGraphEmailSenderManagedIdentity(user)
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
Ensure your managed identity credentials are configured securely.