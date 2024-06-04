package smtp

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/darkrockmountain/gomail"
	"github.com/stretchr/testify/assert"
)

// TestEmailSenderImplementation checks if smtpEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ gomail.EmailSender = (*smtpEmailSender)(nil)
}

// newMockSMTPServer creates a mock SMTP server for testing purposes.
func newMockSMTPServer(t *testing.T, handler func(conn net.Conn)) *mockSMTPServer {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock SMTP server: %v", err)
	}

	server := &mockSMTPServer{
		listener: listener,
		addr:     listener.Addr().String(),
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handler(conn)
		}
	}()

	return server
}

// mockSMTPServer represents a mock SMTP server.
type mockSMTPServer struct {
	listener net.Listener
	addr     string
}

func (s *mockSMTPServer) Close() {
	s.listener.Close()
}

const expectedPassword = "password" // Replace this with the expected password

// smtpHandler is the handler for plain SMTP connections.
func smtpHandler(conn net.Conn) {
	// Send a welcome message to the client when they connect.
	fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

	// Create a buffer to read data from the connection.
	buf := make([]byte, 1024)

	// Infinite loop to handle multiple commands from the client.
	for {
		// Read data from the connection into the buffer.
		n, err := conn.Read(buf)

		if err != nil {
			// If an error occurs (like EOF), break the loop.
			break
		}

		// Convert the read bytes to a string and trim any whitespace.
		cmd := strings.TrimSpace(string(buf[:n]))
		// Handle different SMTP commands based on the command prefix.
		switch {
		// Handle EHLO command by sending a greeting and supported authentication methods.
		case strings.HasPrefix(cmd, "EHLO"):
			fmt.Fprintln(conn, "250-Hello")
			fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")

		// Handle AUTH PLAIN command by sending an authentication successful message.
		case strings.HasPrefix(cmd, "AUTH PLAIN"):
			creds := decodeConnectionCommand("AUTH PLAIN", cmd)
			// Compare the received password with the expected password.
			if len(creds) == 2 && creds[1] == expectedPassword {
				// If the password matches, send an authentication successful message.
				fmt.Fprintln(conn, "235 Authentication successful")
			} else {
				// If the password does not match, send an authentication failed message.
				fmt.Fprintln(conn, "535 Authentication failed")
			}
			// fmt.Fprintln(conn, "235 Authentication successful")

		// Handle AUTH CRAM-MD5 command by sending a challenge and waiting for a response.
		case strings.HasPrefix(cmd, "AUTH CRAM-MD5"):
			// Send a base64-encoded challenge to the client.
			fmt.Fprintln(conn, "334 "+base64.StdEncoding.EncodeToString([]byte("challenge")))
			// Read the client's response to the challenge.
			conn.Read(buf)

			// Send an authentication successful message.
			fmt.Fprintln(conn, "235 Authentication successful")

		// Handle MAIL FROM command by sending an OK response.
		case strings.HasPrefix(cmd, "MAIL FROM"):
			fmt.Fprintln(conn, "250 OK")

		// Handle RCPT TO command by sending an OK response.
		case strings.HasPrefix(cmd, "RCPT TO"):
			fmt.Fprintln(conn, "250 OK")

		// Handle DATA command by indicating readiness to receive email data.
		case strings.HasPrefix(cmd, "DATA"):
			// Indicate the start of data input.
			fmt.Fprintln(conn, "354 Start mail input; end with <CRLF>.<CRLF>")
			// Read the email data (not used in this mock implementation).
			conn.Read(buf)
			// Send an OK response after receiving the data.
			fmt.Fprintln(conn, "250 OK")

		// Handle QUIT command by sending a goodbye message and closing the connection.

		case strings.HasPrefix(cmd, "QUIT"):
			fmt.Fprintln(conn, "221 Bye")
			conn.Close()
			return
		default:
			fmt.Fprintln(conn, "500 Unrecognized command")

		}
	}
}

// tlsHandler is the handler for TLS connections.
func tlsHandler(handler func(conn net.Conn)) func(conn net.Conn) {

	return func(conn net.Conn) {
		cert := generateKeys()

		tlsConn := tls.Server(conn, &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		})
		err := tlsConn.Handshake()
		if err != nil {
			tlsConn.Close()
			return
		}

		handler(tlsConn)
	}
}

