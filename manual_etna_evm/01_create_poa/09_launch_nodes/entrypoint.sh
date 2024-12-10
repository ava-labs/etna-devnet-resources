#!/bin/bash

set -euo pipefail


# Function to convert ENV vars to flags
get_avalanchego_flags() {
    local flags=""
    # Loop through all environment variables
    while IFS='=' read -r name value; do
        # Check if variable starts with AVALANCHEGO_
        if [[ $name == AVALANCHEGO_* ]]; then
            # Convert AVALANCHEGO_DATA_DIR to --data-dir
            flag_name=$(echo "${name#AVALANCHEGO_}" | tr '[:upper:]' '[:lower:]' | tr '_' '-')
            flags+="--$flag_name=$value "
        fi
    done < <(env)
    echo "$flags"
}


# Get flags from environment variables
EXTRA_FLAGS=$(get_avalanchego_flags)

echo "Extra flags: $EXTRA_FLAGS"

# Create data directory if it doesn't exist
mkdir -p $AVALANCHEGO_DATA_DIR/db/

# This speeds up the node startup time
if [ -z "$(ls -A $AVALANCHEGO_DATA_DIR/db/)" ]; then
    wget -O $AVALANCHEGO_DATA_DIR/fuji-latest.tar https://avalanchego-public-database.avax-test.network/p-chain/avalanchego/data-tar/latest.tar
    tar -xvf "$AVALANCHEGO_DATA_DIR/fuji-latest.tar" -C $AVALANCHEGO_DATA_DIR/db/
fi

# Launch avalanchego with dynamic flags
/usr/local/bin/avalanchego $EXTRA_FLAGS
