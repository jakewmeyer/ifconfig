#![deny(clippy::all)]

use std::net::SocketAddr;

use anyhow::Result;
use axum::{routing::get, Router};
use tracing::info;

mod error;
mod routes;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv::dotenv().ok();
    tracing_subscriber::fmt::init();
    let host = std::env::var("HOST").unwrap_or_else(|_| "127.0.0.1".to_string());
    let port = std::env::var("PORT").unwrap_or_else(|_| "7000".to_string());
    let addr: SocketAddr = format!("{}:{}", host, port).parse().unwrap();

    let app = Router::new()
        .route("/", get(routes::get_ip_plaintext))
        .route("/json", get(routes::get_ip_json));

    info!("Listening on {}", addr);
    axum::Server::try_bind(&addr)?
        .serve(app.into_make_service())
        .await?;
    Ok(())
}
