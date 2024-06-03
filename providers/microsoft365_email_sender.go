package providers

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/darkrockmountain/gomail"
	azureAuth "github.com/microsoft/kiota-authentication-azure-go"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

// mSGraphEmailSender defines an abstract class for sending emails using the Microsoft Graph API.
// This class contains the SendEmail method and a *graphusers.UserItemRequestBuilder for interacting with the user's endpoint.
type mSGraphEmailSender struct {
	userRequestBuilder *graphusers.UserItemRequestBuilder
}

// SendEmail sends an email using the Microsoft Graph API.
//
// Parameters:
//   - message: An EmailMessage struct that contains the details of the email to be sent, including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - error: An error if the email sending fails, otherwise nil.
func (s *mSGraphEmailSender) SendEmail(message gomail.EmailMessage) error {
	if s.userRequestBuilder == nil {
		return fmt.Errorf("error: no user request builder available")
	}

	msMessage := composeMsMessage(message)
	// Send the message
	sendMailPostRequestBody := graphusers.NewItemSendmailSendMailPostRequestBody()
	sendMailPostRequestBody.SetMessage(msMessage)

	err := s.userRequestBuilder.SendMail().Post(context.Background(), sendMailPostRequestBody, nil)
	if err != nil {
		return fmt.Errorf("unable to send message: %v", err)
	}

	return nil
}

// composeMsMessage creates a new email message based on the provided message.
//
// Parameters:
//   - message: An EmailMessage struct that contains the details of the email to be sent,
//     including the sender, recipients, subject, body content, and any attachments.
//
// Returns:
//   - *models.Message: A pointer to the composed Message object ready to be sent via the Microsoft Graph API.
func composeMsMessage(message gomail.EmailMessage) *models.Message {

	// Create the message
	body := models.NewItemBody()

	if message.GetHTML() != "" {
		htmlBodyType := models.HTML_BODYTYPE
		body.SetContentType(&htmlBodyType)
		body.SetContent(strPtr(message.GetHTML()))
	} else {
		textBodyType := models.TEXT_BODYTYPE
		body.SetContentType(&textBodyType)
		body.SetContent(strPtr(message.GetText()))
	}

	msMessage := models.NewMessage()
	msMessage.SetSubject(strPtr(message.GetSubject()))
	msMessage.SetBody(body)

	// Add sender
	from := message.GetFrom()
	if from != "" {
		senderEmailAddress := models.NewEmailAddress()
		senderEmailAddress.SetAddress(&from)
		senderObj := models.NewRecipient()
		senderObj.SetEmailAddress(senderEmailAddress)
		msMessage.SetFrom(senderObj)
	}

	// Add recipients
	optionRecipients := message.GetTo()
	if len(optionRecipients) > 0 {
		toRecipients := make([]models.Recipientable, len(optionRecipients))
		for i, recipient := range optionRecipients {
			emailAddress := models.NewEmailAddress()
			emailAddress.SetAddress(&recipient)
			recipientObj := models.NewRecipient()
			recipientObj.SetEmailAddress(emailAddress)
			toRecipients[i] = recipientObj
		}
		msMessage.SetToRecipients(toRecipients)
	}

	// Set CC recipients if any
	optionCC := message.GetCC()
	if len(optionCC) > 0 {
		bccRecipients := make([]models.Recipientable, len(optionCC))
		for i, recipientCC := range optionCC {
			emailAddress := models.NewEmailAddress()
			emailAddress.SetAddress(&recipientCC)
			recipientObj := models.NewRecipient()
			recipientObj.SetEmailAddress(emailAddress)
			bccRecipients[i] = recipientObj
		}
		msMessage.SetCcRecipients(bccRecipients)
	}

	// Set BCC recipients if any
	optionBCC := message.GetBCC()
	if len(optionBCC) > 0 {
		bccRecipients := make([]models.Recipientable, len(optionBCC))
		for i, recipientBCC := range optionBCC {
			emailAddress := models.NewEmailAddress()
			emailAddress.SetAddress(&recipientBCC)
			recipientObj := models.NewRecipient()
			recipientObj.SetEmailAddress(emailAddress)
			bccRecipients[i] = recipientObj
		}
		msMessage.SetBccRecipients(bccRecipients)
	}

	// Set Reply-To address if provided
	replyTo := message.GetReplyTo()
	if replyTo != "" {
		replyToEmailAddress := models.NewEmailAddress()
		replyToEmailAddress.SetAddress(&replyTo)
		replyToRecipient := models.NewRecipient()
		replyToRecipient.SetEmailAddress(replyToEmailAddress)
		msMessage.SetReplyTo([]models.Recipientable{replyToRecipient})
	}

	// Add attachments if any
	attachments := message.GetAttachments()
	if len(attachments) > 0 {
		msAttachments := make([]models.Attachmentable, len(attachments))
		for i, attachment := range attachments {
			fileAttachment := models.NewFileAttachment()
			fileAttachment.SetName(strPtr(attachment.GetFilename()))
			fileAttachment.SetContentBytes([]byte(attachment.GetBase64StringContent()))
			msAttachments[i] = fileAttachment
		}
		msMessage.SetAttachments(msAttachments)
	}

	return msMessage

}

