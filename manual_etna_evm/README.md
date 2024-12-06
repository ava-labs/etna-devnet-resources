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
cd /teleporter_src/contracts && forge build --extra-output-files bin
```

### 5. 🧱 Generating genesis

Source code: [05_L1_genesis/genesis.go](./05_L1_genesis/genesis.go)

Here we generate the Genesis for our new L1. We will include it in a P-chain create chain transaction in the next step.

Note: Don't confuse your L1 genesis with the Avalanche Fuji genesis.

Read more about genesis here: [https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1](https://docs.avax.network/avalanche-l1s/upgrade/customize-avalanche-l1).


The most important part is `alloc` field. That's where we define the initially deployed smart contracts:

| Name                | Address                                      | Purpose                                                                  |
|---------------------|----------------------------------------------|------------------------------------------------------------------------|
| ValidatorMessages   | `0xca11ab1e00000000000000000000000000000000` | Library contract used by PoAValidatorManager                           |
| PoAValidatorManager | `0xC0DEBA5E00000000000000000000000000000000` | Main validator management contract. References ValidatorMessages       |
| ProxyAdmin          | `0xC0FFEE1234567890aBcDEF1234567890AbCdEf34` | Admin contract for managing the transparent proxy                      |
| TransparentProxy    | `0xFEEDC0DE00000000000000000000000000000000` | Proxy contract that delegates calls to the PoAValidatorManager         |


You probably would like to take a look into the manual linking process in `loadDeployedHexFromJSON` function.

<!--
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
<!-- 
### 7. 🛠️ Compile the Validator Manager Contract

After the Etna upgrade, L1s are managed by Warp messages emitted by L1. Currently, the most functional implementation is the [Validator Manager Contract](https://github.com/ava-labs/icm-contracts/tree/790ccce873f9a904910a0f3ffd783436c920ce97/contracts/validator-manager) in the [Teleporter Repo](https://github.com/ava-labs/icm-contracts).

In this step, we first install the [ava-labs/foundry fork](https://github.com/ava-labs/foundry):

```dockerfile
RUN curl -o install_foundry.sh https://raw.githubusercontent.com/ava-labs/teleporter/${ICM_COMMIT}/scripts/install_foundry.sh && \
    chmod +x install_foundry.sh && \
    ./install_foundry.sh && \
    rm install_foundry.sh
```

Then, download and compile the teleporter repository:
```bash
git clone https://github.com/ava-labs/icm-contracts /teleporter
# ....
cd /teleporter/contracts && forge build --extra-output-files=bin
```

The compiled json would be copied to [07_compile_validator_manager/PoAValidatorManager.sol/PoAValidatorManager.json](./07_compile_validator_manager/PoAValidatorManager.sol/PoAValidatorManager.json). 

### 7. 📦 Deploy the Validator Manager Contract

Source code: [07_depoly_validator_manager/deploy.go](./07_depoly_validator_manager/deploy.go)

Using the pre-compiled PoA Validator Manager from [abi-bindings/go/validator-manager/PoAValidatorManager/PoAValidatorManager.go](https://github.com/ava-labs/icm-contracts/blob/main/abi-bindings/go/validator-manager/PoAValidatorManager/PoAValidatorManager.go) in the ava-labs/teleporter repo, deploy it using standard EVM Go bindings.

> In production, you should put PoAValidatorManager behind a transparent proxy and preferably include it in genesis.

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

### 9. 🔃 Restarting nodes

Source code: none

Runs step 6 again so nodes can pick up changes after the upgrade
<!-- 
### 9. 🏥 Checking subnet health

Source code: [09_check_subnet_health/health.go](./09_check_subnet_health/health.go)

Polls `http://127.0.0.1:9650/ext/bc/[CHAIN_ID]/rpc` and requests the EVM chainID until it receives a response. The endpoint becomes available once the node is fully booted and synced, which can take a few minutes. You can monitor progress with `docker logs -f node0`.

