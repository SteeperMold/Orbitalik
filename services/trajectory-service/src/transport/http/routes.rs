use actix_web::web;

use crate::transport::http::handlers::health::health_check;

pub fn configure(config: &mut web::ServiceConfig) {
    config.service(web::resource("/health").to(health_check));
}
