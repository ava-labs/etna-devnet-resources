## Guide to Etna L1s for Go Developers

This repository provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

**Requirements:**
- Fresh Docker installation (verify by running `docker compose ls` without any dashes)
- Go 1.22.8+

Run everything at once: `./run.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (preserves your keys)

> Note: This guide uses simplified, linear code with hardcoded values to demonstrate concepts clearly. Not intended as a library.

For up to date steps, check the [./run.sh](./run.sh) file.

### 1. 🔑 Generating Keys

Source code: [01_generate_keys/generate.go](./01_generate_keys/generate.go)

Generates a validator manager private key if you don't have one yet.

The key method you'll need is `secp256k1.NewPrivateKey()` from package `github.com/ava-labs/avalanchego/utils/crypto/secp256k1`.


### 2. 💰 Transfer AVAX between C and P chains

Source code: [02_transfer_balance/balance.go](./02_transfer_balance/balance.go)

- Checks your C-chain and P-chain balance
- Attempts to export all C-Chain funds to P-chain
- If insufficient funds, directs you to the [Fuji faucet](https://test.core.app/tools/testnet-faucet/?subnet=c&token=c)

This provides a good example of checking balances and transferring AVAX between C and P chains.

> There is a bug in counting the amounts somewhere. Running this step 2 times would solve it.


### 3. 🕸️ Creating subnet

Source code: [03_create_subnet/create.go](./03_create_subnet/create.go)

Creating a subnet requires only an owner. Here's how it works:

```golang
owner := &secp256k1fx.OutputOwners{
    Locktime:  0,
    Threshold: 1,
    Addrs:     []ids.ShortID{subnetOwner},
}

createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
```

A subnet is a group of validators that agree to validate the same chains. A chain can only belong to one subnet, but a subnet can have multiple chains. Each validator must validate all chains within their subnet. Validators can participate in multiple subnets simultaneously.

### 4. 📝 Compiling smart contracts code

Source code: [04_compile_validator_manager/compile.sh](./04_compile_validator_manager/compile.sh), [04_compile_validator_manager/entrypoint.sh](./04_compile_validator_manager/entrypoint.sh)


Simplified workflow:
```bash
git clone https://github.com/ava-labs/icm-contracts /teleporter_src 
cd /teleporter_src
git submodule update --init --recursive
./scripts/install_foundry.sh
cd /teleporter_src/contracts && forge build
```
### 5. 🧱 Generating Genesis

Source code: [05_L1_genesis/genesis.go](./05_L1_genesis/genesis.go)

Here we generate the genesis for our new L1. We will include it in a P-chain create chain transaction in the next step.

Note: Don't confuse your L1 genesis with the Avalanche Fuji genesis.

Read more about genesis here: [https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1](https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1).

The most important part is the `alloc` field. That's where we define the initially deployed smart contracts:

| Name                | Address                                      | Purpose                                                                  |
|---------------------|----------------------------------------------|------------------------------------------------------------------------|
| ValidatorMessages   | `0xca11ab1e00000000000000000000000000000000` | Library contract used by PoAValidatorManager                           |
| PoAValidatorManager | `0xC0DEBA5E00000000000000000000000000000000` | Main validator management contract. References ValidatorMessages       |
| ProxyAdmin          | `0xC0FFEE1234567890aBcDEF1234567890AbCdEf34` | Admin contract for managing the transparent proxy                      |
| TransparentProxy    | `0xFEEDC0DE00000000000000000000000000000000` | Proxy contract that delegates calls to the PoAValidatorManager         |

You may want to take a look at the manual linking process in the `loadDeployedHexFromJSON` function.

### 6. ⛓️  Creating chain

Source code: [06_create_chain/chain.go](./06_create_chain/chain.go)

A chain belongs to a subnet. A node validating a subnet will validate all chains belonging to that subnet. Please note that you also submit the genesis here. Later on the nodes we will create, will pull this genesis from the P-chain.

```golang
createChainTx, err := pWallet.IssueCreateChainTx(
    subnetID,               // Transaction id from 2 steps ago
    genesisBytes,           // L1 genesis
    constants.SubnetEVMID,  // Could be any cb58 string, but for EVM you should use this one
    nil,                    //
    "My L1",                // Just a string
)
```

### 7. 🚀 Launching nodes

Source code: [07_launch_nodes/launch.sh](./07_launch_nodes/launch.sh) with argument `node0`

Setup Steps:
- Writes [07_launch_nodes/evm_debug_config.json](./07_launch_nodes/evm_debug_config.json) to `data/chains/[chainID]/config.json` to enable EVM debugging. This will be used later with the `--chain-config-dir` flag in avalanchego.
- Uses `CURRENT_UID` and `CURRENT_GID` to prevent write access permission issues. This is not specific to avalanchego.
- `TRACK_SUBNETS` loads the current subnet ID so the node can track the subnet and all chains belonging to it.

Launches [07_launch_nodes/docker-compose.yml](./07_launch_nodes/docker-compose.yml). It contains only one node for simplicity. Mounts local `./data/` folder as `/data/`.

### 8. 🔮 Converting chain

Source code: [08_convert_chain/convert.go](./08_convert_chain/convert.go)

This converts your Chain to the new Avalanche L1, introduced at Etna upgrade. 

```golang
tx, err := wallet.P().IssueConvertSubnetToL1Tx(
		subnetID, // Transaction hash from the "Create Subnet" step
		chainID, // Transaction hash from the "Create Subnet" step
		managerAddress.Bytes(), // The address of your manager contract. We added this in genesis
		avaGoBootstrapValidators, // The initial list of validators. Keep it, we will need it in initialization 
		options...,
	)
