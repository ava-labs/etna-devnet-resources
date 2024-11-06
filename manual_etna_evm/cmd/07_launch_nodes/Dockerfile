# # # # # # # # # # # # # # 
# AvalancheGo builder
# # # # # # # # # # # # # # 

FROM golang:1.22-bookworm AS avalanchego-builder

WORKDIR /app

RUN git clone https://github.com/ava-labs/avalanchego.git
WORKDIR /app/avalanchego
RUN git checkout v1.12.0-initial-poc.6
RUN ./scripts/build.sh

# # # # # # # # # # # # # # 
# Subnet-EVM builder
# # # # # # # # # # # # # # 

FROM golang:1.22-bookworm AS subnet-evm-builder

WORKDIR /app
RUN git clone https://github.com/ava-labs/subnet-evm.git
WORKDIR /app/subnet-evm
RUN git checkout v0.6.11

RUN go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin

# # # # # # # # # # # # # # 
# Execution
# # # # # # # # # # # # # # 

FROM debian:bookworm-slim

# Create nobody group if it doesn't exist
RUN groupadd -r nobody || true

COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego
COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

