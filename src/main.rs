#![deny(clippy::all)]
#![forbid(unsafe_code)]

use anyhow::Result;
use clap::Parser;
use ifconfig::config::Config;
use ifconfig::http;

#[tokio::main]
async fn main() -> Result<()> {
    dotenvy::dotenv().ok();

    tracing_subscriber::fmt::init();

    let config = Config::parse();

    http::serve(config).await?;

    Ok(())
}
