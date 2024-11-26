#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./03_create_subnet/

echo -e "\n🧱 Generating genesis\n"
go run ./04_L1_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./05_create_chain/

echo -e "\n🏗️  Setting up node configs\n"
go run ./06_node_configs/

echo -e "\n🚀 Launching nodes\n"
./07_launch_nodes/launch.sh

echo -e "\n🔄 Converting chain\n"
go run ./08_convert_chain/

echo -e "\n🚀 Starting nodes again with a new subnet\n"
./07_launch_nodes/launch.sh

echo -e "\n🏥 Checking subnet health\n"
go run ./11_check_subnet_health/

# echo -e "\n💸 Sending some test coins\n"
# go run ./12_evm_transfer/

# echo -e "\n🔄 Waiting for the transaction to be included\n"
# sleep 30

# echo -e "\n🔄 Initializing PoA\n"
# go run ./13_init_poa/
