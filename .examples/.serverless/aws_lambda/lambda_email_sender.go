package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/providers"
)

// handler is the main Lambda function handler that processes the incoming API Gateway request
// and sends an email using the SMTP settings configured in environment variables.
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the JSON request body into the EmailMessage struct.
	// Note: The body is expected to contain email message in JSON format, but it can be adapted to handle any other structure as needed.
	var emailReq gomail.EmailMessage
	if err := json.Unmarshal([]byte(request.Body), &emailReq); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request payload",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Retrieve the configuration values from the environment or from a secret manager.
	// These values are chosen for the SMTP setup but should be adapted to the ones you need.
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	authMethod := providers.AuthMethod(os.Getenv("SMTP_AUTH_METHOD"))

	// Initialize the SMTP email sender with the retrieved configuration values.
	// You can choose any other email sender implementation by replacing NewSmtpEmailSender with another constructor.
	sender, err := providers.NewSmtpEmailSender(
		host,
		port,
		user,
		password,
		authMethod,
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to initialize email sender",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Send the email using the initialized sender.
	if err := sender.SendEmail(emailReq); err != nil {
		fmt.Printf("Error sending email %v, with error %v\n", emailReq, err)
		return events.APIGatewayProxyResponse{
			Body:       "Failed to send email",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Return a successful response.
	return events.APIGatewayProxyResponse{
		Body:       "Email sent successfully",
		StatusCode: 200,
	}, nil
}

// main function starts the Lambda function handler.
func main() {
	lambda.Start(handler)
}
