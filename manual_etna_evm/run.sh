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

echo -e "\n🚀 Launching nodes\n"
./06_launch_nodes/launch.sh

echo -e "\n🛠️ Compiling validator manager\n"
./07_compile_validator_manager/compile.sh

echo -e "\n📦 Deploy Validator Manager\n"
go run ./08_depoly_validator_manager/

echo -e "\n🔮 Converting chain into L1\n"
go run ./09_convert_chain/

# echo -e "\n🔃 Restarting nodes\n"
# ./06_launch_nodes/launch.sh

# echo -e "\n🏥 Checking subnet health\n"
# go run ./09_check_subnet_health/

# echo -e "\n💸 Sending some test coins\n"
# go run ./10_evm_transfer/

# echo -e "\n🎯 Activate ProposerVM fork\n"
# go run ./11_activate_proposer_vm/
