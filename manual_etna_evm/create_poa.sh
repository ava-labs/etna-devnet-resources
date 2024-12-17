# #!/bin/bash

set -euo pipefail

go run . generate-keys
go run . transfer-coins
go run . create-subnet
go run . generate-genesis
go run . create-chain
go run . convert-to-L1
go run . launch-node
go run . deploy-validator-manager --validator-type=poa
go run . validator-manager-init --validator-type=poa

go run . initialize-validator-set

sleep 30

go run . print-p-chain-info
go run . print-contract-logs
