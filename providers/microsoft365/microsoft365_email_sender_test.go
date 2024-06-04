package microsoft365

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/darkrockmountain/gomail"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	khttp "github.com/microsoft/kiota-http-go"
	ksForm "github.com/microsoft/kiota-serialization-form-go"
	ksJson "github.com/microsoft/kiota-serialization-json-go"
	ksMP "github.com/microsoft/kiota-serialization-multipart-go"
	ksText "github.com/microsoft/kiota-serialization-text-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/stretchr/testify/assert"
)

type MockMSAuthenticationProvider struct {
}

func (m *MockMSAuthenticationProvider) AuthenticateRequest(context context.Context, request *abstractions.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	return nil
}

// MSMockRoundTripper implements the http.RoundTripper interface
type MSMockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *MSMockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
	}, nil
}

func mockItemSendMailRequestBuild(err error) *users.UserItemRequestBuilder {

	defaultClient := &http.Client{
		Transport: &MSMockRoundTripper{
			Err: err,
		},
	}

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(&MockMSAuthenticationProvider{}, nil, nil, defaultClient)
	if err != nil {
		panic(err)
	}

	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory { return ksJson.NewJsonSerializationWriterFactory() })
	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory { return ksText.NewTextSerializationWriterFactory() })
	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory { return ksForm.NewFormSerializationWriterFactory() })
	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory { return ksMP.NewMultipartSerializationWriterFactory() })
	abstractions.RegisterDefaultDeserializer(func() serialization.ParseNodeFactory { return ksJson.NewJsonParseNodeFactory() })
	abstractions.RegisterDefaultDeserializer(func() serialization.ParseNodeFactory { return ksText.NewTextParseNodeFactory() })
	abstractions.RegisterDefaultDeserializer(func() serialization.ParseNodeFactory { return ksForm.NewFormParseNodeFactory() })

	sendMail := users.UserItemRequestBuilder{
		BaseRequestBuilder: abstractions.BaseRequestBuilder{
			RequestAdapter: adapter,
			UrlTemplate:    "{+baseurl}/users{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24top}",
			PathParameters: map[string]string{"user%2Did": "testUserID"},
		},
	}
	return &sendMail
}

func TestSendEmail(t *testing.T) {
	attachment := gomail.Attachment{
		Filename: "test.txt",
		Content:  []byte("test_attachment_content"),
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment)

	reqBuild := mockItemSendMailRequestBuild(nil)
	sender := &mSGraphEmailSender{
		userRequestBuilder: reqBuild,
	}

	err := sender.SendEmail(message)
	assert.NoError(t, err)

}

func TestComposeMessage_PlainTextEmail(t *testing.T) {
	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	)

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "Test Subject", *msMessage.GetSubject())
	assert.Equal(t, "Test Body", *msMessage.GetBody().GetContent())
	assert.Equal(t, models.TEXT_BODYTYPE, *msMessage.GetBody().GetContentType())
	assert.Equal(t, "sender@example.com", *msMessage.GetFrom().GetEmailAddress().GetAddress())
	assert.Equal(t, "recipient@example.com", *msMessage.GetToRecipients()[0].GetEmailAddress().GetAddress())
}

func TestComposeMessage_HTMLEmail(t *testing.T) {
	message := *gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Subject", "<p>Test Body</p>")

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "Test Subject", *msMessage.GetSubject())
	assert.Equal(t, "<p>Test Body</p>", *msMessage.GetBody().GetContent())
	assert.Equal(t, models.HTML_BODYTYPE, *msMessage.GetBody().GetContentType())
}

func TestComposeMessage_WithCC(t *testing.T) {
	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).SetCC([]string{"cc@example.com"})

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "cc@example.com", *msMessage.GetCcRecipients()[0].GetEmailAddress().GetAddress())
}

func TestComposeMessage_WithBCC(t *testing.T) {
	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).SetBCC([]string{"bcc@example.com"})

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "bcc@example.com", *msMessage.GetBccRecipients()[0].GetEmailAddress().GetAddress())
}

func TestComposeMessage_WithReplyTo(t *testing.T) {
	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).SetReplyTo("replyto@example.com")

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "replyto@example.com", *msMessage.GetReplyTo()[0].GetEmailAddress().GetAddress())
}

func TestComposeMessage_WithAttachments(t *testing.T) {

	attachment := gomail.Attachment{
		Filename: "test.txt",
		Content:  []byte("test_attachment_content"),
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment)

	msMessage := composeMsMessage(message)

	attachmentBase64 := []byte(base64.StdEncoding.EncodeToString([]byte("test_attachment_content")))

	assert.NotNil(t, msMessage)
	assert.Equal(t, "test.txt", *msMessage.GetAttachments()[0].(*models.FileAttachment).GetName())
	assert.Equal(t, attachmentBase64, msMessage.GetAttachments()[0].(*models.FileAttachment).GetContentBytes())
}

