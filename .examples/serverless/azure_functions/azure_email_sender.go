package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/darkrockmountain/gomail"
	"github.com/darkrockmountain/gomail/providers/smtp"
)

// SendEmail is the HTTP handler that processes the incoming request
// and sends an email using the SMTP settings configured in environment variables.
func SendEmail(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into the EmailMessage struct.
	// Note: The body is expected to contain email message in JSON format, but it can be adapted to handle any other structure as needed.
	var emailReq gomail.EmailMessage
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the configuration values from the environment or from a secret manager.
	// These values are chosen for the SMTP setup but should be adapted to the ones you need.
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		http.Error(w, "Invalid port number", http.StatusInternalServerError)
		return
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	authMethod := smtp.AuthMethod(os.Getenv("SMTP_AUTH_METHOD"))

	// Initialize the SMTP email sender with the retrieved configuration values.
	// You can choose any other email sender implementation by replacing NewSmtpEmailSender with another constructor.
	sender, err := smtp.NewSmtpEmailSender(
		host,
		port,
		user,
		password,
		authMethod,
	)
	if err != nil {
		http.Error(w, "Failed to initialize email sender", http.StatusInternalServerError)
		return
	}

	// Send the email using the initialized sender.
	if err := sender.SendEmail(&emailReq); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Return a successful response.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

// getPort determines the port number on which the HTTP server should listen.
// It first checks if the `FUNCTIONS_CUSTOMHANDLER_PORT` environment variable is set,
// which is specific to Azure Functions. If the environment variable is found,
// the function uses its value as the port number. If the environment variable is not set,
// it defaults to port 8080.
//
// This function is essential for running custom handlers in Azure Functions, as Azure
// dynamically assigns a port for the function to use. The assigned port number is
// provided via the `FUNCTIONS_CUSTOMHANDLER_PORT` environment variable.
//
// Returns:
//   - A string representing the port number in the format ":<port-number>".
//
// Example:
//
//	port := getPort()
//	http.ListenAndServe(port, nil)
//
// In the example above, the HTTP server will listen on the port specified by the
// `FUNCTIONS_CUSTOMHANDLER_PORT` environment variable, or port 8080 if the environment
// variable is not set.
func getPort() string {
	// Default port is set to 8080
	port := ":8080"

	// Check if the environment variable FUNCTIONS_CUSTOMHANDLER_PORT is set
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		// If set, use the value of the environment variable as the port number
		port = ":" + val
	}

	// Return the port number in the format ":<port-number>"
	return port
}

func main() {
	// The pattern for Azure Functions HTTP triggers has the form /api/<function-route>
	// In our example, the route is specified in SendEmail/function.json under bindings->route
	http.HandleFunc("/api/send_email", SendEmail)

	// Start the HTTP server on the port specified by the FUNCTIONS_CUSTOMHANDLER_PORT
	// environment variable. This is required for Azure Functions custom handlers.
	http.ListenAndServe(getPort(), nil)
}
