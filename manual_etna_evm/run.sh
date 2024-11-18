#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./cmd/01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./cmd/02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./cmd/03_create_subnet/

echo -e "\n🧱 Generating genesis\n"
go run ./cmd/04_gen_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./cmd/05_create_chain/

echo -e "\n🏗️  Setting up node configs\n"
go run ./cmd/06_node_configs/

echo -e "\n🚀 Launching nodes\n"
export CURRENT_UID=$(id -u)
export CURRENT_GID=$(id -g)
docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d --build

echo -e "\n🔄 Converting chain\n"
go run ./cmd/08_convert_chain/

echo -e "\n🔄 Updating node configs\n"
go run ./cmd/09_update_configs/

echo -e "\n🚀 Stopping nodes\n"
docker compose -f ./cmd/07_launch_nodes/docker-compose.yml down

echo -e "\n🚀 Starting nodes again with a new subnet\n"
docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d

echo -e "\n🏥 Checking subnet health\n"
go run ./cmd/11_check_subnet_health/

echo -e "\n💸 Sending some test coins\n"
go run ./cmd/12_evm_transfer/

echo -e "\n🔄 Initializing PoA\n"
go run ./cmd/13_init_poa/
