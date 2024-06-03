package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/providers"
)

type ResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type Response struct {
	Body       string          `json:"body"`
	StatusCode string          `json:"statusCode"`
	Headers    ResponseHeaders `json:"headers"`
}

func generateResponse(code int, body string) Response {
	return Response{
		Headers:    ResponseHeaders{ContentType: "text/html"},
		StatusCode: fmt.Sprintf("%d", code),
		Body:       body,
	}
}

func Main(ctx context.Context, emailReq gomail.EmailMessage) Response {

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return generateResponse(http.StatusInternalServerError, "Invalid port number")
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	authMethod := providers.AuthMethod(os.Getenv("SMTP_AUTH_METHOD"))

	sender, err := providers.NewSmtpEmailSender(host, port, user, password, authMethod)
	if err != nil {
		return generateResponse(http.StatusInternalServerError, "Failed to initialize email sender")

	}

	if err := sender.SendEmail(emailReq); err != nil {
		return generateResponse(http.StatusInternalServerError, "Failed to send email")
	}

	return generateResponse(http.StatusOK, "Email sent successfully")
}
