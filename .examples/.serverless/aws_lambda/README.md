# AWS Lambda Email Sender Example

This example demonstrates how to use the `EmailSender` interface in an AWS Lambda function.

## Files

- `lambda_email_sender.go`: The main Go code for the Lambda function.
- `go.mod`: Go modules file for dependency management.
- `Makefile`: Makefile to build, package, and deploy the Lambda function.
- `README.md`: This documentation file.

## Deployment

1. **Build and Package the Function**:

    ```sh
    make package
    ```

2. **Deploy the Function**:

    ```sh
    make deploy
    ```

3. **Set Up Environment Variables**:

    Set up the following environment variables in the AWS Lambda console:
     _For the example should be:_
    - `SMTP_HOST`
    - `SMTP_PORT`
    - `SMTP_USER`
    - `SMTP_PASSWORD`
    - `SMTP_AUTH_METHOD`

4. **Test the Function**:

    Invoke the Lambda function with a JSON payload matching the `EmailRequest` structure.
