#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
CMD="go run ./cmd/01_generate_keys/"
echo $CMD
$CMD

echo -e "\n💰 Checking balance\n"
CMD="go run ./cmd/02_check_balance/"
echo $CMD
$CMD

echo -e "\n🕸️  Creating subnet\n"
CMD="go run ./cmd/03_create_subnet/"
echo $CMD
$CMD

echo -e "\n🔗 Preparing chain\n"
CMD="go run ./cmd/04_prep_chain/"
echo $CMD
$CMD

echo -e "\n⛓️  Creating chain\n"
CMD="go run ./cmd/05_create_chain/"
echo $CMD
$CMD

echo -e "\n⚙️  Setting up cluster configs\n"
CMD="go run ./cmd/06_cluster_configs/"
echo $CMD
$CMD

echo -e "\n🚀 Launching nodes\n"
CMD="docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d --build"
echo $CMD
$CMD
