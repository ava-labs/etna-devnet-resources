#!/bin/bash

SCRIPT_DIR=$(dirname "$0")

docker build -t validator-manager-compiler "$SCRIPT_DIR"
docker run -it --rm -v "$SCRIPT_DIR/teleporter":/teleporter -e TELEPORTER_COMMIT=790ccce validator-manager-compiler

cp -r "$SCRIPT_DIR/teleporter/out/PoAValidatorManager.sol" $SCRIPT_DIR/../data/