FIXME: [Health API](https://docs.avax.network/api-reference/health-api) is a better option.

### 10. 💸 Sending some test coins

Source code: [10_evm_transfer/transfer.go](./10_evm_transfer/transfer.go)

Sends a test transfer using the generic EVM API. This double checks that the chain is operational. 

### 10. 🎯 Activate ProposerVM fork

Source code: [10_activate_proposer_vm/proposer.go](./10_activate_proposer_vm/proposer.go)

Sends test transactions to activate the ProposerVM fork.

- FIXME: Add more details about ProposerVM fork
- FIXME: Investigate if this can be combined with EVM transfers to eliminate this step

### 11. Initialize PoA validator manager contract

Source code: [11_validator_manager_initialize/initialize.go](./11_validator_manager_initialize/initialize.go)

TODO: Describe this step

### 13. Initialize validator set

Source code: [12_initialize_validator_set/initialize_validator_set.go](./12_initialize_validator_set/initialize_validator_set.go)

TODO: Describe this step

### 14. Add 2 more validators

TODO: Implementation pending

### 15. Remove a validator

TODO: Implementation pending


**Example logs:**

```bash

$ time ./run.sh 

🔑 Generating keys

2024/12/05 09:42:07 POA validator manager key already exists in ./data/ folder

💰 Checking balance

2024/12/05 09:42:11 fetched state in 2.076972955s
2024/12/05 09:42:11 P-chain balance: 16.897684872 AVAX
2024/12/05 09:42:11 ✅ P-chain balance sufficient

🕸️  Creating subnet

2024/12/05 09:42:14 Synced wallet in 1.953804954s
2024/12/05 09:42:17 ✅ Created new subnet 2UMiGz9yn7SiPDY175V5r3MPeh4TDH5RMVS7AJ6bBAPqXLz6Qe in 2.467001749s

🛠️ Using hardcoded smart contracts code


🧱 Generating genesis

2024/12/05 09:42:19 ✅ Successfully wrote genesis to data/L1-genesis.json

⛓️  Creating chain

2024/12/05 09:42:20 Using vmID: srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
2024/12/05 09:42:22 synced wallet in 2.23792055s
2024/12/05 09:42:25 ✅ Created new chain MCBThPucMkznGTsZwhgdJ5MbcXzGMY7zKRFcNu3GC2D8mbLHJ in 2.621411955s
2024/12/05 09:42:25 ✅ Saved chain ID to file

🚀 Launching nodes

WARN[0000] Warning: No resource found to remove for project "07_launch_nodes". 
[+] Building 1.2s (21/21) FINISHED                                                                                                                           docker:default
 => [node0 internal] load build definition from Dockerfile                                                                                                             0.0s
 => => transferring dockerfile: 1.34kB                                                                                                                                 0.0s
 => [node0 internal] load metadata for docker.io/library/debian:bookworm-slim                                                                                          1.1s
 => [node0 internal] load metadata for docker.io/library/golang:1.22-bookworm                                                                                          1.1s
 => [node0 internal] load .dockerignore                                                                                                                                0.0s
 => => transferring context: 2B                                                                                                                                        0.0s
 => [node0 stage-2 1/9] FROM docker.io/library/debian:bookworm-slim@sha256:1537a6a1cbc4b4fd401da800ee9480207e7dc1f23560c21259f681db56768f63                            0.0s
 => [node0 internal] load build context                                                                                                                                0.0s
 => => transferring context: 35B                                                                                                                                       0.0s
 => [node0 subnet-evm-builder 1/3] FROM docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                 0.0s
 => => resolve docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                                          0.0s
 => CACHED [node0 stage-2 2/9] RUN apt-get update                                                                                                                      0.0s
 => CACHED [node0 stage-2 3/9] RUN apt-get install -y wget                                                                                                             0.0s
 => CACHED [node0 stage-2 4/9] RUN groupadd -r nobody || true                                                                                                          0.0s
 => CACHED [node0 avalanchego-builder 2/3] WORKDIR /app                                                                                                                0.0s
 => CACHED [node0 avalanchego-builder 3/3] RUN git clone https://github.com/ava-labs/avalanchego.git && cd avalanchego && git checkout v1.12.0 && ./scripts/build.sh   0.0s
 => CACHED [node0 stage-2 5/9] COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego                                           0.0s
 => CACHED [node0 subnet-evm-builder 2/3] RUN git clone https://github.com/ava-labs/subnet-evm.git /app/subnet-evm && cd /app/subnet-evm && git checkout v0.6.12       0.0s
 => CACHED [node0 subnet-evm-builder 3/3] RUN cd /app/subnet-evm && go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin                     0.0s
 => CACHED [node0 stage-2 6/9] COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7F  0.0s
 => CACHED [node0 stage-2 7/9] RUN wget -O /fuji-latest.tar https://avalanchego-public-database.avax-test.network/p-chain/avalanchego/data-tar/latest.tar              0.0s
 => CACHED [node0 stage-2 8/9] COPY entrypoint.sh /entrypoint.sh                                                                                                       0.0s
 => CACHED [node0 stage-2 9/9] RUN chmod +x /entrypoint.sh                                                                                                             0.0s
 => [node0] exporting to image                                                                                                                                         0.0s
 => => exporting layers                                                                                                                                                0.0s
 => => writing image sha256:f00a821550ec2a11996e6163b8d275b3964c3ad6edf55bf8030a7faf22c12c29                                                                           0.0s
 => => naming to docker.io/library/07_launch_nodes-node0                                                                                                               0.0s
 => [node0] resolving provenance for metadata file                                                                                                                     0.0s
[+] Running 1/1
 ✔ Container node0  Started                                                                                                                                            0.1s 
Waiting for subnet to become available...
🌱 Subnet is still starting up (attempt 1 of 100)
./fuji/v1.4.5/030245.ldb
🌱 Subnet is still starting up (attempt 2 of 100)
[12-05|09:42:34.473] INFO network/ip_tracker.go:542 reset validator tracker bloom filter {"currentCount": 184}
✅ Subnet is healthy and responding
Chain ID (decimal): 12345
To see logs, run: docker logs -f node0

🔮 Converting chain into L1

Using changeOwnerAddress: P-fuji1rn8whlwk3f53yua6wly82hhdn2ms6f9py5ss4q
2024/12/05 09:42:51 Getting node info from http://127.0.0.1:9650
2024/12/05 09:42:51 Issuing convert subnet tx
subnetID: 2UMiGz9yn7SiPDY175V5r3MPeh4TDH5RMVS7AJ6bBAPqXLz6Qe
chainID: MCBThPucMkznGTsZwhgdJ5MbcXzGMY7zKRFcNu3GC2D8mbLHJ
managerAddress: 0feedc0de0000000000000000000000000000000
avaGoBootstrapValidators[0]:
        NodeID: 49d526a32e09dde119dd84179ce9234cecb27c86
        BLS Public Key: 8ae72249dace07b3e9bbc886da5e6b3c5f2df3af459673f2535b4d054cf57c14dbfe88fe3035b4863d966b419aaadc2f
        Weight: 100
        Balance: 1000000000

✅ Convert subnet tx ID: 2KX22D1QXKUpn47fkTDNpGhnHRHNnEGoBZgV4PmPduqQWgPXLn

🚀 Restarting nodes

[+] Running 1/1
 ✔ Container node0  Removed                                                                                                                                           10.1s 
[+] Building 0.6s (21/21) FINISHED                                                                                                                           docker:default
 => [node0 internal] load build definition from Dockerfile                                                                                                             0.0s
 => => transferring dockerfile: 1.34kB                                                                                                                                 0.0s
 => [node0 internal] load metadata for docker.io/library/debian:bookworm-slim                                                                                          0.5s
 => [node0 internal] load metadata for docker.io/library/golang:1.22-bookworm                                                                                          0.5s
 => [node0 internal] load .dockerignore                                                                                                                                0.0s
 => => transferring context: 2B                                                                                                                                        0.0s
 => [node0 stage-2 1/9] FROM docker.io/library/debian:bookworm-slim@sha256:1537a6a1cbc4b4fd401da800ee9480207e7dc1f23560c21259f681db56768f63                            0.0s
 => [node0 internal] load build context                                                                                                                                0.0s
 => => transferring context: 35B                                                                                                                                       0.0s
 => [node0 subnet-evm-builder 1/3] FROM docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                 0.0s
 => => resolve docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                                          0.0s
 => CACHED [node0 stage-2 2/9] RUN apt-get update                                                                                                                      0.0s
 => CACHED [node0 stage-2 3/9] RUN apt-get install -y wget                                                                                                             0.0s
 => CACHED [node0 stage-2 4/9] RUN groupadd -r nobody || true                                                                                                          0.0s
 => CACHED [node0 avalanchego-builder 2/3] WORKDIR /app                                                                                                                0.0s
 => CACHED [node0 avalanchego-builder 3/3] RUN git clone https://github.com/ava-labs/avalanchego.git && cd avalanchego && git checkout v1.12.0 && ./scripts/build.sh   0.0s
 => CACHED [node0 stage-2 5/9] COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego                                           0.0s
 => CACHED [node0 subnet-evm-builder 2/3] RUN git clone https://github.com/ava-labs/subnet-evm.git /app/subnet-evm && cd /app/subnet-evm && git checkout v0.6.12       0.0s
 => CACHED [node0 subnet-evm-builder 3/3] RUN cd /app/subnet-evm && go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin                     0.0s
 => CACHED [node0 stage-2 6/9] COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7F  0.0s
 => CACHED [node0 stage-2 7/9] RUN wget -O /fuji-latest.tar https://avalanchego-public-database.avax-test.network/p-chain/avalanchego/data-tar/latest.tar              0.0s
 => CACHED [node0 stage-2 8/9] COPY entrypoint.sh /entrypoint.sh                                                                                                       0.0s
 => CACHED [node0 stage-2 9/9] RUN chmod +x /entrypoint.sh                                                                                                             0.0s
 => [node0] exporting to image                                                                                                                                         0.0s
 => => exporting layers                                                                                                                                                0.0s
 => => writing image sha256:f00a821550ec2a11996e6163b8d275b3964c3ad6edf55bf8030a7faf22c12c29                                                                           0.0s
 => => naming to docker.io/library/07_launch_nodes-node0                                                                                                               0.0s
 => [node0] resolving provenance for metadata file                                                                                                                     0.0s
[+] Running 1/1
 ✔ Container node0  Started                                                                                                                                            0.1s 
Waiting for subnet to become available...
🌱 Subnet is still starting up (attempt 1 of 100)
🌱 Subnet is still starting up (attempt 2 of 100)
[12-05|09:43:08.904] INFO network/ip_tracker.go:542 reset validator tracker bloom filter {"currentCount": 184}
✅ Subnet is healthy and responding
Chain ID (decimal): 12345
To see logs, run: docker logs -f node0

🎯 Activate ProposerVM fork

Initial block height: 0
Block height after activation: 2
2024/12/05 09:43:26 ✅ Successfully activated proposer VM fork

🔌 Initialize Validator Manager

✅ Transaction receipt: {
  "type": "0x2",
  "root": "0x",
  "status": "0x1",
  "cumulativeGasUsed": "0x1e3b0",
  "logsBloom": "0x00000000000000080000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000008000000000000000020000000000000000000800000000000000000000000000000000400000000000000000000800400200000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000400000000000004000000000000000020000000000000000000000000000000000000000000000000000000000000000000",
  "logs": [
    {
      "address": "0x0feedc0de0000000000000000000000000000000",
      "topics": [
        "0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0",
        "0x0000000000000000000000000000000000000000000000000000000000000000",
        "0x00000000000000000000000073c07d5e006e99323075e6a7b53d94c27db24c08"
      ],
      "data": "0x",
      "blockNumber": "0x3",
      "transactionHash": "0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24",
      "transactionIndex": "0x0",
      "blockHash": "0xcb7838c5c76ccc8741901d741d74f833a3c0faac3efa9bfe3d18c79904966506",
      "logIndex": "0x0",
      "removed": false
    },
    {
      "address": "0x0feedc0de0000000000000000000000000000000",
      "topics": [
        "0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2"
      ],
      "data": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "blockNumber": "0x3",
      "transactionHash": "0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24",
      "transactionIndex": "0x0",
      "blockHash": "0xcb7838c5c76ccc8741901d741d74f833a3c0faac3efa9bfe3d18c79904966506",
      "logIndex": "0x1",
      "removed": false
    }
  ],
  "transactionHash": "0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24",
  "contractAddress": "0x0000000000000000000000000000000000000000",
  "gasUsed": "0x1e3b0",
  "effectiveGasPrice": "0x5d21dba01",
  "blockHash": "0xcb7838c5c76ccc8741901d741d74f833a3c0faac3efa9bfe3d18c79904966506",
  "blockNumber": "0x3",
  "transactionIndex": "0x0"
}
✅ Validator manager initialized at: 0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24
2024/12/05 09:43:28 ✅ Validator manager initialized

👥 Initialize validator set

2024/12/05 09:43:31 Getting node info from http://127.0.0.1:9650
{"level":"info","timestamp":"2024-12-05T09:43:31.224Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:348","msg":"Signing subnet not found, requesting from PChain","blockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"warn","timestamp":"2024-12-05T09:43:31.839Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:269","msg":"Failed to make async request to node","nodeID":"NodeID-7jPfzQLKZHgBf6Jp3PHxm1Mz4vFcN1sCG"}
{"level":"info","timestamp":"2024-12-05T09:43:33.842Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:311","msg":"Created signed message.","warpMessageID":"2PNM5tYVHA9mJg8fxZk9wjU93U7CrrjgGj66frXuVhcxFyzPZh","signatureWeight":100,"sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
❌ Failed to initialize validator set: failed to initialize validator set: invalid warp message
Attempt 2/3 (will sleep for 10s before retry)
2024/12/05 09:43:44 Getting node info from http://127.0.0.1:9650
{"level":"info","timestamp":"2024-12-05T09:43:44.894Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:348","msg":"Signing subnet not found, requesting from PChain","blockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"warn","timestamp":"2024-12-05T09:43:45.367Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:269","msg":"Failed to make async request to node","nodeID":"NodeID-7jPfzQLKZHgBf6Jp3PHxm1Mz4vFcN1sCG"}
{"level":"info","timestamp":"2024-12-05T09:43:47.369Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:311","msg":"Created signed message.","warpMessageID":"2PNM5tYVHA9mJg8fxZk9wjU93U7CrrjgGj66frXuVhcxFyzPZh","signatureWeight":100,"sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
✅ Successfully initialized validator set. Transaction hash: 0xaca00b00c02ada5bec42fa66e257f47fef87f1442fea900434b6d0714f71f81a

📄 Reading contract logs

P-Chain State:
------------------------

Validators:
NodeID: NodeID-7jPfzQLKZHgBf6Jp3PHxm1Mz4vFcN1sCG
  Public Key: 0x8ae72249dace07b3e9bbc886da5e6b3c5f2df3af459673f2535b4d054cf57c14dbfe88fe3035b4863d966b419aaadc2f
  Weight: 100

Subnet Info:
Is Permissioned: false
Control Keys: [P-fuji1rn8whlwk3f53yua6wly82hhdn2ms6f9py5ss4q]
Threshold: 1
Manager Chain ID: MCBThPucMkznGTsZwhgdJ5MbcXzGMY7zKRFcNu3GC2D8mbLHJ
Manager Address: 0x0feedc0de0000000000000000000000000000000




------------------------
Log TxHash: 0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24
OwnershipTransferred:
  Previous Owner: 0x0000000000000000000000000000000000000000
  New Owner: 0x73c07D5e006E99323075E6A7B53D94C27dB24C08
------------------------
Log TxHash: 0x8a0f99ce23664f8e64df15dc84131339e56149468ae8fcca5a18a5f2bcd97f24
Initialized:
  Version: 1
------------------------
Log TxHash: 0xaca00b00c02ada5bec42fa66e257f47fef87f1442fea900434b6d0714f71f81a
InitialValidatorCreated:
  ValidationID: 06d0314a8a40890b906a641b026b818a8d5bdfba1cbe1cd72585127f0635cd00
  NodeID: f8a84f640ef9a7db16a2e807934d94d13280b0350c98e883980a62d19794bd15
  Weight: 100

🚀 Starting 1 more node

[+] Running 1/1
 ✔ Container node0  Removed                                                                                                                                           10.1s 
[+] Building 0.7s (26/28)                                                                                                                                    docker:default
 => [node1 internal] load build definition from Dockerfile                                                                                                             0.0s
 => => transferring dockerfile: 1.34kB                                                                                                                                 0.0s
 => [node0 internal] load build definition from Dockerfile                                                                                                             0.0s
 => => transferring dockerfile: 1.34kB                                                                                                                                 0.0s
 => [node1 internal] load metadata for docker.io/library/debian:bookworm-slim                                                                                          0.5s
 => [node1 internal] load metadata for docker.io/library/golang:1.22-bookworm                                                                                          0.6s
 => [node1 internal] load .dockerignore                                                                                                                                0.0s
 => => transferring context: 2B                                                                                                                                        0.0s
 => [node0 internal] load .dockerignore                                                                                                                                0.0s
 => => transferring context: 2B                                                                                                                                        0.0s
 => [node0 stage-2 1/9] FROM docker.io/library/debian:bookworm-slim@sha256:1537a6a1cbc4b4fd401da800ee9480207e7dc1f23560c21259f681db56768f63                            0.0s
 => [node1 subnet-evm-builder 1/3] FROM docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                 0.0s
 => => resolve docker.io/library/golang:1.22-bookworm@sha256:0d22c0d84536a5bb9bdd5b65b71fad5df32e648b2dfd10cb3fd87e4063da0f9c                                          0.0s
 => [node1 internal] load build context                                                                                                                                0.0s
 => => transferring context: 35B                                                                                                                                       0.0s
 => [node0 internal] load build context                                                                                                                                0.0s
 => => transferring context: 35B                                                                                                                                       0.0s
 => CACHED [node0 stage-2 2/9] RUN apt-get update                                                                                                                      0.0s
 => CACHED [node0 stage-2 3/9] RUN apt-get install -y wget                                                                                                             0.0s
 => CACHED [node0 stage-2 4/9] RUN groupadd -r nobody || true                                                                                                          0.0s
 => CACHED [node0 avalanchego-builder 2/3] WORKDIR /app                                                                                                                0.0s
 => CACHED [node0 avalanchego-builder 3/3] RUN git clone https://github.com/ava-labs/avalanchego.git && cd avalanchego && git checkout v1.12.0 && ./scripts/build.sh   0.0s
 => CACHED [node0 stage-2 5/9] COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego                                           0.0s
 => CACHED [node0 subnet-evm-builder 2/3] RUN git clone https://github.com/ava-labs/subnet-evm.git /app/subnet-evm && cd /app/subnet-evm && git checkout v0.6.12       0.0s
 => CACHED [node0 subnet-evm-builder 3/3] RUN cd /app/subnet-evm && go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin                     0.0s
 => CACHED [node0 stage-2 6/9] COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7F  0.0s
 => CACHED [node0 stage-2 7/9] RUN wget -O /fuji-latest.tar https://avalanchego-public-database.avax-test.network/p-chain/avalanchego/data-tar/latest.tar              0.0s
 => CACHED [node0 stage-2 8/9] COPY entrypoint.sh /entrypoint.sh                                                                                                       0.0s
 => CACHED [node1 stage-2 9/9] RUN chmod +x /entrypoint.sh                                                                                                             0.0s
 => [node0] exporting to image                                                                                                                                         0.0s
 => => exporting layers                                                                                                                                                0.0s
 => => writing image sha256:f00a821550ec2a11996e6163b8d275b3964c3ad6edf55bf8030a7faf22c12c29                                                                           0.0s
 => => naming to docker.io/library/07_launch_nodes-node0                                                                                                               0.0s
 => [node1] exporting to image                                                                                                                                         0.0s
 => => exporting layers                                                                                                                                                0.0s
 => => writing image sha256:5ed1034f2076dc40e88f41808c7ea35d0247e8b06eeae25c3dffd92bd00379a7                                                                           0.0s
 => => naming to docker.io/library/07_launch_nodes-node1                                                                                                               0.0s
 => [node0] resolving provenance for metadata file                                                                                                                     0.0s
 => [node1] resolving provenance for metadata file                                                                                                                     0.0s
[+] Running 2/2
 ✔ Container node0  Started                                                                                                                                            0.1s 
 ✔ Container node1  Started                                                                                                                                            0.1s 
Waiting for subnet to become available...
🌱 Subnet is still starting up (attempt 1 of 100)
🌱 Subnet is still starting up (attempt 2 of 100)
[12-05|09:44:05.973] INFO network/ip_tracker.go:542 reset validator tracker bloom filter {"currentCount": 184}
✅ Subnet is healthy and responding
Chain ID (decimal): 12345
To see logs, run: docker logs -f node0

👥 Add validator - initialize registration

2024/12/05 09:44:22 Adding validator on port 9652
2024/12/05 09:44:22 Getting node info from http://127.0.0.1:9652
2024/12/05 09:44:23 ✅ Validator registration initialized: 0x70a91c82ebffbe86a41d8e0616626af0e69d7d240ccb01b4f79693e20a7b99c5
{"level":"info","timestamp":"2024-12-05T09:44:23.988Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:348","msg":"Signing subnet not found, requesting from PChain","blockchainID":"MCBThPucMkznGTsZwhgdJ5MbcXzGMY7zKRFcNu3GC2D8mbLHJ"}
{"level":"warn","timestamp":"2024-12-05T09:44:24.624Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:269","msg":"Failed to make async request to node","nodeID":"NodeID-7jPfzQLKZHgBf6Jp3PHxm1Mz4vFcN1sCG"}
{"level":"info","timestamp":"2024-12-05T09:44:26.626Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:311","msg":"Created signed message.","warpMessageID":"2V29MEUnoc6eVa4sd3UjbAoEW3Mz56jKpQPAdMYGvSQpg4Nk4Q","signatureWeight":100,"sourceBlockchainID":"MCBThPucMkznGTsZwhgdJ5MbcXzGMY7zKRFcNu3GC2D8mbLHJ"}
validationID: sr7dRbzYtxepJ79Emr6eoqExurGYSHVg5VUm1J288sGJo9GYh

👥 Add validator - register on P-chain

2024/12/05 09:44:28 Attempting to register L1 validator on P-chain...
2024/12/05 09:44:28 Attempt 1/3
2024/12/05 09:44:33 ✅ Successfully registered L1 validator on P-chain

👥 Add validator - complete validator registration

{"level":"info","timestamp":"2024-12-05T09:44:35.368Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:348","msg":"Signing subnet not found, requesting from PChain","blockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"warn","timestamp":"2024-12-05T09:44:36.001Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:269","msg":"Failed to make async request to node","nodeID":"NodeID-7jPfzQLKZHgBf6Jp3PHxm1Mz4vFcN1sCG"}
{"level":"info","timestamp":"2024-12-05T09:44:44.019Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:311","msg":"Created signed message.","warpMessageID":"2TaLHcFrFndM7SP1TfWZjrFD4RsPrfRBBeGWG4M1bjfvaahRcb","signatureWeight":100,"sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
\p�/�]w�.5^�����ϗK���m���|�߂�GZ+9�*弹?����%/���U(��dQ��O�����M�`�3~����\6{�5��~���

real    2m38.504s
user    0m26.618s
sys     0m8.638s
```
