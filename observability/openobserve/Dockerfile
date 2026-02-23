# Build OpenObserve from source
# Based on: https://github.com/openobserve/openobserve/blob/main/deploy/build/Dockerfile
#
# This is a heavy build (~15-30 min first time) but produces a native binary
# for whatever architecture you're running on (amd64/arm64).

ARG OPENOBSERVE_VERSION=main

# ---------- Stage 1: clone ----------
FROM alpine:3 AS source
ARG OPENOBSERVE_VERSION
RUN apk add --no-cache git
WORKDIR /src
RUN git clone --depth 1 --branch ${OPENOBSERVE_VERSION} \
    https://github.com/openobserve/openobserve.git .

# ---------- Stage 2: build web UI ----------
FROM node:24-trixie AS webbuilder
WORKDIR /web
COPY --from=source /src/web /web/
RUN npm install
RUN NODE_OPTIONS="--max-old-space-size=8192" npm run build

# ---------- Stage 3: build Rust backend ----------
FROM rust:1-slim-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler cmake make g++ pkg-config \
    libssl-dev libsasl2-dev libzstd-dev \
    && rm -rf /var/lib/apt/lists/*
RUN rustc --version
WORKDIR /openobserve
COPY --from=source /src /openobserve
COPY --from=webbuilder /web/dist web/dist
RUN cargo build --release

# ---------- Stage 4: runtime ----------
FROM debian:trixie-slim AS runtime
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY --from=builder /openobserve/target/release/openobserve /
RUN ["/openobserve", "init-dir", "-p", "/data/"]
CMD ["/openobserve"]
