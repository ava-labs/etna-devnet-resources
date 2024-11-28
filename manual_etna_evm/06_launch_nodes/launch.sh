#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

export CURRENT_UID=$(id -u)
export CURRENT_GID=$(id -g)
export TRACK_SUBNETS=$(cat "${SCRIPT_DIR}/../data/subnet.txt" | tr -d '\n')
export CHAIN_ID=$(cat "${SCRIPT_DIR}/../data/chain.txt" | tr -d '\n')

mkdir -p "${SCRIPT_DIR}/../data/chains/${CHAIN_ID}"
cp "${SCRIPT_DIR}/evm_debug_config.json" "${SCRIPT_DIR}/../data/chains/${CHAIN_ID}/config.json"

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" down || true
docker compose -f "${SCRIPT_DIR}/docker-compose.yml" up -d --build

# Add health check loop
echo "Waiting for subnet to become available..."
max_attempts=60
attempt=1
while [ $attempt -le $max_attempts ]; do
    if curl -s -X POST \
        -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
        "http://127.0.0.1:9650/ext/bc/${CHAIN_ID}/rpc" > /dev/null 2>&1; then
        echo "✅ Subnet is healthy and responding"
        exit 0
    fi
    
    echo "🌱 Subnet is still starting up (attempt $attempt of $max_attempts)"
    docker logs node0 2>&1 | tail -n 1
    
    sleep 2
    attempt=$((attempt + 1))
done

echo "❌ Subnet failed to become healthy after $max_attempts attempts"
exit 1

