## Guide to Etna L1s for Go Developers

This repository provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

**Requirements:**
- Fresh Docker installation (verify by running `docker compose ls` without any dashes)
- Go 1.22.8+

Run everything at once: `./run.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (preserves your keys)

### 1. 🔑 Generating Keys

Source code: [01_generate_keys/generate.go](./01_generate_keys/generate.go)

Generates a validator manager private key if you don't have one yet.

The key method you'll need is `secp256k1.NewPrivateKey()` from package `github.com/ava-labs/avalanchego/utils/crypto/secp256k1`.

### 2. 💰 Checking balance

Source code: [02_check_balance/balance.go](./02_check_balance/balance.go)

- Checks your P-chain balance
- Attempts to export funds from C-chain to P-chain if balance is less than 1.1 AVAX (required for subnet creation)
- If insufficient funds, directs you to the [Fuji faucet](https://test.core.app/tools/testnet-faucet/?subnet=c&token=c)

This provides a good example of checking balances and transferring AVAX between C and P chains.

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

### 4. 🧱 Generating genesis

Source code: [04_L1_genesis/genesis.go](./04_L1_genesis/genesis.go)

Here we generate the Genesis for our new L1. We will include it in a P-chain create chain transaction in the next step.

Note: Don't confuse your L1 genesis with the Avalanche Fuji genesis. Your node will need both.

~~The most important function calls are `validatormanager.AddPoAValidatorManagerContractToAllocations` and `validatormanager.AddTransparentProxyContractToAllocations` from the `github.com/ava-labs/avalanche-cli/pkg/validatormanager` package.~~

Read more about genesis here: [https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1](https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1).

> Normally, you would want to include the Transparent Proxy and Validator manager contracts in genesis, but in this tutorial, for the purpose of a more granular workflow, we are going to deploy them manually in the later steps.

### 5. ⛓️  Creating chain

Source code: [05_create_chain/chain.go](./05_create_chain/chain.go)

```golang
createChainTx, err := pWallet.IssueCreateChainTx(
    subnetID,               // Transaction id from 2 steps ago
    genesisBytes,           // L1 genesis
    constants.SubnetEVMID,  // Could be any cb58 string, but for EVM you should use this one
    nil,                    // FIXME: Document fixture extension usage
    "My L1",                // Just a string
)
```

### 6. 🚀 Launching nodes

Source code: [06_launch_nodes/launch.sh](./06_launch_nodes/launch.sh)

Setup Steps:
- Writes [06_launch_nodes/evm_debug_config.json](./06_launch_nodes/evm_debug_config.json) to `data/chains/[chainID]/config.json` to enable EVM debugging. This will be used later with the `--chain-config-dir` flag in avalanchego.
- Uses `CURRENT_UID` and `CURRENT_GID` to prevent write access permission issues. This is not specific to avalanchego.
- `TRACK_SUBNETS` loads the current subnet ID so the node can track the subnet and all chains belonging to it.

Launches [06_launch_nodes/docker-compose.yml](./06_launch_nodes/docker-compose.yml). It contains only one node for simplicity. Mounts local `./data/` folder as `/data/`

### 7. 🔮 Converting chain

Source code: [07_convert_chain/convert.go](./07_convert_chain/convert.go)

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

### 8. 🔃 Restarting nodes

Source code: none

Runs step 6 again so nodes can pick up changes after the upgrade

### 9. 🏥 Checking subnet health

Source code: [09_check_subnet_health/health.go](./09_check_subnet_health/health.go)

Polls `http://127.0.0.1:6550/ext/bc/[CHAIN_ID]/rpc` and requests the EVM chainID until it receives a response. The endpoint becomes available once the node is fully booted and synced, which can take a few minutes. You can monitor progress with `docker logs -f node0`.

FIXME: [Health API](https://docs.avax.network/api-reference/health-api) is a better option.

### 10. 💸 Sending some test coins

Source code: [10_evm_transfer/transfer.go](./10_evm_transfer/transfer.go)

Sends a test transfer using the generic EVM API. This double checks that the chain is operational.

### 11. 🎯 Activate ProposerVM fork

Sends test transactions to activate the ProposerVM fork.

- FIXME: Add more details about ProposerVM fork
- FIXME: Investigate if this can be combined with EVM transfers to eliminate this step

### 12. Initialize PoA validator manager contract

TODO: Implementation pending

### 13. Initialize validator set

TODO: Implementation pending

### 14. Add 2 more validators

TODO: Implementation pending

### 15. Remove a validator

TODO: Implementation pending
