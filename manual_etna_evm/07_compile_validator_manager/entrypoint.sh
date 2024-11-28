#!/bin/bash

set -exu

if [ ! -d "/teleporter/contracts" ]; then
    git clone https://github.com/ava-labs/teleporter /teleporter && \
    cd /teleporter && \
    git checkout $TELEPORTER_COMMIT
fi

cd /teleporter/contracts && forge build --extra-output-files=bin