// createGraphClient initializes a Microsoft Graph service client using the provided authentication provider.
//
// Parameters:
//   - authProvider: A pointer to the authentication provider for Microsoft Graph.
//
// Returns:
//   - *msgraph.GraphServiceClient: A pointer to the initialized Microsoft Graph service client.
//   - error: An error if the initialization fails, otherwise nil.
func createGraphClient(authProvider *azureAuth.AzureIdentityAuthenticationProvider) (*msgraph.GraphServiceClient, error) {
	adapter, err := msgraph.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, fmt.Errorf("unable to create graph adapter: %v", err)
	}

	graphClient := msgraph.NewGraphServiceClient(adapter)
	return graphClient, nil
}

// NewMSGraphEmailSenderAppAuth creates a new instance of mSGraphEmailSender with application authentication.
// It initializes the Microsoft Graph email sender with the provided Azure AD client credentials.
//
// Parameters:
//   - clientId: The client ID of the Azure AD application.
//   - tenantId: The tenant ID of the Azure AD directory.
//   - clientSecret: The client secret of the Azure AD application.
//   - user: The user account from where the emails will be sent. Can be a real user or a shared mailbox.
//
// Returns:
//   - *mSGraphEmailSender: A pointer to the initialized mSGraphEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMSGraphEmailSenderAppAuth(clientId, tenantId, clientSecret, user string) (*mSGraphEmailSender, error) {
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create credential: %v", err)
	}

	authProvider, err := azureAuth.NewAzureIdentityAuthenticationProvider(cred)
	if err != nil {
		return nil, fmt.Errorf("unable to create auth provider: %v", err)
	}

	graphClient, err := createGraphClient(authProvider)
	if err != nil {
		return nil, err
	}

	emailSender := &mSGraphEmailSender{
		userRequestBuilder: graphClient.Users().ByUserId(user),
	}

	return emailSender, nil
}

// NewMSGraphEmailSenderUserAuth creates a new instance of mSGraphEmailSender with user (delegated) authentication.
// It initializes the Microsoft Graph email sender with the provided Azure AD client credentials using delegated authentication.
//
// Parameters:
//   - clientId: The client ID of the Azure AD application.
//   - tenantId: The tenant ID of the Azure AD directory.
//
// Returns:
//   - *mSGraphEmailSender: A pointer to the initialized mSGraphEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMSGraphEmailSenderUserAuth(clientId, tenantId string) (*mSGraphEmailSender, error) {
	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		TenantID: tenantId,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create device code credential: %v", err)
	}

	authProvider, err := azureAuth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{"mail.send"})
	if err != nil {
		return nil, fmt.Errorf("unable to create auth provider: %v", err)
	}

	graphClient, err := createGraphClient(authProvider)
	if err != nil {
		return nil, err
	}

	emailSender := &mSGraphEmailSender{
		userRequestBuilder: graphClient.Me(),
	}

	return emailSender, nil
}

// NewMSGraphEmailSenderManagedIdentity creates a new instance of mSGraphEmailSender using Managed Identity authentication.
// It initializes the Microsoft Graph email sender with the provided Azure AD managed identity credentials.
//
// Parameters:
//   - user: The user account from where the emails will be sent. Can be a real user or a shared mailbox.
//
// Returns:
//   - *mSGraphEmailSender: A pointer to the initialized mSGraphEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMSGraphEmailSenderManagedIdentity(user string) (*mSGraphEmailSender, error) {
	cred, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create managed identity credential: %v", err)
	}

	authProvider, err := azureAuth.NewAzureIdentityAuthenticationProvider(cred)
	if err != nil {
		return nil, fmt.Errorf("unable to create auth provider: %v", err)
	}

	graphClient, err := createGraphClient(authProvider)
	if err != nil {
		return nil, err
	}

	emailSender := &mSGraphEmailSender{
		userRequestBuilder: graphClient.Users().ByUserId(user),
	}

	return emailSender, nil
}

// NewMSGraphEmailSenderROPC creates a new instance of mSGraphEmailSender using Resource Owner Password Credentials (ROPC) authentication.
// It initializes the Microsoft Graph email sender with the provided Azure AD client credentials and user credentials.
//
// Parameters:
//   - clientId: The client ID of the Azure AD application.
//   - tenantId: The tenant ID of the Azure AD directory.
//   - username: The username of the user.
//   - password: The password of the user.
//
// Returns:
//   - *mSGraphEmailSender: A pointer to the initialized mSGraphEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewMSGraphEmailSenderROPC(clientId, tenantId, username, password string) (*mSGraphEmailSender, error) {
	cred, err := azidentity.NewUsernamePasswordCredential(tenantId, clientId, username, password, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create username/password credential: %v", err)
	}

	authProvider, err := azureAuth.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, fmt.Errorf("unable to create auth provider: %v", err)
	}

	graphClient, err := createGraphClient(authProvider)
	if err != nil {
		return nil, err
	}

	emailSender := &mSGraphEmailSender{
		userRequestBuilder: graphClient.Me(),
	}

	return emailSender, nil
}
