package credentials

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// ParseCredentials parses the OAuth2 credentials JSON byte slice and returns an *oauth2.Config.
//
// Parameters:
//   - credentialsJSON: A byte slice containing the OAuth2 client credentials JSON.
//
// Returns:
//   - *oauth2.Config: The parsed OAuth2 config.
//   - error: An error if parsing the credentials fails.
func ParseCredentials(credentialsJSON []byte) (*oauth2.Config, error) {
	config, err := google.ConfigFromJSON(credentialsJSON, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	return config, nil
}

// ParseToken parses the OAuth2 token JSON byte slice and returns an *oauth2.Token.
//
// Parameters:
//   - tokenJSON: A byte slice containing the OAuth2 token JSON.
//
// Returns:
//   - *oauth2.Token: The parsed OAuth2 token.
//   - error: An error if parsing the token fails.
func ParseToken(tokenJSON []byte) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	if err := json.Unmarshal(tokenJSON, token); err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}
	return token, nil
}
