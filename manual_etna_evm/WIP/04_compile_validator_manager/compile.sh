#!/bin/bash
set -exu -o pipefail

SCRIPT_DIR=$(dirname "$0")
TELEPORTER_COMMIT=790ccce873f9a904910a0f3ffd783436c920ce97

# Get current user and group IDs
CURRENT_UID=$(id -u)
CURRENT_GID=$(id -g)

docker build -t validator-manager-compiler --build-arg TELEPORTER_COMMIT=$TELEPORTER_COMMIT "$SCRIPT_DIR"
docker run -it --rm \
    -v "$SCRIPT_DIR/bindings":/bindings \
    -v "$SCRIPT_DIR/compiled":/compiled \
    -v "$SCRIPT_DIR/teleporter_src":/teleporter_src \
    -e TELEPORTER_COMMIT=$TELEPORTER_COMMIT \
    -e HOST_UID=$CURRENT_UID \
    -e HOST_GID=$CURRENT_GID \
    validator-manager-compiler
