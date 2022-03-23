use actix_web::{
    middleware::{Compress, Logger, NormalizePath},
    App, HttpServer,
};
use dotenv::dotenv;
use env_logger::Env;

mod errors;
mod routes;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(Env::default().default_filter_or("info"));
    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .wrap(Compress::default())
            .wrap(NormalizePath::trim())
            .service(routes::get_ip_plaintext)
            .service(routes::get_ip_json)
    })
    .bind(("0.0.0.0", 7000))?
    .run()
    .await
}
