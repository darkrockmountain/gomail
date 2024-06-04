# Obtaining AWS SES Credentials

To send emails using the AWS SES API, you need to set up AWS credentials and obtain access keys.

## Steps

### 1. Sign up for AWS:
- Go to [AWS](https://aws.amazon.com/).
- Sign up for an account if you don't already have one.

### 2. Create an IAM User with SES Permissions:
- Navigate to the IAM service in the AWS Management Console.
- Click on "Users" and then "Add user".
- Enter a user name and select "Programmatic access" under "Access type".
- Click "Next: Permissions" and attach the "AmazonSESFullAccess" policy (or create a custom policy with necessary permissions).
- Click "Next: Tags", "Next: Review", and then "Create user".
- Download the generated AWS access key ID and secret access key.

### 3. Update the SES settings in the `ses_email_sender.go` file:
```go
accessKeyID := "your-aws-access-key-id"
secretAccessKey := "your-aws-secret-access-key"
region := "your-aws-region"
sender := "your-email@example.com"
```
### 4. Run the example function to test the email sending:
```go
// Replace these with your actual values (in production, retrieve from a secure file or secret manager)
accessKeyID := "your-aws-access-key-id"
secretAccessKey := "your-aws-secret-access-key"
region := "your-aws-region"
sender := "your-email@example.com"

// Create SESEmailSender with the credentials
emailSender, err := NewSESEmailSender(region, sender, accessKeyID, secretAccessKey)
if err != nil {
    log.Fatalf("Failed to create email sender: %v", err)
}

// Read attachment content
attachmentContent, err := os.ReadFile("path/to/attachment.jpg")
if err != nil {
    log.Fatalf("Failed to read attachment: %v", err)
}

// Define email message
message := *gomail.NewEmailMessage(sender,[]string{"recipient@example.com"}, "Test Email with attachment", "This is the plain text part of the email.").
		SetHTML("<p>This is the <b>HTML</b> part of the <i>email</i>.</p>").
		AddAttachments(gomail.Attachment{
			Filename: "attachment.jpg",  Content: attachmentContent,
		})



// Send email
if err := emailSender.SendEmail(message); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}

fmt.Println("Email sent successfully")
```