func TestComposeMessage_EmptyFields(t *testing.T) {
	message := *gomail.NewEmailMessage("", []string{}, "", "")
	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Empty(t, msMessage.GetSubject())
	assert.Empty(t, msMessage.GetBody().GetContent())
	assert.Empty(t, msMessage.GetFrom())
	assert.Nil(t, msMessage.GetToRecipients())
	assert.Nil(t, msMessage.GetCcRecipients())
	assert.Nil(t, msMessage.GetBccRecipients())
	assert.Nil(t, msMessage.GetReplyTo())
	assert.Nil(t, msMessage.GetAttachments())
}

func TestNewMSGraphEmailSenderAppAuth(t *testing.T) {
	clientId := "test-client-id"
	tenantId := "test-tenant-id"
	clientSecret := "test-client-secret"
	user := "test-user"

	sender, err := NewMSGraphEmailSenderAppAuth(clientId, tenantId, clientSecret, user)
	assert.NoError(t, err)
	assert.NotNil(t, sender)
	assert.NotNil(t, sender.userRequestBuilder)
}

func TestNewMSGraphEmailSenderUserAuth(t *testing.T) {
	clientId := "test-client-id"
	tenantId := "test-tenant-id"

	sender, err := NewMSGraphEmailSenderUserAuth(clientId, tenantId)
	assert.NoError(t, err)
	assert.NotNil(t, sender)
	assert.NotNil(t, sender.userRequestBuilder)
}

func TestNewMSGraphEmailSenderManagedIdentity(t *testing.T) {
	user := "test-user"

	sender, err := NewMSGraphEmailSenderManagedIdentity(user)
	assert.NoError(t, err)
	assert.NotNil(t, sender)
	assert.NotNil(t, sender.userRequestBuilder)
}

func TestNewMSGraphEmailSenderROPC(t *testing.T) {
	clientId := "test-client-id"
	tenantId := "test-tenant-id"
	username := "test-username"
	password := "test-password"

	sender, err := NewMSGraphEmailSenderROPC(clientId, tenantId, username, password)
	assert.NoError(t, err)
	assert.NotNil(t, sender)
	assert.NotNil(t, sender.userRequestBuilder)
}

///////////////

func TestSendEmail_NoUserRequestBuilder(t *testing.T) {
	attachment := gomail.Attachment{
		Filename: "test.txt",
		Content:  []byte("test_attachment_content"),
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment)

	sender := &mSGraphEmailSender{
		userRequestBuilder: nil,
	}

	err := sender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "error: no user request builder available", err.Error())
}

func TestComposeMessage_LargeAttachment(t *testing.T) {
	largeContent := make([]byte, 10*1024*1024) // 10 MB
	attachment := gomail.Attachment{
		Filename: "large.txt",
		Content:  largeContent,
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment)

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "large.txt", *msMessage.GetAttachments()[0].(*models.FileAttachment).GetName())
	assert.True(t, bytes.Equal(attachment.GetBase64Content(), msMessage.GetAttachments()[0].(*models.FileAttachment).GetContentBytes()), "Attachment do not match")
}

func TestComposeMessage_MultipleAttachments(t *testing.T) {
	attachment1 := gomail.Attachment{
		Filename: "test1.txt",
		Content:  []byte("test_attachment_content1"),
	}
	attachment2 := gomail.Attachment{
		Filename: "test2.txt",
		Content:  []byte("test_attachment_content2"),
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment1).AddAttachment(attachment2)

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Equal(t, "test1.txt", *msMessage.GetAttachments()[0].(*models.FileAttachment).GetName())
	assert.Equal(t, attachment1.GetBase64Content(), msMessage.GetAttachments()[0].(*models.FileAttachment).GetContentBytes())
	assert.Equal(t, "test2.txt", *msMessage.GetAttachments()[1].(*models.FileAttachment).GetName())
	assert.Equal(t, attachment2.GetBase64Content(), msMessage.GetAttachments()[1].(*models.FileAttachment).GetContentBytes())
}

func TestComposeMessage_MissingFields(t *testing.T) {
	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"", "")

	msMessage := composeMsMessage(message)

	assert.NotNil(t, msMessage)
	assert.Empty(t, msMessage.GetSubject())
	assert.Empty(t, msMessage.GetBody().GetContent())
	assert.Equal(t, "sender@example.com", *msMessage.GetFrom().GetEmailAddress().GetAddress())
	assert.Equal(t, "recipient@example.com", *msMessage.GetToRecipients()[0].GetEmailAddress().GetAddress())
	assert.Nil(t, msMessage.GetCcRecipients())
	assert.Nil(t, msMessage.GetBccRecipients())
	assert.Nil(t, msMessage.GetReplyTo())
	assert.Nil(t, msMessage.GetAttachments())
}

func TestSendEmail_FailedSend(t *testing.T) {
	attachment := gomail.Attachment{
		Filename: "test.txt",
		Content:  []byte("test_attachment_content"),
	}

	message := *gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	).AddAttachment(attachment)

	clientError := fmt.Errorf("mock error")
	reqBuild := mockItemSendMailRequestBuild(clientError)
	sender := &mSGraphEmailSender{
		userRequestBuilder: reqBuild,
	}

	err := sender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf(`unable to send message: Post "/users/testUserID/sendMail": %v`, clientError), err)
}
