use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum APIError {
    #[error("Invalid IP address")]
    InvalidIp(#[from] std::net::AddrParseError),
    #[error("No IP address found")]
    NoIp,
}
