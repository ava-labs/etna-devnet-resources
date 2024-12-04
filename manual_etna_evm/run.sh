# #!/bin/bash

# set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./03_create_subnet/
echo -e "\n🛠️ Using hardcoded smart contracts code\n"

echo -e "\n🧱 Generating genesis\n"
go run ./05_L1_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./06_create_chain/

echo -e "\n🚀 Launching nodes\n"
./07_launch_nodes/launch.sh 

echo -e "\n🔮 Converting chain into L1\n"
go run ./08_convert_chain/

echo -e "\n🚀 Restarting nodes\n"
./09_restart_nodes/restart.sh 

# echo -e "\n🎯 Activate ProposerVM fork\n"
# go run ./10_activate_proposer_vm/

# echo -e "\n📦 Deploy Validator Manager\n"
# go run ./07_depoly_validator_manager/

# echo -e "\n🔃 Restarting nodes\n"
# ./06_launch_nodes/launch.sh # Reuse the script to restart nodes

# echo -e "\n🔌 Initialize Validator Manager\n"
# go run ./11_validator_manager_initialize/ 

# echo -e "\n👥 Initialize validator set\n"
# go run ./12_initialize_validator_set


