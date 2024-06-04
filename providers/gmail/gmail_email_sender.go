package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/credentials"
	"github.com/darkrockmountain/gomail/providers"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// gmailMessageSenderWrapper wraps the Gmail UsersMessagesService.
type gmailMessageSenderWrapper struct {
	messageSender *gmail.UsersMessagesService
	user          string
}

// send sends a Gmail message.
// Parameters:
//   - message: The Gmail message to be sent.
//
// Returns:
//   - *gmail.Message: The sent Gmail message.
//   - error: An error if sending the message fails.
func (ms *gmailMessageSenderWrapper) send(message *gmail.Message) (*gmail.Message, error) {
	if ms.messageSender == nil {
		return nil, errors.New("error no UsersMessagesService initiated")
	}
	user := ms.user
	if user == "" {
		user = "me"
	}
	return ms.messageSender.Send(user, message).Do()
}

// SendEmail sends an email using the Gmail API.
// It constructs the MIME message from the given EmailMessage and sends it using the Gmail API.
// Parameters:
//   - message: The EmailMessage containing the email details.
//
// Returns:
//   - error: An error if sending the email fails.
func (s *gmailMessageSenderWrapper) SendEmail(message gomail.EmailMessage) error {
	mimeMessage, err := providers.BuildMimeMessage(message)
	if err != nil {
		return fmt.Errorf("unable to build MIME message: %w", err)
	}

	bccs := message.GetBCC()
	if len(bccs) > 0 {
		var msg bytes.Buffer
		msg.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(bccs, ",")))
		mimeMessage = append(msg.Bytes(), mimeMessage...)
	}

	gMessage := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(mimeMessage),
	}

	_, err = s.send(gMessage)
	if err != nil {
		return fmt.Errorf("unable to send message: %v", err)
	}

	return nil
}

// GmailTokenManager defines an interface for obtaining OAuth2 tokens.
type GmailTokenManager interface {
	GetToken() ([]byte, error)
}

// NewGmailEmailSenderOauth2 initializes a new gmailEmailSenderOauth2 instance using OAuth2 credentials and token management.
// Parameters:
//   - ctx: The context for the API requests.
//   - configJson: A byte slice containing the OAuth2 client credentials JSON.
//   - tokenManager: An implementation of the GmailTokenManager interface.
//   - user: The user for whom the email is being sent (usually "me").
//
// Returns:
//   - *gmailMessageSenderWrapper: A pointer to the initialized gmailMessageSenderWrapper.
//   - error: An error if the initialization fails.
func NewGmailEmailSenderOauth2(ctx context.Context, configJson []byte, tokenManager GmailTokenManager, user string) (*gmailMessageSenderWrapper, error) {
	config, err := credentials.ParseCredentials(configJson)
	if err != nil {
		return nil, err
	}

	tokBytes, err := tokenManager.GetToken()
	if err != nil {
		return nil, err
	}
	tok, err := credentials.ParseToken(tokBytes)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to start Gmail service: %v", err)
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderServiceAccount initializes a new gmailEmailSenderServiceAccount instance using service account JSON credentials.
// Parameters:
//   - ctx: The context for the API requests.
//   - jsonCredentials: A byte slice containing the JSON service account credentials.
//   - user: The user for whom the email is being sent.
//
// Returns:
//   - *gmailMessageSenderWrapper: A pointer to the initialized gmailMessageSenderWrapper.
//   - error: An error if the initialization fails.
func NewGmailEmailSenderServiceAccount(ctx context.Context, jsonCredentials []byte, user string) (*gmailMessageSenderWrapper, error) {
	params := google.CredentialsParams{
		Scopes:  []string{gmail.GmailSendScope},
		Subject: user,
	}
	creds, err := google.CredentialsFromJSONWithParams(ctx, jsonCredentials, params)
	if err != nil {
		return nil, fmt.Errorf("unable to parse service account key file: %w", err)
	}
	srv, err := gmail.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %w", err)
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderAPIKey initializes a new gmailEmailSenderAPIKey instance using an API key.
// Parameters:
//   - ctx: The context for the API requests.
//   - apiKey: The API key to be used for authentication.
//   - user: The user for whom the email is being sent (usually "me").
//
// Returns:
//   - *gmailMessageSenderWrapper: A pointer to the initialized gmailMessageSenderWrapper.
//   - error: An error if the initialization fails.
func NewGmailEmailSenderAPIKey(ctx context.Context, apiKey, user string) (*gmailMessageSenderWrapper, error) {
	srv, err := gmail.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderJWT initializes a new gmailEmailSenderJWT instance using JWT configuration.
// Parameters:
//   - ctx: The context for the API requests.
//   - configJson: A byte slice containing the JWT client credentials JSON.
//   - user: The user for whom the email is being sent (usually "me").
//
// Returns:
//   - *gmailMessageSenderWrapper: A pointer to the initialized gmailMessageSenderWrapper.
//   - error: An error if the initialization fails.
func NewGmailEmailSenderJWT(ctx context.Context, configJson []byte, user string) (*gmailMessageSenderWrapper, error) {
	config, err := google.JWTConfigFromJSON(configJson)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JWT credentials: %v", err)
	}

	client := config.Client(ctx)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to start Gmail service: %v", err)
	}
	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderJWTAccess initializes a new gmailEmailSenderJWTAccess instance using a JWT access token.
// Parameters:
//   - ctx: The context for the API requests.
//   - jsonCredentials: A byte slice containing the JSON service account credentials.
//   - user: The user for whom the email is being sent.
//
// Returns:
//   - *gmailMessageSenderWrapper: A pointer to the initialized gmailMessageSenderWrapper.
//   - error: An error if the initialization fails.
func NewGmailEmailSenderJWTAccess(ctx context.Context, jsonCredentials []byte, user string) (*gmailMessageSenderWrapper, error) {
	tokenSource, err := google.JWTAccessTokenSourceFromJSON(jsonCredentials, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JWT credentials: %v", err)
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}
	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}
