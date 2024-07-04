package common_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/darkrockmountain/gomail/common"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSONCustom(t *testing.T) {
	t.Run("Marshal EmailMessage with attachments", func(t *testing.T) {
		email := common.NewFullEmailMessage(
			"sender@example.com",
			[]string{"recipient@example.com"},
			"Subject",
			[]string{"cc@example.com"},
			[]string{"bcc@example.com"},
			"replyto@example.com",
			"This is the email content.",
			"<p>This is the email content.</p>",
			[]common.Attachment{*common.NewAttachment("attachment1.txt", []byte("file content"))},
		)
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)
		expected := `{"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>","attachments":[{"filename":"attachment1.txt","content":"ZmlsZSBjb250ZW50"}]}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("Marshal EmailMessage without attachments", func(t *testing.T) {
		email := common.NewFullEmailMessage(
			"sender@example.com",
			[]string{"recipient@example.com"},
			"Subject",
			[]string{"cc@example.com"},
			[]string{"bcc@example.com"},
			"replyto@example.com",
			"This is the email content.",
			"<p>This is the email content.</p>",
			nil,
		)
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)
		expected := `{"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>"}`
		assert.JSONEq(t, expected, string(jsonData))
	})
}

func TestUnmarshalJSONCustom(t *testing.T) {
	t.Run("Unmarshal EmailMessage with attachments", func(t *testing.T) {
		jsonData := `{
			"from": "sender@example.com",
			"to": ["recipient@example.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>",
			"attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
		}`
		var email common.EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "sender@example.com", email.GetFrom())
		assert.Equal(t, []string{"recipient@example.com"}, email.GetTo())
		assert.Equal(t, []string{"cc@example.com"}, email.GetCC())
		assert.Equal(t, []string{"bcc@example.com"}, email.GetBCC())
		assert.Equal(t, "replyto@example.com", email.GetReplyTo())
		assert.Equal(t, "Subject", email.GetSubject())
		assert.Equal(t, "This is the email content.", email.GetText())
		assert.Equal(t, "<p>This is the email content.</p>", email.GetHTML())
		expectedAttachment := *common.NewAttachment("attachment1.txt", []byte("file content"))
		assert.Equal(t, []common.Attachment{expectedAttachment}, email.GetAttachments())
	})

	t.Run("Unmarshal EmailMessage without attachments", func(t *testing.T) {
		jsonData := `{
			"from": "sender@example.com",
			"to": ["recipient@example.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>"
		}`
		var email common.EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "sender@example.com", email.GetFrom())
		assert.Equal(t, []string{"recipient@example.com"}, email.GetTo())
		assert.Equal(t, []string{"cc@example.com"}, email.GetCC())
		assert.Equal(t, []string{"bcc@example.com"}, email.GetBCC())
		assert.Equal(t, "replyto@example.com", email.GetReplyTo())
		assert.Equal(t, "Subject", email.GetSubject())
		assert.Equal(t, "This is the email content.", email.GetText())
		assert.Equal(t, "<p>This is the email content.</p>", email.GetHTML())
		assert.Nil(t, email.GetAttachments())
	})
}

func TestMarshalJSONEdgeCases(t *testing.T) {
	t.Run("nil EmailMessage", func(t *testing.T) {
		var email *common.EmailMessage
		result, err := json.Marshal(email)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})

	t.Run("nil Attachment", func(t *testing.T) {
		var attachment *common.Attachment
		result, err := json.Marshal(attachment)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})
}

func TestUnmarshalJSONEdgeCases(t *testing.T) {
	t.Run("empty JSON EmailMessage", func(t *testing.T) {
		jsonData := `{}`
		var email common.EmailMessage
		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "", email.GetFrom())
		assert.Equal(t, "", email.GetSubject())
		assert.Equal(t, "", email.GetText())
		assert.Equal(t, "", email.GetHTML())
	})

	t.Run("invalid JSON EmailMessage", func(t *testing.T) {
		// Example of JSON data that will cause the unmarshaling to fail
		// 'to' should be an array, not a string
		invalidJSONData := `{
        "from": "sender@example.com",
        "to": "invalid_recipient@example.com",
        "cc": ["cc@example.com"],
        "bcc": ["bcc@example.com"],
        "replyTo": "replyto@example.com",
        "subject": "Subject",
        "text": "This is the email content.",
        "html": "<p>This is the email content.</p>",
        "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
    }`

		var email common.EmailMessage
		err := json.Unmarshal([]byte(invalidJSONData), &email)
		assert.Error(t, err)
	})

	t.Run("invalid JSON Attachment", func(t *testing.T) {
		// Example of JSON data that will cause the unmarshaling to fail
		// 'filename' should be an string, not an integer
		jsonData := `{"filename": 123456789, "content": "invalid_base64"}`
		var attachment common.Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.Error(t, err)
	})

	t.Run("empty JSON Attachment", func(t *testing.T) {
		jsonData := `{}`
		var attachment common.Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.Nil(t, err)
		assert.Equal(t, "", attachment.GetFilename())
		assert.Equal(t, []byte{}, attachment.GetRawContent())
	})

	t.Run("invalid base64 content Attachment", func(t *testing.T) {
		jsonData := `{"filename": "file.txt", "content": "invalid_base64"}`
		var attachment common.Attachment
		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.NotNil(t, err)
	})

}

func ExampleEmailMessage_MarshalJSON() {

	email := common.NewFullEmailMessage(
		"sender@example.com",
		[]string{"recipient@example.com"},
		"Subject",
		[]string{"cc@example.com"},
		[]string{"bcc@example.com"},
		"replyto@example.com",
		"This is the email content.",
		"<p>This is the email content.</p>",
		[]common.Attachment{
			*common.NewAttachment("attachment1.txt", []byte("file content")),
		},
	)
	jsonData, err := json.Marshal(email)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println("JSON output:", string(jsonData))
	// Output: JSON output: {"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"\u003cp\u003eThis is the email content.\u003c/p\u003e","attachments":[{"filename":"attachment1.txt","content":"ZmlsZSBjb250ZW50"}]}
}

func ExampleEmailMessage_UnmarshalJSON() {

	jsonInput := `{
	    "from": "sender@example.com",
	    "to": ["recipient@example.com"],
	    "cc": ["cc@example.com"],
	    "bcc": ["bcc@example.com"],
	    "replyTo": "replyto@example.com",
	    "subject": "Subject",
	    "text": "This is the email content.",
	    "html": "<p>This is the email content.</p>",
	    "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
	}`
	var email common.EmailMessage
	err := json.Unmarshal([]byte(jsonInput), &email)
	if err != nil {
		fmt.Println("Error unmarshaling from JSON:", err)
		return
	}

	jsonData, err := json.Marshal(&email)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println("JSON output:", string(jsonData))
	// Output: JSON output: {"from":"sender@example.com","to":["recipient@example.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"\u003cp\u003eThis is the email content.\u003c/p\u003e","attachments":[{"filename":"attachment1.txt","content":"ZmlsZSBjb250ZW50"}]}

}
