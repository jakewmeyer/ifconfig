use crate::config::Config;
use anyhow::{Context, Result};
use axum::{routing::get, Router};
use std::net::SocketAddr;
use tokio::signal;
use tower_default_headers::DefaultHeadersLayer;
use tower_http::compression::CompressionLayer;
use tower_http::cors::CorsLayer;
use tower_http::request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer};
use tower_http::trace::TraceLayer;
use tracing::info;

mod error;
mod routes;

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    info!("Signal received, starting graceful shutdown...");
}

pub async fn serve(config: Config) -> Result<()> {
    let addr: SocketAddr = format!("{}:{}", config.host, config.port).parse().unwrap();

    let app = Router::new()
        .route("/", get(routes::get_ip_plaintext))
        .route("/json", get(routes::get_ip_json))
        .route("/health", get(routes::health))
        .layer(CompressionLayer::new())
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::new())
        .layer(DefaultHeadersLayer::new(owasp_headers::headers()))
        .layer(SetRequestIdLayer::x_request_id(MakeRequestUuid))
        .layer(PropagateRequestIdLayer::x_request_id());

    info!("Listening on {}", addr);
    axum::Server::try_bind(&addr)?
        .serve(app.into_make_service())
        .with_graceful_shutdown(shutdown_signal())
        .await
        .context("Failed to start http server")
}
