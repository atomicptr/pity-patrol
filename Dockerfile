FROM rust:1.93-alpine AS builder

WORKDIR /app

COPY . /app
RUN cargo build --release --target x86_64-unknown-linux-musl

FROM alpine:latest

ENV PITY_PATROL_CONFIG="/app/config/config.toml"
ENV RUST_LOG="info"

WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/pity-patrol /app/pity-patrol

CMD ["/app/pity-patrol"]

