#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
TELEPORTER_COMMIT=790ccce873f9a904910a0f3ffd783436c920ce97

docker build -t validator-manager-compiler --build-arg TELEPORTER_COMMIT=$TELEPORTER_COMMIT "$SCRIPT_DIR"
docker run -it --rm -v "$SCRIPT_DIR/teleporter":/teleporter -e TELEPORTER_COMMIT=$TELEPORTER_COMMIT validator-manager-compiler

cp -r "$SCRIPT_DIR/teleporter/out/PoAValidatorManager.sol" $SCRIPT_DIR/

mkdir -p "$SCRIPT_DIR/bindings/povalidatormanager"

go run github.com/ava-labs/subnet-evm/cmd/abigen@v0.6.12 \
    --abi "$SCRIPT_DIR/PoAValidatorManager.sol/PoAValidatorManager.abi" \
    --pkg povalidatormanager \
    --type PoAValidatorManager \
    --out "$SCRIPT_DIR/bindings/povalidatormanager/PoAValidatorManager.go" \
    --bin "$SCRIPT_DIR/teleporter/out/PoAValidatorManager.sol/PoAValidatorManager.bin"
