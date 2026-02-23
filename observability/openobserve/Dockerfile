# Minimal runtime image for OpenObserve.
# The binary must be pre-downloaded into this directory before building.
# See download-openobserve.sh for fetching it.
FROM debian:trixie-slim AS runtime
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY openobserve /
RUN chmod +x /openobserve
RUN ["/openobserve", "init-dir", "-p", "/data/"]
CMD ["/openobserve"]
