//! Application config parameters
//! 
//! Config args can be passed through the command line, or passed
//! via environment variables. Dotenv support is included in main.
#[derive(clap::Parser)]
pub struct Config {
    /// The host the server will bind to, any valid
    /// IpAddr will suffice
    #[clap(long, env, default_value = "0.0.0.0")]
    pub host: String,

    // The port the server will bind to
    #[clap(long, env, default_value = "7000")]
    pub port: String,
}
