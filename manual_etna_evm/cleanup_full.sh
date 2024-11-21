#!/bin/bash

docker compose -f ./cmd/07_launch_nodes/docker-compose.yml down
rm -rf ./data
docker rm -f $(docker ps -a -q)
