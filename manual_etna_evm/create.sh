# #!/bin/bash

set -exuo pipefail

export L1_VALIDATOR_TYPE="poa"

echo "Building etnacli"
go build -o ./etnacli .

./etnacli generate-keys
./etnacli transfer-coins
./etnacli create-subnet
./etnacli generate-genesis
./etnacli create-chain
./etnacli convert-to-L1
./etnacli launch-node
./etnacli deploy-validator-manager --validator-type=${L1_VALIDATOR_TYPE}
./etnacli validator-manager-init --validator-type=${L1_VALIDATOR_TYPE}

./etnacli initialize-validator-set

sleep 30

./etnacli validators
./etnacli logs 9650
