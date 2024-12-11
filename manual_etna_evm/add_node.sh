#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Initialize registration\n"
go run ./02_add_poa_validator/01_init_registration

echo -e "\n⛓️ Register on P-chain\n"
go run ./02_add_poa_validator/02_register_on_p_chain

echo -e "\n✅ Complete validator registration\n"
go run ./02_add_poa_validator/03_complete_validator_registration

echo -e "\n⏳ Waiting 60 seconds for transaction finalization...\n"
sleep 60

echo -e "\n📋 Reading contract logs again\n"
go run ./00_tools/check_validator_set

echo -e "\n🖥️ Launching node\n"
echo "Run the following command to start the node on a remote machine:"
go run ./02_add_poa_validator/04_print_launch_script
