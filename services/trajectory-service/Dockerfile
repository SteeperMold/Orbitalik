FROM rust:1.89 AS cargo-tools

RUN cargo install sccache --version ^0.7 && \
    cargo install cargo-chef --version ^0.1

ENV RUSTC_WRAPPER=sccache SCCACHE_DIR=/sccache

FROM cargo-tools AS planner
WORKDIR /app

COPY /services/trajectory-service .

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=$SCCACHE_DIR,sharing=locked \
    cargo chef prepare --recipe-path recipe.json

FROM cargo-tools AS builder
WORKDIR /app

COPY --from=planner /app/recipe.json recipe.json

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=$SCCACHE_DIR,sharing=locked \
    cargo chef cook --release --recipe-path recipe.json

RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

COPY /services/trajectory-service .

COPY /proto /app/proto

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=$SCCACHE_DIR,sharing=locked \
    cargo build --release

FROM debian:bookworm-slim AS runtime
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/trajectory-service /usr/local/bin/trajectory-service

EXPOSE 8080

CMD ["trajectory-service"]
