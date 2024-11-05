# AvalancheGo builder
FROM golang:1.22-bookworm AS avalanchego-builder

WORKDIR /app

RUN git clone https://github.com/ava-labs/avalanchego.git
WORKDIR /app/avalanchego
RUN git checkout v1.12.0-initial-poc.6
RUN ./scripts/build.sh

# Execution
FROM debian:bookworm-slim

# Create nobody group if it doesn't exist
RUN groupadd -r nobody || true

COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego

