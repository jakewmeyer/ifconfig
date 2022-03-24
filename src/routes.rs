use actix_web::{get, HttpRequest, HttpResponse, Responder};
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::{net::IpAddr, str::FromStr};

use crate::errors::APIError;

#[derive(Debug, Serialize, Deserialize)]
struct IpResponse {
    ip: IpAddr,
}

fn parse_ip_from_request(req: &HttpRequest) -> Result<IpAddr, APIError> {
    if let Some(ip) = req.connection_info().realip_remote_addr() {
        match IpAddr::from_str(ip) {
            Ok(ip) => Ok(ip),
            Err(e) => Err(APIError::InvalidIp(e)),
        }
    } else {
        Err(APIError::NoIp)
    }
}

#[get("/")]
pub async fn get_ip_plaintext(req: HttpRequest) -> impl Responder {
    match parse_ip_from_request(&req) {
        Ok(ip) => HttpResponse::Ok().body(format!("{}", ip)),
        Err(e) => HttpResponse::BadRequest().body(e.to_string()),
    }
}

#[get("/json")]
pub async fn get_ip_json(req: HttpRequest) -> impl Responder {
    match parse_ip_from_request(&req) {
        Ok(ip) => HttpResponse::Ok().json(IpResponse { ip }),
        Err(e) => HttpResponse::BadRequest().body(e.to_string()),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::{test, web::Bytes, App};

    #[test]
    async fn test_parse_ip_from_request() {
        let input = "192.168.1.1";
        let expected = IpAddr::from_str(input).unwrap();
        let req = test::TestRequest::default()
            .insert_header(("X-Forwarded-For", input))
            .to_http_request();
        let res = parse_ip_from_request(&req);
        assert_eq!(res.unwrap(), expected);
    }

    #[test]
    async fn test_parse_ip_from_request_no_ip() {
        let req = test::TestRequest::default().to_http_request();
        let res = parse_ip_from_request(&req);
        assert_eq!(res.unwrap_err(), APIError::NoIp);
    }

    #[test]
    async fn test_parse_ip_from_request_invalid_ip() {
        let req = test::TestRequest::default()
            .insert_header(("X-Forwarded-For", "invalid-ip"))
            .to_http_request();
        let res = parse_ip_from_request(&req);
        assert!(matches!(res.unwrap_err(), APIError::InvalidIp(_)));
    }

    #[test]
    async fn test_get_ip_plaintext() {
        let input = "192.168.1.1";
        let app = test::init_service(App::new().service(get_ip_plaintext)).await;
        let req = test::TestRequest::get()
            .uri("/")
            .insert_header(("X-Forwarded-For", input))
            .to_request();
        let res = test::call_and_read_body(&app, req).await;
        assert_eq!(res, Bytes::from_static(input.as_bytes()));
    }

    #[test]
    async fn test_get_ip_plaintext_no_ip() {
        let app = test::init_service(App::new().service(get_ip_plaintext)).await;
        let req = test::TestRequest::get().uri("/").to_request();
        let res = test::call_and_read_body(&app, req).await;
        assert_eq!(res, Bytes::from_static("No IP address found".as_bytes()));
    }

    #[test]
    async fn test_get_ip_json() {
        let input = "192.168.1.1";
        let expected = IpAddr::from_str(input).unwrap();
        let app = test::init_service(App::new().service(get_ip_json)).await;
        let req = test::TestRequest::get()
            .uri("/json")
            .insert_header(("X-Forwarded-For", input))
            .to_request();
        let res: IpResponse = test::call_and_read_body_json(&app, req).await;
        assert_eq!(res.ip, expected);
    }

    #[test]
    async fn test_get_ip_json_no_ip() {
        let app = test::init_service(App::new().service(get_ip_json)).await;
        let req = test::TestRequest::get().uri("/json").to_request();
        let res = test::call_and_read_body(&app, req).await;
        assert_eq!(res, Bytes::from_static("No IP address found".as_bytes()));
    }
}
