#!/bin/bash

docker rm -f node0 node1 node2 node3 node4

sudo rm -rf ./data/node*/chainData
sudo rm -rf ./data/node*/logs
sudo rm -rf ./data/node*/db
sudo rm -rf ./data/node*/process.json

docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d

echo -e "\n🏥 Checking subnet health\n"
go run ./cmd/11_check_subnet_health/

echo -e "\n💸 Sending some test coins\n"
go run ./cmd/12_evm_transfer/
