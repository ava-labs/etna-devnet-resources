## Guide to Etna L1s for Go Developers

This repository provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

**Requirements:**
- Fresh Docker installation (verify by running `docker compose ls` without any dashes)
- Go 1.22.10+

Run everything at once: `./create.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (preserves your keys)

Run `go run . validators` to print validators

Run `go run . logs 9650` to print contract logs from node0, `go run . logs 9652` to print contract logs from node1, etc.
