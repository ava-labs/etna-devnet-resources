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
jq .abi /teleporter/out/PoSValidatorManager.sol/PoSValidatorManager.json > /teleporter/out/PoSValidatorManager.sol/PoSValidatorManager.abi

mkdir -p "/bindings/povalidatormanager"
mkdir -p "/bindings/posvalidatormanager"
abigen \
    --abi "/teleporter/out/PoAValidatorManager.sol/PoAValidatorManager.abi" \
    --pkg povalidatormanager \
    --type PoAValidatorManager \
    --out "/bindings/povalidatormanager/PoAValidatorManager.go" \
    --bin "/teleporter/out/PoAValidatorManager.sol/PoAValidatorManager.bin"

abigen \
    --abi "/teleporter/out/PoSValidatorManager.sol/PoSValidatorManager.abi" \
    --pkg posvalidatormanager \
    --type PoSValidatorManager \
    --out "/bindings/posvalidatormanager/PoSValidatorManager.go" \
    --bin "/teleporter/out/PoSValidatorManager.sol/PoSValidatorManager.bin"

cp -r /teleporter/out/PoAValidatorManager.sol/*.json /compiled/
cp -r /teleporter/out/PoSValidatorManager.sol/*.json /compiled/

chmod -R 777 /compiled /bindings
