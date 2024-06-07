// OAuth2 credentials JSON (in production, retrieve from a secure file or secret manager)
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
attachmentContent, err := os.ReadFile("path```markdown
# Obtaining Gmail API Credentials for OAuth 2.0

To send emails using the Gmail API with OAuth 2.0, you need to set up OAuth2 credentials in the Google Cloud Console and obtain an access token.

## Steps

### 1. Create a project in Google Cloud Console:

- Go to [Google Cloud Console](https://console.cloud.google.com/).
- Create a new project.

### 2. Enable the Gmail API:

- In the Google Cloud Console, navigate to "APIs & Services" > "Library".
- Search for "Gmail API" and enable it for your project.

### 3. Set up OAuth 2.0 credentials:

- Go to "APIs & Services" > "Credentials".
- Click "Create credentials" and select "OAuth 2.0 Client IDs".
- Configure the consent screen and create OAuth 2.0 Client ID.
- Download the credentials JSON file.
- Gmail web credentials do not have the "redirect_uris" by default; however, Golang [oauth2/google](https://github.com/golang/oauth2/blob/0f29369cfe4552d0e4bcddc57cc75f4d7e672a33/google/google.go#L61) requires you to have an array of "redirect_uris":["http://redirect_uri",] to your credentials JSON. So add them to your OAuth 2.0 credentials.

### 4. Obtain an access token (if needed):

- Use the [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/) or a similar tool to obtain an access token using the credentials file.
- Save the access token in a secure place.

### 5. Run the example function to test the email sending:

```go
// OAuth2 credentials JSON (in production, retrieve from a secure file or secret manager)
credentials := []byte(`{
	"installed": {
		"client_id": "your-client-id",
		"project_id": "your-project-id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_secret": "your-client-secret",
		"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
	}
}`)

// Implementing the token manager
type MyTokenManager struct {
	config *oauth2.Config
}

func (m *MyTokenManager) GetToken() ([]byte, error) {
	// Implement your token retrieval logic here 
	// or request it from Google using the following code:

	token, err := loadTokenFromFile() //(in production, retrieve from a secure file or secret manager)
	if err != nil {
		return nil, fmt.Errorf("unable to load token from file: %v", err)
	}

	// If the token is expired, refresh it
	if !token.Valid() {
		tokenSource := m.config.TokenSource(context.Background(), token)
		token, err = tokenSource.Token()
		if err != nil {
			// Request user manual input if token refresh fails
			// Generate the consent URL
			authURL := m.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			fmt.Printf("- Go to the following link in your browser: \n%v\n", authURL)
			fmt.Println("- Once done, type the authorization code:")
			// Read the authorization code
			var authCode string
			if _, err := fmt.Scan(&authCode); err != nil {
				return nil, fmt.Errorf("unable to read authorization code: %v", err)
			}
			fmt.Printf("- Authorization received:\n- '%v'\n", authCode)
			// Exchange the authorization code for an access token
			token, err = m.config.Exchange(context.Background(), authCode)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
			}
		}

		// Save the refreshed token for future use
		//(in production, store in a secure file or secret manager)
		if err := saveTokenToFile(token); err != nil {
			return nil, err
		}
	}

	// Convert token to JSON
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal token to JSON: %v", err)
	}

	return tokenJSON, nil
}

config, err := credentials.ParseCredentials(credentials)
if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}
tokenManager := &MyTokenManager{
	config: config,
}

// Create GmailEmailSender with the parsed credentials and token
emailSender, err := gmail.NewGmailEmailSenderOauth2(credentials, tokenManager, "me")
if err != nil {
	log.Fatalf("Failed to create email sender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
	log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message using an alias email address
message := gomail.NewEmailMessage("your-email_or_alias@example.com",[]string{"recipient@example.com"}, "Test Email with attachment", "This is the plain text part of the email.").
		SetHTML("<p>This is the <b>HTML</b> part of the <i>email</i>.</p>").AddAttachments(*gomail.NewAttachment("attachment.jpg", attachmentContent))

// Send email
if err := emailSender.SendEmail(message); err != nil {
	log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")

```

Ensure your OAuth2 tokens are refreshed appropriately and stored securely.