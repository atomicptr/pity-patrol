FROM rust:1.93-alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY . /app

RUN if [[ "$TARGETARCH" = "arm64" ]]; then \
        rustup target add aarch64-unknown-linux-musl; \
        cargo build --release --target aarch64-unknown-linux-musl; \
        mv /app/target/aarch64-unknown-linux-musl/release/pity-patrol /app/pity-patrol; \
    else \
        rustup target add x86_64-unknown-linux-musl; \
        cargo build --release --target x86_64-unknown-linux-musl; \
        mv /app/target/x86_64-unknown-linux-musl/release/pity-patrol /app/pity-patrol; \
    fi

FROM alpine:latest

ENV PITY_PATROL_CONFIG="/app/config/config.toml"
ENV RUST_LOG="info"

WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /app/pity-patrol /app/pity-patrol

CMD ["/app/pity-patrol"]

