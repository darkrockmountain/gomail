package credentials_test

import (
	"testing"

	"github.com/darkrockmountain/gomail/credentials"
	"github.com/stretchr/testify/assert"
)

func TestParseCredentials(t *testing.T) {
	credentialsJSON := []byte(`{
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
	_, err := credentials.ParseCredentials(credentialsJSON)
	assert.NoError(t, err)
}

func TestParseToken(t *testing.T) {
	tokenJSON := []byte(`{
		"access_token": "mock_access_token",
		"token_type": "Bearer",
		"expiry": "2023-10-02T15:00:00Z"
	}`)

	_, err := credentials.ParseToken(tokenJSON)
	assert.NoError(t, err)
}

func TestParseCredentialsInvalidJson(t *testing.T) {
	credentialsJSON := []byte(`{invalid_json}`)

	_, err := credentials.ParseCredentials(credentialsJSON)
	assert.Error(t, err)
}

func TestParseTokenInvalidJson(t *testing.T) {
	tokenJSON := []byte(`{invalid_json}`)

	_, err := credentials.ParseToken(tokenJSON)
	assert.Error(t, err)
}
