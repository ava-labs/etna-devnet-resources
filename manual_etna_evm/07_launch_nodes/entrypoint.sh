#!/bin/bash

set -euo pipefail
# Create data directory if it doesn't exist
mkdir -p /data/node0/db/

# Only extract if db directory is empty
if [ -z "$(ls -A /data/node0/db/)" ]; then
    tar -xvf "/fuji-latest.tar" -C /data/node0/db/
fi

/usr/local/bin/avalanchego --chain-config-dir=/data/chains --network-id=fuji --data-dir=/data/node0 --plugin-dir=/plugins/ --http-port=9650 --staking-port=9651 --track-subnets=${TRACK_SUBNETS} --http-allowed-hosts=* --http-host=0.0.0.0
