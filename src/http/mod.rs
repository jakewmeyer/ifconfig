//! http contains a serve function that constructs a new
//! Axum app from a Config and attempts to serve it

use crate::config::Config;
use anyhow::Result;
use axum::{routing::get, Router};
use std::net::SocketAddr;
use std::sync::Arc;
use tokio::net::TcpListener;
use tokio::{select, signal};
use tokio_util::sync::CancellationToken;
use tower_http::compression::CompressionLayer;
use tower_http::cors::CorsLayer;
use tower_http::request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer};
use tower_http::trace::TraceLayer;
use tracing::info;

mod error;
mod routes;

#[derive(Clone)]
pub struct ApiContext {
    pub shutdown: CancellationToken,
    pub config: Arc<Config>,
}

/// Creates a signal handler for graceful shutdown.
async fn shutdown_signal(ctx: ApiContext) {
    // Handle SIGINT
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    // Handle SIGTERM
    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    info!("Starting graceful server cleanup...");

    // Trigger cancellation token
    ctx.shutdown.cancel();

    info!("Starting hyper graceful shutdown...");
}

/// Create and serve an Axum server with pre-registered routes
/// and middleware
pub async fn serve(config: Config) -> Result<()> {
    let addr: SocketAddr = format!("{}:{}", config.host, config.port).parse().unwrap();

    let state = ApiContext {
        shutdown: CancellationToken::new(),
        config: Arc::new(config),
    };

    let app = Router::new()
        .route("/", get(routes::get_ip_plaintext))
        .route("/json", get(routes::get_ip_json))
        .route("/health", get(routes::health))
        .layer(TraceLayer::new_for_http())
        .layer(CompressionLayer::new())
        .layer(CorsLayer::new())
        .layer(PropagateRequestIdLayer::x_request_id())
        .layer(SetRequestIdLayer::x_request_id(MakeRequestUuid));

    let listener = TcpListener::bind(addr).await?;
    info!("Listening on http://{}", addr);

    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<SocketAddr>(),
    )
    .with_graceful_shutdown(shutdown_signal(state.clone()))
    .await?;

    Ok(())
}
