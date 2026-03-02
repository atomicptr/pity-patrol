FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG GIT_COMMIT

WORKDIR /app

COPY go.mod go.sum /app
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build \
        -ldflags="\
            -s -w \
            -X 'github.com/atomicptr/pity-patrol/pkgs/meta.Version=${VERSION}' \
            -X 'github.com/atomicptr/pity-patrol/pkgs/meta.GitCommit=${GIT_COMMIT}'" \
        -o pity-patrol cmd/pity-patrol/main.go

FROM alpine:latest

ENV PITY_PATROL_CONFIG="/app/config/config.toml"

WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /app/pity-patrol /app/pity-patrol

RUN chmod +x /app/pity-patrol

CMD ["/app/pity-patrol"]
