#![deny(clippy::all)]

use actix_web::{
    middleware::{Compress, NormalizePath},
    App, HttpServer,
};
use miette::{Result, IntoDiagnostic};
use tracing_actix_web::TracingLogger;

mod errors;
mod routes;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv::dotenv().ok();
    tracing_subscriber::fmt::init();
    let host = std::env::var("HOST").unwrap_or_else(|_| "127.0.0.1".to_string());
    let port = std::env::var("PORT").unwrap_or_else(|_| "7100".to_string());
    let addr = format!("{}:{}", host, port);
    HttpServer::new(|| {
        App::new()
            .wrap(TracingLogger::default())
            .wrap(Compress::default())
            .wrap(NormalizePath::trim())
            .service(routes::get_ip_plaintext)
            .service(routes::get_ip_json)
    })
    .bind(addr).into_diagnostic()?
    .run()
    .await.into_diagnostic()?;
    Ok(())
}
