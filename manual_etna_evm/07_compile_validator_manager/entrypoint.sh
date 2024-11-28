#!/bin/bash

set -exu

if [ ! -d "/teleporter/contracts" ]; then
    git clone https://github.com/ava-labs/teleporter /teleporter && \
    cd /teleporter && \
    git checkout $TELEPORTER_COMMIT
fi

cd /teleporter/contracts && forge build --extra-output-files=bin


# Extract ABI from the compiled JSON file
jq .abi /teleporter/out/PoAValidatorManager.sol/PoAValidatorManager.json > /teleporter/out/PoAValidatorManager.sol/PoAValidatorManager.abi
