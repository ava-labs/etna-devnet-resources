> Warning! This readme might be outdated. Step numbers might not match the code, folder names, or scripts. But overall, the concepts should be the same.


## Add PoA validator to an existing L1

FIXME: this readme needs some love  

### 01. 👾 Initialize registration

Source code: [01_init_registration/add.go](./01_init_registration/add.go)

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


### 02. 📝 Register on P-chain

Source code: [02_register_on_p_chain/register.go](./02_register_on_p_chain/register.go)

Registers the validator on the P-chain by:

1. Loading the BLS info and signed Warp message from files
2. Creating a wallet using the validator manager key
3. Building and signing a P-chain transaction to register the L1 validator with:
   - 1 AVAX staking amount
   - BLS proof of possession
   - Signed Warp message containing validator details
4. Issuing the transaction to register the validator
5. Retrying up to 3 times if registration fails

### 03. 🏰 Complete registration

Source code: [03_complete_validator_registration/finish.go](./03_complete_validator_registration/finish.go)

Completes the validator registration process by:

1. Loading the validation ID from file
2. Generating a signed Warp message proving the validator was registered on the P-chain
3. Loading the validator manager contract address and private key
4. Calling `completeValidatorRegistration` on the validator manager contract with:
   - The signed Warp message proving P-chain registration
   - A registration index of 0
5. Waiting for the transaction to be confirmed

This final step links the P-chain registration with the validator manager contract, allowing the validator to participate in consensus.

### 04. 🚀 Start another node

Source code: [04_launch_nodes/launch.sh](./04_launch_nodes/launch.sh)

Prints a start command for another node. This command can be executed on either a remote or local machine. Example:

```bash
docker rm -f NodeID-9GEcBKKBr7RfEZPio4SLnkXNgwmqaPysd || true; \
docker run -d \
  --name NodeID-9GEcBKKBr7RfEZPio4SLnkXNgwmqaPysd \
  --network host \
  -e AVALANCHEGO_NETWORK_ID=fuji \
  -e AVALANCHEGO_HTTP_PORT=9652 \
  -e AVALANCHEGO_STAKING_PORT=9653 \
  -e AVALANCHEGO_TRACK_SUBNETS=Eqf7CViYbedHJm2cTtHeueAeVQdVjuXQshWxvYYuywaVtinrw \
  -e AVALANCHEGO_HTTP_ALLOWED_HOSTS=* \
  -e AVALANCHEGO_HTTP_HOST=0.0.0.0 \
  -e AVALANCHEGO_STAKING_TLS_CERT_FILE_CONTENT=LS0tLS1..... \
  -e AVALANCHEGO_STAKING_TLS_KEY_FILE_CONTENT=LS0tLS1..... \
  -e BLS_KEY_BASE64=LS0tLS1..... \
  -e AVALANCHEGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme \
  containerman17/avalanchego-subnetevm:v1.12.0_v0.6.12 ;\
  docker logs -f NodeID-9GEcBKKBr7RfEZPio4SLnkXNgwmqaPysd=

```
