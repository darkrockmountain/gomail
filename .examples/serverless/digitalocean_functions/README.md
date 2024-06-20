# DigitalOcean Functions Email Sender Example

This example demonstrates how to use the `EmailSender` interface in a DigitalOcean Function using environment variables for configuration.

## Files

- `do_email_sender.go`: The main Go code for the function.
- `go.mod`: Go modules file for dependency management.
- `Makefile`: Makefile to build and deploy the function.
- `project.yml`: DigitalOcean Project configuration file.
- `Dockerfile`: Dockerfile to containerize the Go application.
- `README.md`: This documentation file.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Node.js and npm are installed.
- Go is installed and configured.
- DigitalOcean CLI (for deployment to DigitalOcean).

## Deployment

1. **Install the DigitalOcean CLI (doctl)**:

    ```sh
    snap install doctl
    ```

2. **Authenticate with DigitalOcean**:

    ```sh
    doctl auth init
    ```

3. **Set Up Environment Variables**:

    Set up the following environment variables in your DigitalOcean App Platform, or create an .env file at the root of this folder:
    _For the example should be:_ 
    - `SMTP_HOST`
    - `SMTP_PORT`
    - `SMTP_USER`
    - `SMTP_PASSWORD`
    - `SMTP_AUTH_METHOD`

4. **Build and Deploy the Function**:

    ```sh
    make deploy
    ```

5. **Test the Function**:

    Invoke the function with a JSON payload matching the `EmailRequest` structure.
