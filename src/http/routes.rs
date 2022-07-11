use axum::{
    body::Body,
    http::{Request, StatusCode},
    response::IntoResponse,
    Json,
};
use serde::Serialize;
use std::net::IpAddr;

use crate::http::error::Error;

const X_REAL_IP: &str = "x-real-ip";
const X_FORWARDED_FOR: &str = "x-forwarded-for";

#[derive(Debug, Serialize)]
struct IpResponse {
    ip: IpAddr,
}

fn parse_ip_from_request(req: Request<Body>) -> Result<IpAddr, Error> {
    let headers = req.headers();

    if let Some(real_ip) = headers
        .get(X_REAL_IP)
        .and_then(|value| value.to_str().ok())
        .and_then(|s| s.parse::<IpAddr>().ok())
    {
        return Ok(real_ip);
    }

    if let Some(forwarded_for) = headers
        .get(X_FORWARDED_FOR)
        .and_then(|value| value.to_str().ok())
        .and_then(|s| s.split(',').find_map(|s| s.trim().parse::<IpAddr>().ok()))
    {
        return Ok(forwarded_for);
    }

    Err(Error::NotFound)
}

pub async fn get_ip_plaintext(req: Request<Body>) -> Result<impl IntoResponse, Error> {
    let ip = parse_ip_from_request(req)?;
    Ok(format!("{}", ip))
}

pub async fn get_ip_json(req: Request<Body>) -> Result<impl IntoResponse, Error> {
    let ip = parse_ip_from_request(req)?;
    Ok(Json(IpResponse { ip }))
}

pub async fn health() -> impl IntoResponse {
    StatusCode::OK
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use pretty_assertions::assert_eq;
    use std::str::FromStr;

    #[tokio::test]
    async fn parse_ip_from_x_forwarded_for() {
        let input = "192.168.1.1";
        let expected = IpAddr::from_str(input).unwrap();
        let req = Request::builder()
            .uri("/")
            .header(X_FORWARDED_FOR, input)
            .body(Body::empty())
            .unwrap();
        let res = parse_ip_from_request(req);
        assert_eq!(res.unwrap(), expected);
    }

    #[tokio::test]
    async fn parse_ip_from_x_real_ip() {
        let input = "192.168.1.3";
        let expected = IpAddr::from_str(input).unwrap();
        let req = Request::builder()
            .uri("/")
            .header(X_FORWARDED_FOR, "192.168.1.1,192.168.1.2")
            .header(X_REAL_IP, input)
            .body(Body::empty())
            .unwrap();
        let res = parse_ip_from_request(req);
        assert_eq!(res.unwrap(), expected);
    }
}
