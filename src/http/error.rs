use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use thiserror::Error;
use tracing::error;

/// Common Error type that allows us to return `Result` in handler functions.
///
/// User facing errors are defined with a corresponding status code and user
/// friendly message, while any one off `Anyhow` errors are automatically
/// considered to be 500, and the resulting error is only logged application
/// side for security purposes.
#[derive(Error, Debug)]
pub enum Error {
    #[error("No IP address found")]
    NotFound,
    #[error("An internal server error occurred")]
    Anyhow(#[from] anyhow::Error),
}

impl Error {
    /// Map defined errors to HTTP status codes
    pub fn status_code(&self) -> StatusCode {
        match self {
            Self::NotFound => StatusCode::NOT_FOUND,
            Self::Anyhow(_) => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }
}

/// Implement Axum's `IntoResponse` trait for our errors, so we can
/// return `Result` from handlers
impl IntoResponse for Error {
    fn into_response(self) -> Response {
        match self {
            Self::Anyhow(ref e) => {
                error!("Anyhow error: {:?}", e)
            }
            Error::NotFound => (),
        }
        (self.status_code(), self.to_string()).into_response()
    }
}
