# # # # # # # # # # # # # # 
# AvalancheGo builder
# # # # # # # # # # # # # # 

FROM golang:1.22-bookworm AS avalanchego-builder

WORKDIR /app

RUN git clone https://github.com/ava-labs/avalanchego.git && cd avalanchego && git checkout v1.12.0-fuji && ./scripts/build.sh

# # # # # # # # # # # # # # 
# Subnet-EVM builder
# # # # # # # # # # # # # # 

FROM golang:1.22-bookworm AS subnet-evm-builder

RUN git clone https://github.com/ava-labs/subnet-evm.git /app/subnet-evm && cd /app/subnet-evm && git checkout v0.6.12

RUN cd /app/subnet-evm && go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin

# # # # # # # # # # # # # # 
# Execution
# # # # # # # # # # # # # # 

FROM debian:bookworm-slim

# Create nobody group if it doesn't exist
RUN groupadd -r nobody || true

COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego
COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

