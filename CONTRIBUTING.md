# Contributing to GoMail - Email Library

Thank you for considering contributing to our project! We welcome contributions of all kinds, including bug reports, feature suggestions, and code improvements. To make the contribution process smooth for everyone, please follow these guidelines.

## Table of Contents

- [How to Report a Bug](#how-to-report-a-bug)
- [How to Suggest a Feature](#how-to-suggest-a-feature)
- [How to Set Up Your Development Environment](#how-to-set-up-your-development-environment)
- [Branch Naming Conventions](#branch-naming-conventions)
- [Coding Guidelines](#coding-guidelines)
- [How to Submit a Pull Request](#how-to-submit-a-pull-request)
- [Review and Approval Process](#review-and-approval-process)
- [Community Guidelines](#community-guidelines)

## How to Report a Bug

1. **Search Existing Issues**: Before reporting a new bug, check if it has already been reported.
2. **Open a New Issue**: If it’s a new bug, open an issue and include:
   - A clear and descriptive title.
   - Steps to reproduce the bug.
   - Expected and actual results.
   - Any relevant logs, screenshots, or other information.
   
   Use our [Bug Report Template](./.github/ISSUE_TEMPLATE/bug_report.md) for reference.

## How to Suggest a Feature

1. **Search Existing Issues**: Before suggesting a new feature, check if it has already been suggested.
2. **Open a New Issue**: If it’s a new feature, open an issue and include:
   - A clear and descriptive title.
   - A detailed description of the feature.
   - Any potential use cases or benefits.

   Use our [Feature Request Template](./.github/ISSUE_TEMPLATE/feature_request.md) for reference.

## How to Set Up Your Development Environment

1. **Fork the Repository**: Fork the repository to your GitHub account.
2. **Clone the Repository**: Clone your forked repository to your local machine.
    ```bash
    git clone github.com/darkrockmountain/gomail.git
    ```
3. **Navigate to the Project Directory**:
    ```bash
    cd gomail
    ```
4. **Install Dependencies**:
    ```bash
    go mod tidy
    ```
5. **Run Tests**: Ensure everything is working by running the tests.
    ```bash
    go test ./...
    ```

## Branch Naming Conventions

To keep the repository organized, please follow these branch naming conventions:
- **Feature branches**: `feature/your-feature-name`
- **Bugfix branches**: `bugfix/your-bug-name`
- **Documentation branches**: `docs/your-doc-name`
- **Hotfix branches**: `hotfix/your-hotfix-name`

## Coding Guidelines

1. **Code Style**: Follow the standard Go coding style. Run `go fmt` before committing your code.
2. **Documentation**: Ensure all public methods and packages have clear comments and documentation.
3. **Testing**: Write tests for new features and bug fixes. Ensure existing tests pass before submitting a pull request.
4. **Commits**: Write clear and concise commit messages. Follow the convention:
    ```
    type(scope): description
    ```
    Example:
    ```
    feat(providers): add support for new email provider
    ```

## How to Submit a Pull Request

1. **Create a Branch**: Create a new branch for your feature or bugfix.
    ```bash
    git checkout -b feature/my-new-feature
    ```
2. **Commit Your Changes**: Commit your changes with a clear message.
    ```bash
    git commit -m "feat: add new feature"
    ```
3. **Push to Your Fork**: Push your changes to your forked repository.
    ```bash
    git push origin feature/my-new-feature
    ```
4. **Open a Pull Request**: Go to the original repository on GitHub and open a pull request. Provide a detailed description of your changes.

   **Important**: Ensure that the pull request is made against the `develop` branch, not the `main` branch. Otherwise, it will be automatically rejected.
   
   Use our [Pull Request Template](./.github/PULL_REQUEST_TEMPLATE.md) for reference.

## Review and Approval Process

1. **Review**: Once a pull request is submitted, it will be reviewed by project maintainers. They will check for code quality, adherence to guidelines, and overall impact.
2. **Approval**: If the pull request meets all requirements, it will be approved and merged. If changes are needed, feedback will be provided for updates.
3. **Timeline**: Reviews typically take up to a few days. Please be patient as maintainers work through the review process.

## Community Guidelines

- **Respect**: Be respectful and considerate in your interactions with others.
- **Inclusivity**: Ensure your contributions are inclusive and accessible to a diverse audience.
- **Collaboration**: Work collaboratively and be open to feedback and suggestions.
- **Confidentiality**: Avoid publishing others' private information, such as (but not limited to) a physical or email address, without their explicit permission.
- **Code of Conduct**: By participating in this project, you agree to abide by our [Code of Conduct](./CODE_OF_CONDUCT.md).

By following these guidelines, you will help us maintain a welcoming and productive environment for all contributors. Thank you for your contributions!
