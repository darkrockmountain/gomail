# Google Cloud Functions Email Sender Example

This example demonstrates how to use the `EmailSender` interface in a Google Cloud Function using environment variables for configuration.

## Files

- `gcf_email_sender.go`: The main Go code for the Cloud Function.
- `go.mod`: Go modules file for dependency management.
- `Makefile`: Makefile to build and deploy the function.
- `README.md`: This documentation file.

## Deployment

1. **Install the Google Cloud SDK**:

    ```sh
    curl https://sdk.cloud.google.com | bash
    exec -l $SHELL
    gcloud init
    ```

2. **Set Up Environment Variables**:

    Set up the following environment variables in your Google Cloud Function:
    _For the example should be:_
    - `SMTP_HOST`
    - `SMTP_PORT`
    - `SMTP_USER`
    - `SMTP_PASSWORD`
    - `SMTP_AUTH_METHOD`

3. **Build and Deploy the Function**:

    ```sh
    make deploy region=REGION #Replace REGION with the name of the Google Cloud region where you want to deploy your function (for example, us-west1).
    ```

4. **Test the Function**:

    Invoke the Cloud Function with a JSON payload matching the `EmailRequest` structure.
