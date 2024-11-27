#!/bin/bash

set -exu

if [ ! -d "/teleporter/contracts" ]; then
    git clone https://github.com/ava-labs/teleporter /teleporter && \
    cd /teleporter && \
    git checkout $TELEPORTER_COMMIT
fi

if [ ! -f "/root/.foundry/bin/forge" ]; then
    cd /teleporter && ./scripts/install_foundry.sh
fi

export PATH="/root/.foundry/bin/:${PATH}"

cd /teleporter/contracts && forge build --extra-output-files=bin
