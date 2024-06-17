package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/darkrockmountain/gomail"
)

// smtpEmailSender is responsible for sending emails using SMTP.
type smtpEmailSender struct {
	host             string           // The SMTP server host.
	port             int              // The SMTP server port.
	user             string           // The username for authentication.
	password         string           // The password for authentication.
	authMethod       AuthMethod       // The authentication method to use.
	connectionMethod ConnectionMethod // The connection method to use (by default implicit).

}

// NewSmtpEmailSender creates a new instance of smtpEmailSender.
// It initializes the SMTP email sender with the provided SMTP server settings.
//
// Parameters:
//   - host: The SMTP server host (e.g., "smtp.example.com").
//   - port: The SMTP server port (e.g., 587).
//   - user: The username for authentication.
//   - password: The password for authentication.
//   - authMethod: The authentication method to use (AUTH_CRAM_MD5,  AUTH_PLAIN).

// Returns:
//   - *smtpEmailSender: A pointer to the initialized smtpEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSmtpEmailSender(host string, port int, user, password string, authMethod AuthMethod) (*smtpEmailSender, error) {
	// Return the initialized smtpEmailSender with Implicit TLS connection
	return NewSmtpEmailSenderWithConnMethod(host, port, user, password, authMethod, CONN_IMPLICIT)
}

// NewSmtpEmailSenderWithConnMethod creates a new instance of smtpEmailSender.
// It initializes the SMTP email sender with the provided SMTP server settings.
//
// Parameters:
//   - host: The SMTP server host (e.g., "smtp.example.com").
//   - port: The SMTP server port (e.g., 587).
//   - user: The username for authentication.
//   - password: The password for authentication.
//   - authMethod: The authentication method to use (AUTH_CRAM_MD5,  AUTH_PLAIN).
//   - connectionMethod: The connection method to use (CONN_IMPLICIT, CONN_TLS).

// Returns:
//   - *smtpEmailSender: A pointer to the initialized smtpEmailSender.
//   - error: An error if the initialization fails, otherwise nil.
func NewSmtpEmailSenderWithConnMethod(host string, port int, user, password string, authMethod AuthMethod, connectionMethod ConnectionMethod) (*smtpEmailSender, error) {
	// Initialize the smtpEmailSender with the provided settings and return it
	return &smtpEmailSender{
		host,
		port,
		user,
		password,
		authMethod,
		connectionMethod,
	}, nil
}

// AuthMethod defines the authorization method for SMTP.
type AuthMethod string

const (
	AUTH_CRAM_MD5 AuthMethod = "CRAM-MD5" // CRAM-MD5 authentication method.
	AUTH_PLAIN    AuthMethod = "PLAIN"    // Plain authentication method.
)

type ConnectionMethod string

const (
	CONN_IMPLICIT ConnectionMethod = "IMPLICIT" // IMPLICIT TLS authentication method.
	CONN_TLS      ConnectionMethod = "TLS"      // TLS authentication method.

)

// SendEmail sends an email using the specified SMTP settings and authentication method.
func (s *smtpEmailSender) SendEmail(message *gomail.EmailMessage) error {
	// Include CC and BCC recipients in the SMTP envelope
	sendMailTo := message.GetTo()
	sendMailTo = append(sendMailTo, message.GetCC()...)
	sendMailTo = append(sendMailTo, message.GetBCC()...)
	msg, err := gomail.BuildMimeMessage(message)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	if s.authMethod == AUTH_CRAM_MD5 {
		auth = smtp.CRAMMD5Auth(s.user, s.password)
	}

	if s.connectionMethod == CONN_TLS {
		skipInsecure := false
		if os.Getenv("APP_ENV") == "development" {
			skipInsecure = true
		}

		tlsconfig := &tls.Config{
			InsecureSkipVerify: skipInsecure,
			ServerName:         s.host,
		}

		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), tlsconfig)
		if err != nil {
			return err
		}

		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return err
		}

		if err = client.Auth(auth); err != nil {
			return err
		}

		if err = client.Mail(message.GetFrom()); err != nil {
			return err
		}
		for _, addr := range sendMailTo {
			if err = client.Rcpt(addr); err != nil {
				return err
			}
		}

		w, err := client.Data()
		if err != nil {
			return err
		}

		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}

		return client.Quit()
	} else {
		err = smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, message.GetFrom(), sendMailTo, msg)
	}

	return err
}
