#!/bin/bash

set -euo pipefail


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
go run ./13_check_validator_set
