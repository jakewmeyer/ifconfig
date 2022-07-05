use axum::{
    body::Body,
    http::{Request, StatusCode},
    response::IntoResponse,
    Json,
};
use miette::Result;
use serde::Serialize;
use std::net::IpAddr;

use crate::error::Error;

const X_REAL_IP: &str = "x-real-ip";
const X_FORWARDED_FOR: &str = "x-forwarded-for";

#[derive(Debug, Serialize)]
struct IpResponse {
    ip: Option<IpAddr>,
}

fn parse_ip_from_request(req: Request<Body>) -> Result<IpAddr, Error> {
    let headers = req.headers();

    if let Some(real_ip) = headers
        .get(X_REAL_IP)
        .and_then(|value| value.to_str().ok())
        .and_then(|s| s.parse::<IpAddr>().ok())
    {
        return Ok(real_ip);`
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

pub async fn get_ip_plaintext(req: Request<Body>) -> impl IntoResponse {
    match parse_ip_from_request(req) {
        Ok(ip) => (StatusCode::OK, format!("{}", ip)),
        Err(e) => (StatusCode::BAD_REQUEST, e.to_string()),
    }
}

pub async fn get_ip_json(req: Request<Body>) -> impl IntoResponse {
    match parse_ip_from_request(req) {
        Ok(ip) => (
            StatusCode::OK,
            Json(IpResponse {
                ip: Some(ip),
                error: None,
            }),
        ),
        Err(e) => (
            StatusCode::BAD_REQUEST,
            Json(IpResponse {
                ip: None,
                error: Some(e),
            }),
        ),
    }
}

// #[cfg(test)]
// mod tests {
//     use super::*;
//     use actix_web::{test, web::Bytes, App};
//     use pretty_assertions::assert_eq;

//     #[test]
//     async fn test_parse_ip_from_request() {
//         let input = "192.168.1.1";
//         let expected = IpAddr::from_str(input).unwrap();
//         let req = test::TestRequest::default()
//             .insert_header(("X-Forwarded-For", input))
//             .to_http_request();
//         let res = parse_ip_from_request(&req);
//         assert_eq!(res.unwrap(), expected);
//     }

//     #[test]
//     async fn test_parse_ip_from_request_no_ip() {
//         let req = test::TestRequest::default().to_http_request();
//         let res = parse_ip_from_request(&req);
//         assert_eq!(res.unwrap_err(), APIError::NoIp);
//     }

//     #[test]
//     async fn test_parse_ip_from_request_invalid_ip() {
//         let req = test::TestRequest::default()
//             .insert_header(("X-Forwarded-For", "invalid-ip"))
//             .to_http_request();
//         let res = parse_ip_from_request(&req);
//         assert!(matches!(res.unwrap_err(), APIError::InvalidIp(_)));
//     }

//     #[test]
//     async fn test_get_ip_plaintext() {
//         let input = "192.168.1.1";
//         let app = test::init_service(App::new().service(get_ip_plaintext)).await;
//         let req = test::TestRequest::get()
//             .uri("/")
//             .insert_header(("X-Forwarded-For", input))
//             .to_request();
//         let res = test::call_and_read_body(&app, req).await;
//         assert_eq!(res, Bytes::from_static(input.as_bytes()));
//     }

//     #[test]
//     async fn test_get_ip_plaintext_no_ip() {
//         let app = test::init_service(App::new().service(get_ip_plaintext)).await;
//         let req = test::TestRequest::get().uri("/").to_request();
//         let res = test::call_and_read_body(&app, req).await;
//         assert_eq!(res, Bytes::from_static("No IP address found".as_bytes()));
//     }

//     #[test]
//     async fn test_get_ip_json() {
//         let input = "192.168.1.1";
//         let expected = IpAddr::from_str(input).unwrap();
//         let app = test::init_service(App::new().service(get_ip_json)).await;
//         let req = test::TestRequest::get()
//             .uri("/json")
//             .insert_header(("X-Forwarded-For", input))
//             .to_request();
//         let res: IpResponse = test::call_and_read_body_json(&app, req).await;
//         assert_eq!(res.ip, expected);
//     }

//     #[test]
//     async fn test_get_ip_json_no_ip() {
//         let app = test::init_service(App::new().service(get_ip_json)).await;
//         let req = test::TestRequest::get().uri("/json").to_request();
//         let res = test::call_and_read_body(&app, req).await;
//         assert_eq!(res, Bytes::from_static("No IP address found".as_bytes()));
//     }
// }
