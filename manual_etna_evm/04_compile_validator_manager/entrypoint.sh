#!/bin/bash

set -exu

# downlaod source code if not already present
if [ ! -d "/teleporter_src/contracts" ]; then
    git clone https://github.com/ava-labs/icm-contracts /teleporter_src 
    cd /teleporter_src && git submodule update --init --recursive
fi

# Add foundry to PATH
export PATH="/root/.foundry/bin/:${PATH}"

# Install foundry if not already installed
if ! command -v forge &> /dev/null; then
    cd /teleporter_src && ./scripts/install_foundry.sh
fi

# Build contracts
cd /teleporter_src/contracts && forge build --extra-output-files bin

# Extract ABI from the compiled JSON file
jq .abi /teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.json > /teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.abi

mkdir -p "/bindings/povalidatormanager"
abigen \
    --abi "/teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.abi" \
    --pkg povalidatormanager \
    --type PoAValidatorManager \
    --out "/bindings/povalidatormanager/PoAValidatorManager.go" \
    --bin "/teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.bin"

cp -r /teleporter_src/out/PoAValidatorManager.sol/*.json /compiled/
cp -r /teleporter_src/out/ValidatorMessages.sol/*.json /compiled/

chown -R $HOST_UID:$HOST_GID /bindings /compiled
