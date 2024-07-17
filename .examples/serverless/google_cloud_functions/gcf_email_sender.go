package function

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/providers/smtp"
)

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var emailReq gomail.EmailMessage
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var sender gomail.EmailSender
	// Initialize sender with an EmailSender implementation located in providers package (or implement your own one)
	// For example:
	// you can retrieve the configuration values form the environment or from a secret manager
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		http.Error(w, "Invalid port number", http.StatusInternalServerError)
		return
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	authMethod := smtp.AuthMethod(os.Getenv("SMTP_AUTH_METHOD"))

	sender, err = smtp.NewSmtpEmailSender(
		host,
		port,
		user,
		password,
		authMethod,
	)

	if err := sender.SendEmail(emailReq); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func init() {
	functions.HTTP("SendEmail", SendEmail)
}
