## Guide to Etna L1s for Go Developers

This repository provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

**Requirements:**
- Fresh Docker installation (verify by running `docker compose ls` without any dashes)
- Go 1.22.8+

- Part 1: [Create POA L1](./01_create_poa/README.md)
- Part 2: [Add POA validator](./02_add_poa_validator/README.md)

Run everything at once: `./add_node.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (preserves your keys)

### 14. 🚀 Start another node

Source code: none

```bash
./07_launch_nodes/launch.sh "node0 node1"
```

Starts another node to test adding a validator.

### 15. 👾 Add validator - initialize registration

//FIXME: Readme for the last 3 steps needs improvement

Source code: [15_add_validator_init_registration/add.go](./15_add_validator_init_registration/add.go)

Initializes the validator registration process by:

1. Getting the node ID and BLS key info from the node
2. Generating an expiry timestamp for the registration
3. Loading the validator manager key and contract address
4. Calling `initializeValidatorRegistration` on the validator manager contract with:
   - Node ID
   - BLS public key
   - Registration expiry
   - Balance owners (who can claim remaining balance)
   - Disable owners (who can disable the validator)
   - Validator weight
5. Generating and signing a Warp message containing the validator registration data
6. Saving the signed Warp message, BLS info, and validation ID to files


### 16. 📝 Add validator - register on P-chain

Source code: [16_add_validator_register_on_p_chain/register.go](./16_add_validator_register_on_p_chain/register.go)

Registers the validator on the P-chain by:

1. Loading the BLS info and signed Warp message from files
2. Creating a wallet using the validator manager key
3. Building and signing a P-chain transaction to register the L1 validator with:
   - 1 AVAX staking amount
   - BLS proof of possession
   - Signed Warp message containing validator details
4. Issuing the transaction to register the validator
5. Retrying up to 3 times if registration fails

### 17. 🏰 Add validator - complete registration

Source code: [17_add_validator_complete_validator_registration/finish.go](./17_add_validator_complete_validator_registration/finish.go)

Completes the validator registration process by:

1. Loading the validation ID from file
2. Generating a signed Warp message proving the validator was registered on the P-chain
3. Loading the validator manager contract address and private key
4. Calling `completeValidatorRegistration` on the validator manager contract with:
   - The signed Warp message proving P-chain registration
   - A registration index of 0
5. Waiting for the transaction to be confirmed

This final step links the P-chain registration with the validator manager contract, allowing the validator to participate in consensus.
