# Azure Functions Email Sender Example

This example demonstrates how to use the `EmailSender` interface in an Azure Function using environment variables for configuration.

## Files

- `azure_email_sender.go`: The main Go code for the Azure Function.
- `go.mod`: Go modules file for dependency management.
- `Makefile`: Makefile to build and deploy the function.
- `host.json`: Configuration file for the Azure Functions host.
- `SendEmail/function.json`: Configuration file for the specific Azure Function.
- `README.md`: This documentation file.

## Deployment

1. **Install the Azure Functions Core Tools**:

    ```sh
    npm install -g azure-functions-core-tools@3
    ```

2. **Set Up Environment Variables**:

    Set up the following environment variables in your Azure Function App:
    _For the example should be:_
    - `SMTP_HOST`
    - `SMTP_PORT`
    - `SMTP_USER`
    - `SMTP_PASSWORD`
    - `SMTP_AUTH_METHOD`

3. **Build and Deploy the Function**:

    ```sh
    make deploy
    ```

4. **Test the Function**:

    Invoke the Azure Function with a JSON payload matching the `EmailRequest` structure.
