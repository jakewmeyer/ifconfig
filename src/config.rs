#[derive(clap::Parser)]
pub struct Config {
    #[clap(long, env, default_value = "0.0.0.0")]
    pub host: String,

    #[clap(long, env, default_value = "7001")]
    pub port: String,
}
