use std::{
    pin::Pin,
    task::{Context, Poll},
};
use tokio::time::Instant;
use tower::{Layer, Service};

#[derive(Clone, Default)]
pub struct LoggingMiddlewareLayer {}

impl<S> Layer<S> for LoggingMiddlewareLayer {
    type Service = LoggingMiddleware<S>;

    fn layer(&self, service: S) -> Self::Service {
        LoggingMiddleware { inner: service }
    }
}

#[derive(Clone)]
pub struct LoggingMiddleware<S> {
    inner: S,
}

type BoxFuture<'a, T> = Pin<Box<dyn Future<Output = T> + Send + 'a>>;

impl<S, ReqBody, ResBody> Service<http::Request<ReqBody>> for LoggingMiddleware<S>
where
    S: Service<http::Request<ReqBody>, Response = http::Response<ResBody>> + Clone + Send + 'static,
    S::Future: Send + 'static,
    ReqBody: Send + 'static,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: http::Request<ReqBody>) -> Self::Future {
        let clone = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, clone);

        let method_path = req.uri().path().to_string();
        let start = Instant::now();

        Box::pin(async move {
            let result = inner.call(req).await;

            let duration = start.elapsed();

            match &result {
                Ok(resp) => {
                    let status = resp.status().as_u16();
                    tracing::info!(
                        "gRPC call succeeded: method={method_path} status={status} duration={duration:?}"
                    );
                }
                Err(_) => {
                    tracing::info!("gRPC call failed: method={method_path} duration={duration:?}");
                }
            }

            result
        })
    }
}
