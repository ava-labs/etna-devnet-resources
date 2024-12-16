# #!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./01_create_poa/01_generate_keys/

echo -e "\n💰 Transferring AVAX between C and P chains\n" 
go run ./01_create_poa/02_transfer_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./01_create_poa/03_create_subnet/

echo -e "\n🧱 Generating genesis\n"
go run ./01_create_poa/05_L1_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./01_create_poa/06_create_chain/

echo -e "\n🔮 Converting chain into L1\n"
go run ./01_create_poa/07_convert_chain/

echo -e "\n🚀 Launching node0\n"
for i in {1..3}; do
  if ./01_create_poa/08_launch_nodes/launch.sh "node0"; then
    break
  elif [ $i -eq 3 ]; then
    echo "Failed to launch node0 after 3 attempts"
    exit 1
  else
    echo "Attempt $i failed, retrying..."
    sleep 5
  fi
done

echo -e "\n📝 Deploying Validator Manager\n"
go run ./01_create_poa/09_deploy_validator_manager/

echo -e "\n🔌 Initialize Validator Manager\n"
go run ./01_create_poa/10_validator_manager_initialize/

echo -e "\n⏳ Waiting for P-chain transactions to be mined...\n"
sleep 30

echo -e "\n👥 Initialize validator set\n"
go run ./01_create_poa/11_initialize_validator_set/

echo -e "\n⏳ Waiting for P-chain transactions to be mined...\n"
sleep 30

echo -e "\n📄 Reading contract logs\n"
go run ./00_tools/check_validator_set
