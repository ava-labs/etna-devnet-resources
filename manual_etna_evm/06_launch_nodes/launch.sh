#!/bin/bash

set -exuo pipefail

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

export CURRENT_UID=$(id -u)
export CURRENT_GID=$(id -g)
export TRACK_SUBNETS=$(cat "${SCRIPT_DIR}/../data/subnet.txt" | tr -d '\n')
export CHAIN_ID=$(cat "${SCRIPT_DIR}/../data/chain.txt" | tr -d '\n')

mkdir -p "${SCRIPT_DIR}/../data/chains/${CHAIN_ID}"
cp "${SCRIPT_DIR}/evm_debug_config.json" "${SCRIPT_DIR}/../data/chains/${CHAIN_ID}/config.json"

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" up -d --build
