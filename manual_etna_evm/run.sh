# #!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./03_create_subnet/

echo -e "\n🛠️ Compiling smart contracts code\n"
./04_compile_validator_manager/compile.sh

echo -e "\n🧱 Generating genesis\n"
go run ./05_L1_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./06_create_chain/

echo -e "\n🚀 Launching nodes\n"
./07_launch_nodes/launch.sh "node0"

echo -e "\n🔮 Converting chain into L1\n"
go run ./08_convert_chain/

echo -e "\n🚀 Restarting nodes\n"
./07_launch_nodes/launch.sh "node0"

echo -e "\n🎯 Activate ProposerVM fork\n"
go run ./10_activate_proposer_vm/

echo -e "\n🔌 Initialize Validator Manager\n"
go run ./11_validator_manager_initialize/ 

echo -e "\n👥 Initialize validator set\n"
go run ./12_initialize_validator_set

echo -e "\n📄 Reading contract logs\n"
go run ./13_read_contract_logs

echo -e "\n🚀 Starting 1 more node\n"
./07_launch_nodes/launch.sh "node0 node1"

echo -e "\n👥 Add validator - initialize registration\n"
go run ./15_add_validator_init_registration/

echo -e "\n👥 Add validator - register on P-chain\n"
go run ./16_add_validator_register_on_p_chain/

echo -e "\n👥 Add validator - complete validator registration\n"
go run ./17_add_validator_complete_validator_registration/

echo -e "\n🎉 Everything is done! Waiting for 1 minute before reading contract logs again\n"
sleep 60

echo -e "\n📄 Reading contract logs again\n"
go run ./13_read_contract_logs
