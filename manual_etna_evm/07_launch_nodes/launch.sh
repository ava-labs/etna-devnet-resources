#!/bin/bash

set -exuo pipefail

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

export CURRENT_UID=$(id -u)
export CURRENT_GID=$(id -g)
export TRACK_SUBNETS=$(cat "${SCRIPT_DIR}/../data/subnet.txt" | tr -d '\n')

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" up -d --build