```

`avaGoBootstrapValidators` is formed using HTTP requests to nodes, like this one:

```bash
curl -X POST --data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"info.getNodeID"
}' -H 'content-type:application/json' 127.0.0.1:9650/ext/info
```

### 9. 🔃 Restarting the node

Source code: none

Runs step 7 again so nodes can pick up changes after the conversion:

```bash
./07_launch_nodes/launch.sh "node0"
```


### 10. 🎯 Activate ProposerVM fork

Source code: [10_activate_proposer_vm/proposer.go](./10_activate_proposer_vm/proposer.go)

Sends test transactions to activate the ProposerVM fork.

- FIXME: Add more details about ProposerVM fork

### 11. 🔌 Initialize PoA validator manager contract

Source code: [11_validator_manager_initialize/initialize.go](./11_validator_manager_initialize/initialize.go)

Here using [bindings for PoAValidatorManager](https://github.com/ava-labs/icm-contracts/blob/main/abi-bindings/go/validator-manager/PoAValidatorManager/PoAValidatorManager.go) to initialize the contract.

```go
contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
if err != nil {
    log.Fatalf("failed to deploy contract: %s\n", err)
}

tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
    L1ID:                   subnetID,
    ChurnPeriodSeconds:     0,
    MaximumChurnPercentage: 20,
}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
```

### 12. 👥 Initialize validator set

Source code: [12_initialize_validator_set/initialize_validator_set.go](./12_initialize_validator_set/initialize_validator_set.go)

Here we call `initializeValidatorSet` method with the same data we used in the "Converting chain". That is very important - if one byte differs, the hash would be different.

```golang
tx, _, err := contract.TxToMethodWithWarpMessage(
		fmt.Sprintf("http://127.0.0.1:9650/ext/bc/%s/rpc", chainID),
		strings.TrimSpace(privateKey),
		managerAddress,
		subnetConversionSignedMessage,
		big.NewInt(0),
		"initialize validator set",
		validatormanager.ErrorSignatureToError,
		"initializeValidatorSet((bytes32,bytes32,address,[(bytes,bytes,uint64)]),uint32)",
		subnetConversionDataPayload,
		uint32(0),
	)
```
### 13. 📄 Check validator set

Source code: [13_check_validator_set/check_validator_set.go](./13_check_validator_set/check_validator_set.go)

Just checks the state of the validator set on the P-chain and the logs in the smart contract:
```golang
validatorsPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "platform.getValidatorsAt",
		"params": map[string]string{
			"height":   "proposed",
			"subnetID": subnetID.String(),
		},
		"id": 1,
	}
// Get validators
validatorsResp, err := makeJSONRPCRequest(client, fujiPChainURL, validatorsPayload)
if err != nil {
  return fmt.Errorf("failed to get validators: %w", err)
}
//...
query := ethereum.FilterQuery{
  Addresses: []common.Address{managerAddress},
}

logs, err := ethClient.FilterLogs(context.Background(), (interfaces.FilterQuery)(query))
if err != nil {
  log.Fatal(err)
}
```

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
