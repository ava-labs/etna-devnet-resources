#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
echo "go run ./cmd/01_generate_keys/generate.go"
go run ./cmd/01_generate_keys/generate.go

echo -e "\n💰 Checking balance\n"
echo "go run ./cmd/02_check_balance/balance.go"
go run ./cmd/02_check_balance/balance.go

echo -e "\n🕸️  Creating subnet\n"
echo "go run ./cmd/03_create_subnet/create.go"
go run ./cmd/03_create_subnet/create.go

echo -e "\n🔗 Preparing chain\n"
echo "go run ./cmd/04_prep_chain/prep.go"
go run ./cmd/04_prep_chain/prep.go
