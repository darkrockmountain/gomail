# Changelog
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
- No changes yet.

## [0.5.2] - 2024-06-03
### Added
- Migration of project repository from GitLab to GitHub.
- Automated CI/CD pipeline using GitHub Actions.

### Deprecated
- Deprecated all former versions to avoid issues with the migration.

## [0.5.1] - 2024-05-30
### Fixed
- Resolved issue with Microsoft365 email sender authentication.

## [0.5.0] - 2024-01-15
### Changed
- Major refactor of the email sending architecture to support plugin-based providers.

### Added
- Support for SparkPost email provider.
- Enhanced logging functionality.

## [0.4.1] - 2023-10-10
### Fixed
- Patched security vulnerability in SMTP email sender.

### Added
- Serverless examples.

## [0.4.0] - 2023-07-22
### Added
- Introduction of the SendGrid email provider.
- Initial support for OAuth2 authentication with Gmail.

## [0.3.4] - 2023-04-18
### Fixed
- Bug fix in Mailgun email sender related to attachment handling.

## [0.3.2] - 2023-01-25
### Added
- Feature to track email open and click rates.

## [0.3.1] - 2022-11-30
### Fixed
- Resolved intermittent issues with Postmark email sender.

## [0.3.0] - 2022-08-14
### Added
- Support for SES email provider.
- Added retry mechanism for failed email sends.

## [0.2.1] - 2022-04-12
### Added
- Initial support for Microsoft365 email provider.

## [0.2.0] - 2021-12-05
### Changed
- Improved error handling across all email providers.

## [0.1.2] - 2021-09-28
### Fixed
- Fixed a bug in Gmail email sender related to large email bodies.

## [0.1.1] - 2021-06-15
### Added
- Added support for Mandrill email provider.

## [0.1.0] - 2021-03-22
### Fixed
- Minor bug fixes and performance improvements.

## [0.0.3] - 2020-12-10
### Added
- Added support for Mailgun email provider.

## [0.0.2] - 2020-10-05
### Added
- Added support for Postmark email provider.
- Added detailed logging for email sending failures.

## [0.0.1] - 2020-05-10
### Added
- Support for Gmail email provider.

## [0.0.0] - 2020-01-20
### Added
- Initial release with support for SMTP email provider.
