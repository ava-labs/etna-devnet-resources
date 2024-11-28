#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
TELEPORTER_COMMIT=790ccce873f9a904910a0f3ffd783436c920ce97

docker build -t validator-manager-compiler --build-arg TELEPORTER_COMMIT=$TELEPORTER_COMMIT "$SCRIPT_DIR"
docker run -it --rm -v "$SCRIPT_DIR/bindings":/bindings -v "$SCRIPT_DIR/teleporter":/teleporter -e TELEPORTER_COMMIT=$TELEPORTER_COMMIT validator-manager-compiler

