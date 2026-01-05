use actix_web::{App, HttpServer, dev::Server};
use actix_web_prom::PrometheusMetricsBuilder;

use crate::domain::errors::HttpServerError;
use crate::transport::http;

pub fn run(port: u16) -> Result<Server, HttpServerError> {
    let prometheus = PrometheusMetricsBuilder::new("api")
        .endpoint("/metrics")
        .build()
        .map_err(|e| HttpServerError::Prometheus(e.to_string()))?;

    let server = HttpServer::new(move || {
        App::new()
            .wrap(prometheus.clone())
            .wrap(actix_web::middleware::Logger::default())
            .configure(http::routes::configure)
    })
    .bind(("0.0.0.0", port))?
    .run();

    Ok(server)
}
