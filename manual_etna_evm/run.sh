#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
CMD="go run ./cmd/01_generate_keys/generate.go"
echo $CMD
$CMD

echo -e "\n💰 Checking balance\n"
CMD="go run ./cmd/02_check_balance/balance.go"
echo $CMD
$CMD

echo -e "\n🕸️  Creating subnet\n"
CMD="go run ./cmd/03_create_subnet/create.go"
echo $CMD
$CMD

echo -e "\n🔗 Preparing chain\n"
CMD="go run ./cmd/04_prep_chain/prep.go"
echo $CMD
$CMD

echo -e "\n🚀 Launching nodes\n"
CMD="docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d --build"
echo $CMD
$CMD
