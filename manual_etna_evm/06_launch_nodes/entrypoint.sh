#!/bin/bash

set -euo pipefail


mkdir -p /data/node0/db/
tar -xvf "/fuji-latest.tar" -C /data/node0/db/

/usr/local/bin/avalanchego --chain-config-dir=/data/chains --network-id=fuji --data-dir=/data/node0 --plugin-dir=/plugins/ --http-port=9650 --staking-port=9651 --track-subnets=${TRACK_SUBNETS} --http-allowed-hosts=* --http-host=0.0.0.0
