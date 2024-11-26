## Guide to Etna L1s for golang devs

This repo provides a detailed, code-first guide for integrating L1 subnet management into your services on Avalanche Fuji after the Etna upgrade. For end-user subnet management, check out [avalanche-cli](https://github.com/ava-labs/avalanche-cli).

** Requirements: **
- Fresh Docker (check by running `docker compose ls` without any dashes)
- Go 1.22.8+

Run all together: `./run.sh` to start a new L1 on Devnet, `./cleanup.sh` to clean up (keeps your keys)

### 1. 🔑 Generating keys

Source code: [01_generate_keys/generate.go](./01_generate_keys/generate.go)

Generates a validator manager private key, if you don't have one yet.

What you are looking for is method `secp256k1.NewPrivateKey()` from package `github.com/ava-labs/avalanchego/utils/crypto/secp256k1`.

### 2. 💰 Checking balance

Source code: [02_check_balance/balance.go](./02_check_balance/balance.go)

- Checks your P-chain balance.
- Tries to export funds from C chain to P Chain, if less than 1.1 AVAX (which is required to create a subnet)
- If not enough funds, recommends you to [visit Fuji faucet](https://test.core.app/tools/testnet-faucet/?subnet=c&token=c). 

A good example of checking balances and moving AVAX between C and P chains.

### 3. 🕸️ Creating subnet

Source code: [03_create_subnet/create.go](./03_create_subnet/create.go)

Creating subnet requires only an owner. Here is the gist of it:

```golang
owner := &secp256k1fx.OutputOwners{
    Locktime:  0,
    Threshold: 1,
    Addrs:     []ids.ShortID{subnetOwner},
}

createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
```

Subnet is a group of validators agreed to validate the same chains. One chain can be only in one subnet, but one subnet could have multiple chains, and each validator has to validate all chains of a subnet. At the same time a validator can be a part of 2 different subnets. 


### 4. 🧱 Generating genesis

Source code: [04_L1_genesis/genesis.go](./04_L1_genesis/genesis.go)

Here we generate Genesis fo our new L1. We will put it in a P chain create chain transaction in the next step. 

Do not confuse your L1 genesis and Avalanche Fuji genesis. For your node you'll need both of them.

It does fill a lot of diferent params, but the most important function calls are `validatormanager.AddPoAValidatorManagerContractToAllocations` and `validatormanager.AddTransparentProxyContractToAllocations` from `github.com/ava-labs/avalanche-cli/pkg/validatormanager` package.

### 5. ⛓️  Creating chain

Source code: [05_create_chain/chain.go](./05_create_chain/chain.go)

```golang
createChainTx, err := pWallet.IssueCreateChainTx(
    subnetID,               // Transaction id from 2 steps ago
    genesisBytes,           // L1 genesis
    constants.SubnetEVMID,  // Really could be any cb58 sting, but for EVM you should use 
    nil,                    // TODO: figure out what fixture extension is
    "My L1",                // Just a string
)
```

### 6. 🚀 Launching nodes

Source code: [06_launch_nodes/launch.sh](./06_launch_nodes/launch.sh)

Preparation: 
- Writes [06_launch_nodes/evm_debug_config.json](./06_launch_nodes/evm_debug_config.json) into `data/chains/[chainID]/config.json` to enable debug in EVM. Will be used later with flag `--chain-config-dir` in avalanchego. 
- `CURRENT_UID` and `CURRENT_GID` are used to avoid write access right problems. Not avalanchego specific.
- `TRACK_SUBNETS` loads current subnetID, so the node could track the subnet and all chains that belong to this subnet

Launches [06_launch_nodes/docker-compose.yml](./06_launch_nodes/docker-compose.yml). It contains only one node for simplicity. Mounts local `./data/` folder as `/data/`

### 7. 🔄 Converting chain

Source code: [07_convert_chain/convert.go](./07_convert_chain/convert.go)

TODO: Describe

### 8. 🔃 Restarting nodes

Source code: none

Runs step 6 again so nodes could pick up changes after upgrade

### 9. 🏥 Checking subnet health

Source code: [9_check_subnet_health/health.go](./9_check_subnet_health/health.go)

Knocks on `http://127.0.0.1:6550/ext/bc/[CHAIN_ID]/rpc` and asks for an EVM chainID intill gets an answer. Once the node is fully booted and synced, it will be available. It might take a couple menites. Monitor with `docker logs -f node0`. 

### 10. 💸 Sending some test coins

Source code: [9_check_subnet_health/health.go](9_check_subnet_health/health.go)

Sending a test transfer using generic EVM API. Just double checks that's the chain is operational. 
