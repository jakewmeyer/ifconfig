# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2022-07-14
### Added
- Support for the `x-real-ip` header, with `x-forwarded-for` as a backup option
- OWASP security headers by default
- Support for manual PORT and HOST overrides via cli args or environment variables
- Support for graceful shutdowns
- Support for CORS requests
- Support for request compression

### Changed
- Migrate web framework from [Actix Web](https://actix.rs/) to [Axum](https://github.com/tokio-rs/axum)
- Updated to Rust v1.62.0
