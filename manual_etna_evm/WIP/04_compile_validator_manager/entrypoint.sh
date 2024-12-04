#!/bin/bash

set -exu

if [ ! -d "/teleporter_src/contracts" ]; then
    git clone https://github.com/ava-labs/teleporter /teleporter_src && \
    cd /teleporter_src && \
    git checkout $TELEPORTER_COMMIT
fi

cd /teleporter_src/contracts

forge build --via-ir

# Extract ABI from the compiled JSON file
jq .abi /teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.json > /teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.abi
jq .abi /teleporter_src/out/PoSValidatorManager.sol/PoSValidatorManager.json > /teleporter_src/out/PoSValidatorManager.sol/PoSValidatorManager.abi

mkdir -p "/bindings/povalidatormanager"
mkdir -p "/bindings/posvalidatormanager"
abigen \
    --abi "/teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.abi" \
    --pkg povalidatormanager \
    --type PoAValidatorManager \
    --out "/bindings/povalidatormanager/PoAValidatorManager.go" \
    --bin "/teleporter_src/out/PoAValidatorManager.sol/PoAValidatorManager.bin"

abigen \
    --abi "/teleporter_src/out/PoSValidatorManager.sol/PoSValidatorManager.abi" \
    --pkg posvalidatormanager \
    --type PoSValidatorManager \
    --out "/bindings/posvalidatormanager/PoSValidatorManager.go" \
    --bin "/teleporter_src/out/PoSValidatorManager.sol/PoSValidatorManager.bin"

cp -r /teleporter_src/out/PoAValidatorManager.sol/*.json /compiled/
cp -r /teleporter_src/out/PoSValidatorManager.sol/*.json /compiled/

chown -R $HOST_UID:$HOST_GID /bindings /compiled
