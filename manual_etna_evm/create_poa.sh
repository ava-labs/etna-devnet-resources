# #!/bin/bash

set -euo pipefail

export L1_VALIDATOR_TYPE="pos-native"

echo "Building etnacli"
go build -o ./etnacli .

./etnacli generate-keys
./etnacli transfer-coins
./etnacli create-subnet
./etnacli generate-genesis
./etnacli create-chain
./etnacli convert-to-L1
./etnacli launch-node
./etnacli deploy-validator-manager 
./etnacli validator-manager-init

./etnacli initialize-validator-set

sleep 30

./etnacli print-p-chain-info
./etnacli print-contract-logs