func decodeConnectionCommand(cmd, message string) []string {
	cleanedCmd := strings.TrimSpace(cmd)
	stringB64 := strings.TrimPrefix(message, cleanedCmd+" ")
	decoded, err := base64.StdEncoding.DecodeString(stringB64)
	if err != nil {
		return nil
	}

	// Split the decoded string by the null character
	trimmedDecoded := strings.Trim(string(decoded), "\x00")
	return strings.Split(trimmedDecoded, "\x00")

}
func TestNewSmtpEmailSender(t *testing.T) {
	emailSender, err := NewSmtpEmailSender("smtp.example.com", 587, "user", "password", AUTH_PLAIN)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestSendEmailPlainAuth(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSender(host, portInt, "user", "password", AUTH_PLAIN)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailCramMD5Auth(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSender(host, portInt, "user", "password", AUTH_CRAM_MD5)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailError(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSender(host, portInt, "user", "wrongpassword", AUTH_PLAIN)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

// //////
func TestSendEmailInvalidServer(t *testing.T) {
	emailSender, err := NewSmtpEmailSender("invalid.server.com", 587, "user", "password", AUTH_PLAIN)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailMissingSettings(t *testing.T) {
	emailSender, err := NewSmtpEmailSender("", 0, "", "", AUTH_PLAIN)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailImplicitTLS(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 465
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_IMPLICIT)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailExplicitTLS(t *testing.T) {
	server := newMockSMTPServer(t, tlsHandler(smtpHandler))
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 587
	fmt.Sscanf(port, "%d", &portInt)

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailExplicitErrorTLS(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 587
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage("sender@example.com", []string{"recipient@example.com"}, "Test Email", "This is a test email.")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendTLSEmailConnectionError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				conn.Close()
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendTLSEmailEHLOError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"Test Email",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "500 Unrecognized command", err.Error())

}

func TestSendTLSEmailAUTHError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "403 Forbidden")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "403 Forbidden", err.Error())

}

func TestSendTLSEmailMailError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "550 MAIL ERROR")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "550 MAIL ERROR", err.Error())
}

func TestSendTLSEmailRcptError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "530 No such user")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "530 No such user", err.Error())

}

func TestSendTLSEmailDataError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "DATA"):
						fmt.Fprintln(conn, "554 Service unavailable")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "554 Service unavailable", err.Error())
}

func TestSendTLSEmailDataWriteError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}
					cmd := strings.TrimSpace(string(buf[:n]))
					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "DATA"):
						fmt.Fprintln(conn, "354 Start mail input; end with <CRLF>.<CRLF>")
						conn.Read(buf)
						fmt.Fprintln(conn, "552 Message size exceeds fixed limit")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()
						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25
	fmt.Sscanf(port, "%d", &portInt)

	emailSender, err := NewSmtpEmailSenderWithConnMethod(host, portInt, "user", "password", AUTH_PLAIN, CONN_TLS)
	assert.NoError(t, err)

	message := gomail.NewEmailMessage(
		"sender@example.com",
		[]string{"nonexistent@example.com"},
		"This is a test email.",
		"This is a test email.",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "552 Message size exceeds fixed limit", err.Error())
}

func generateKeys() tls.Certificate {
	cert, err := tls.X509KeyPair([]byte(ecdsaCertPEM), []byte(ecdsaKeyPEM))
	if err != nil {
		log.Panicf("error X509KeyPair %v", err)
	}
	return cert
}

func testingKey(s string) string { return strings.ReplaceAll(s, "TESTING KEY", "PRIVATE KEY") }

//keys from https://go.dev/src/crypto/tls/tls_test.go

var ecdsaCertPEM = `-----BEGIN CERTIFICATE-----
MIIB/jCCAWICCQDscdUxw16XFDAJBgcqhkjOPQQBMEUxCzAJBgNVBAYTAkFVMRMw
EQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0
eSBMdGQwHhcNMTIxMTE0MTI0MDQ4WhcNMTUxMTE0MTI0MDQ4WjBFMQswCQYDVQQG
EwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lk
Z2l0cyBQdHkgTHRkMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBY9+my9OoeSUR
lDQdV/x8LsOuLilthhiS1Tz4aGDHIPwC1mlvnf7fg5lecYpMCrLLhauAc1UJXcgl
01xoLuzgtAEAgv2P/jgytzRSpUYvgLBt1UA0leLYBy6mQQbrNEuqT3INapKIcUv8
XxYP0xMEUksLPq6Ca+CRSqTtrd/23uTnapkwCQYHKoZIzj0EAQOBigAwgYYCQXJo
A7Sl2nLVf+4Iu/tAX/IF4MavARKC4PPHK3zfuGfPR3oCCcsAoz3kAzOeijvd0iXb
H5jBImIxPL4WxQNiBTexAkF8D1EtpYuWdlVQ80/h/f4pBcGiXPqX5h2PQSQY7hP1
+jwM1FGS4fREIOvlBYr/SzzQRtwrvrzGYxDEDbsC0ZGRnA==
-----END CERTIFICATE-----
`

var ecdsaKeyPEM = testingKey(`-----BEGIN EC PARAMETERS-----
BgUrgQQAIw==
-----END EC PARAMETERS-----
-----BEGIN EC TESTING KEY-----
MIHcAgEBBEIBrsoKp0oqcv6/JovJJDoDVSGWdirrkgCWxrprGlzB9o0X8fV675X0
NwuBenXFfeZvVcwluO7/Q9wkYoPd/t3jGImgBwYFK4EEACOhgYkDgYYABAFj36bL
06h5JRGUNB1X/Hwuw64uKW2GGJLVPPhoYMcg/ALWaW+d/t+DmV5xikwKssuFq4Bz
VQldyCXTXGgu7OC0AQCC/Y/+ODK3NFKlRi+AsG3VQDSV4tgHLqZBBus0S6pPcg1q
kohxS/xfFg/TEwRSSws+roJr4JFKpO2t3/be5OdqmQ==
-----END EC TESTING KEY-----
`)
