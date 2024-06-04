# Obtaining Gmail API Credentials for Service Accounts

To send emails using the Gmail API with a service account, you need to set up service account credentials in the Google Cloud Console and obtain the service account key file.

## Steps

### 1. Create a project in Google Cloud Console:
- Go to [Google Cloud Console](https://console.cloud.google.com/).
- Create a new project.

### 2. Enable the Gmail API:
- In the Google Cloud Console, navigate to "APIs & Services" > "Library".
- Search for "Gmail API" and enable it for your project.

### 3. Set up a service account:
- Go to "IAM & Admin" > "Service Accounts".
- Click "Create Service Account".
- Provide a name and description for the service account, and click "Create".
- Assign the role "Project > Editor" to the service account, and click "Continue".
- Click "Done" to finish creating the service account.

### 4. Create and download a service account key:
- In the "Service Accounts" page, find the created service account.
- Click the three dots on the right, and select "Manage keys".
- Click "Add key" > "Create new key".
- Select "JSON" as the key type, and click "Create".
- Download the JSON key file to a secure location.

### 5. Set up domain-wide delegation (if necessary):

#### Steps in Google Cloud Console:
- In the Google Cloud Console, go to "IAM & Admin" > "Service Accounts".
- Click on the service account you created.
- Click "Edit" and then "Show Domain-wide Delegation".
- Enable "Enable G Suite Domain-wide Delegation" and provide a product name for the consent screen.
- Save your changes.
- Click "View Client ID" and note the "Client ID" for the next steps.

#### Steps in Google Workspace Admin Console:
- Go to your [Google Workspace Admin Console](https://admin.google.com).
- Navigate to "Security" > "API controls" > "Manage domain-wide delegation".
- Click "Add new" and enter the "Client ID" noted earlier.
- In the "OAuth Scopes" field, enter the [scopes](https://developers.google.com/identity/protocols/oauth2/scopes) required for the Gmail API: `https://www.googleapis.com/auth/gmail.send`
- Click "Authorize" to save the settings.

### 6. Run the example function to test the email sending:

```go
// Service account credentials JSON (in production, retrieve from a secure file or secret manager)
jsonCredentials := []byte(`{
    "type": "service_account",
    "project_id": "your-project-id",
    "private_key_id": "your-private-key-id",
    "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----\n",
    "client_email": "your-service-account-email@your-project.iam.gserviceaccount.com",
    "client_id": "your-client-id",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "your-client-secret",
    "redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
}`)

// User for whom the email is being sent
user := "user@domain.com"

// Initialize the GmailEmailSender
emailSender, err := gmail.NewGmailEmailSenderServiceAccount(jsonCredentials, user)
if err != nil {
    log.Fatalf("Failed to create GmailEmailSender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message using an alias email address
message := gomail.EmailMessage{
    From:        "your-email_or_alias@example.com",
    To:          []string{"recipient@example.com"},
    Subject:     "Test Email with Alias",
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

Following these steps ensures that your service account is correctly set up to send emails on behalf of users within your Google Workspace domain.
