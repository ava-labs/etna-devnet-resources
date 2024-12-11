## Guide to Etna L1s for Go Developers

This repository provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

**Requirements:**
- Fresh Docker installation (verify by running `docker compose ls` without any dashes)
- Go 1.22.8+

- Part 1: [Create POA L1](./01_create_poa/README.md)
- Part 2: [Add POA validator](./02_add_poa_validator/README.md)

Run everything at once: `./add_node.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (preserves your keys)
