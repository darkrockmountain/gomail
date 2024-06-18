# GoMail Examples

The .examples directory provides usage examples for the GoMail library.

## Overview

The examples directory contains various example applications demonstrating how to use the
gomail library with different email providers and deployment environments.

Each subdirectory contains examples for a specific environment or provider, showing how to
configure and use the library in different contexts.

## Example Usage

To run an example, navigate to the respective directory and follow the instructions provided
in the README.md file within that directory.

### Example (Using Google Cloud function example):

```sh
cd .examples/serverless/google_cloud_functions

# from https://cloud.google.com/sdk/docs/downloads-interactive
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

gcloud init

make deploy region=REGION #Replace REGION with the name of the Google Cloud region where you want to deploy your function (for example, us-west1).
```

### Examples Included:

- Serverless: Demonstrates usage in serverless environments like:
    - Google Cloud Functions,
    - Azure Functions, 
    - AWS Lambda, 
    - and DigitalOcean Functions.

Each example is self-contained and can be run independently to see the gomail library in action.