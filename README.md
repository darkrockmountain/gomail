# GoMail - Email Sender Library


![Build Status](https://github.com/darkrockmountain/gomail/actions/workflows/ci.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/darkrockmountain/gomail?status.svg)](https://pkg.go.dev/github.com/darkrockmountain/gomail)
[![Go Report Card](https://goreportcard.com/badge/github.com/darkrockmountain/gomail?branch=master&kill_cache=1)](https://goreportcard.com/report/github.com/darkrockmountain/gomail)
[![codecov](https://codecov.io/gh/darkrockmountain/gomail/graph/badge.svg?token=NC0O7RMK2X)](https://codecov.io/gh/darkrockmountain/gomail)
![Vulnerability assessment Status](https://github.com/darkrockmountain/gomail/actions/workflows/govulncheck.yaml/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdarkrockmountain%2Fgomail.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdarkrockmountain%2Fgomail?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdarkrockmountain%2Fgomail.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdarkrockmountain%2Fgomail?ref=badge_shield&issueType=security)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/darkrockmountain/gomail/badge)](https://scorecard.dev/viewer/?uri=github.com/darkrockmountain/gomail)
[![GitHub Release](https://img.shields.io/github/v/release/darkrockmountain/gomail)](https://github.com/darkrockmountain/gomail/releases)
<!-- [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=your-project-key&metric=alert_status)](https://sonarcloud.io/dashboard?id=your-project-key) -->
## Project Description


This project provides implementations for sending emails using different services including SMTP, Gmail API, Microsoft Graph API, SendGrid, AWS SES, Mailgun, Mandrill, Postmark, and SparkPost. Each implementation follows a common interface, allowing for flexibility and easy integration with various email services.

## Features

- Send emails using various providers: SMTP, Gmail API, Microsoft Graph API, SendGrid, AWS SES, Mailgun, Mandrill, Postmark, SparkPost.
- Support for attachments and both plain text and HTML content.
- Easy configuration and setup.

## Getting Started

### Prerequisites

- Go 1.22+
- Access to the relevant email service provider (SMTP server, Gmail, Microsoft 365, SendGrid, AWS SES, Mailgun, Mandrill, Postmark, SparkPost)

### Download
```bash
go get github.com/darkrockmountain/gomail 
```
- **Alternatively** download repository:
    ```bash
    git clone https://github.com/darkrockmountain/gomail.git
    cd gomail
    ```
#### Install dependencies  

```bash
go mod tidy
```

### Usage

#### 1. SMTP Email Sender
- Configure your SMTP server settings in `providers/smpt/smtp_email_sender.go`.
- Refer to the [SMTP Credentials Documentation](./docs/SMTP_Credentials.md) for details on obtaining credentials.
- Run the `smtpExample()` function to send a test email.

#### 2. Gmail Email Sender
- Configure your Gmail API credentials in `providers/gmail/gmail_email_sender.go`.
- Refer to the [Gmail Credentials Documentation](./docs/Gmail_Credentials_API_Key.md) for details on obtaining credentials.
- Run the `gExample()` function to send a test email.

#### 3. Gmail Email Sender using OAuth2
- Configure your Gmail API credentials and token in `providers/gmail/gmail_email_sender_oauth2.go`.
- Refer to the [Gmail OAuth2 Credentials Documentation](./docs/Gmail_Credentials_OAuth2.md) for details on obtaining credentials.
- Run the `gExampleOauth2()` function to send a test email.

#### 4. Microsoft 365 Email Sender
- Configure your Microsoft Graph API credentials in `providers/microsoft365/microsoft365_email_sender.go`.
- Refer to the [Microsoft 365 Credentials Documentation](./docs/Microsoft365_Credentials_ROPC.md) for details on obtaining credentials.
- Run the `msGraphExample()` function to send a test email.

#### 5. SendGrid Email Sender
- Configure your SendGrid API key in `providers/sendgrid/sendgrid_email_sender.go`.
- Refer to the [SendGrid Credentials Documentation](./docs/SendGrid_Credentials.md) for details on obtaining credentials.
- Run the `sendgridExample()` function to send a test email.

#### 6. AWS SES Email Sender
- Configure your AWS SES credentials in `providers/ses/ses_email_sender.go`.
- Refer to the [AWS SES Credentials Documentation](./docs/AWS_SES_Credentials.md) for details on obtaining credentials.
- Run the `sesExample()` function to send a test email.

#### 7. Mailgun Email Sender
- Configure your Mailgun API key in `providers/mailgun/mailgun_email_sender.go`.
- Refer to the [Mailgun Credentials Documentation](./docs/Mailgun_Credentials.md) for details on obtaining credentials.
- Run the `mailgunExample()` function to send a test email.

#### 8. Mandrill Email Sender
- Configure your Mandrill API key in `providers/mandrill/mandrill_email_sender.go`.
- Refer to the [Mandrill Credentials Documentation](./docs/Mandrill_Credentials.md) for details on obtaining credentials.
- Run the `mandrillExample()` function to send a test email.

#### 9. Postmark Email Sender
- Configure your Postmark API key in `providers/postmark/postmark_email_sender.go`.
- Refer to the [Postmark Credentials Documentation](./docs/Postmark_Credentials.md) for details on obtaining credentials.
- Run the `postmarkExample()` function to send a test email.

#### 10. SparkPost Email Sender
- Configure your SparkPost API key in `providers/sparkpost/sparkpost_email_sender.go`.
- Refer to the [SparkPost Documentation](https://developers.sparkpost.com/api/) for details on obtaining credentials.
- Run the `sparkpostExample()` function to send a test email.

## Documentation

For detailed instructions on obtaining the necessary credentials for each implementation, refer to the respective documentation files in the `docs` directory.

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.


## Error Handling and Troubleshooting

### Common Errors

- **Authentication Failed**: Ensure your API keys and credentials are correct and have the necessary permissions.
- **Network Issues**: Verify your network connectivity and ensure your server can reach the email service provider.
- **Invalid Email Addresses**: Check that all email addresses are correctly formatted.
- **Attachment Issues**: Ensure attachments are correctly encoded and within size limits.

### Troubleshooting Tips

1. **Check Logs**: Always check your application logs for detailed error messages.
2. **Validate Credentials**: Double-check your credentials and permissions.
3. **API Limits**: Ensure you are not exceeding API rate limits or quotas.
4. **Service Status**: Verify that the email service provider is operational and not experiencing downtime.

## Security Best Practices

1. **Environment Variables**: Use environment variables to store credentials.
2. **Secret Managers**: Use secret management services like AWS Secrets Manager, Google Secret Manager, or HashiCorp Vault.
3. **Encryption**: Encrypt sensitive information both at rest and in transit.
4. **Least Privilege**: Follow the principle of least privilege for API keys and credentials.
