package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/providers/smtp"
)

const errorOutputMessage = "An unexpected error occurred. Please try again later."

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

func Main(ctx context.Context, jsonData []byte) Response {

	var emailMessage gomail.EmailMessage

	err := json.Unmarshal(jsonData, &emailMessage)
	if err != nil {
		log.Printf("error unmarshaling from JSON: %v", err)
		return generateResponse(http.StatusInternalServerError, errorOutputMessage)

	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("error invalid port number: %v", err)
		return generateResponse(http.StatusInternalServerError, errorOutputMessage)
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	authMethod := smtp.AuthMethod(os.Getenv("SMTP_AUTH_METHOD"))

	sender, err := smtp.NewSmtpEmailSender(host, port, user, password, authMethod)
	if err != nil {
		log.Printf("error failed to initialize email sender %v", err)
		return generateResponse(http.StatusInternalServerError, errorOutputMessage)

	}

	if err := sender.SendEmail(&emailMessage); err != nil {
		log.Printf("error sending email %v", err)
		return generateResponse(http.StatusInternalServerError, "Failed to send email")
	}

	return generateResponse(http.StatusOK, "Email sent successfully")
}
