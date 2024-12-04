#!/bin/bash

set -euo pipefail
# Create data directory if it doesn't exist
mkdir -p /data/$NODE_NAME/db/

# Only extract if db directory is empty
if [ -z "$(ls -A /data/$NODE_NAME/db/)" ]; then
    tar -xvf "/fuji-latest.tar" -C /data/$NODE_NAME/db/
fi

/usr/local/bin/avalanchego --chain-config-dir=/data/chains --network-id=fuji --data-dir=/data/$NODE_NAME --plugin-dir=/plugins/ --http-port=9650 --staking-port=9651 --track-subnets=${TRACK_SUBNETS} --http-allowed-hosts=* --http-host=0.0.0.0
