#!/bin/bash
set -eu -o pipefail

SCRIPT_DIR=$(dirname "$0")
ICM_COMMIT=$(grep "github.com/ava-labs/icm-contracts" "$SCRIPT_DIR/../go.mod" | cut -d'-' -f5)
SUBNET_EVM_VERSION=$(grep "github.com/ava-labs/subnet-evm" "$SCRIPT_DIR/../go.mod" | cut -d' ' -f2)

echo "ICM_COMMIT: $ICM_COMMIT"
echo "SUBNET_EVM_VERSION: $SUBNET_EVM_VERSION"
echo "Warning! If those versions don't look like what you expect, check the first few lines of $0"

# Get current user and group IDs
CURRENT_UID=$(id -u)
CURRENT_GID=$(id -g)

docker build -t validator-manager-compiler --build-arg SUBNET_EVM_VERSION=$SUBNET_EVM_VERSION --build-arg ICM_COMMIT=$ICM_COMMIT "$SCRIPT_DIR"
docker run -it --rm \
    -v "$SCRIPT_DIR/bindings":/bindings \
    -v "$SCRIPT_DIR/compiled":/compiled \
    -v "$SCRIPT_DIR/teleporter_src":/teleporter_src \
    -e ICM_COMMIT=$ICM_COMMIT \
    -e HOST_UID=$CURRENT_UID \
    -e HOST_GID=$CURRENT_GID \
    validator-manager-compiler
