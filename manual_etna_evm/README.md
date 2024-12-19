## Guide to Etna L1s for Developers

This repository provides a detailed, code-first guide to understanding L1 subnet management on Avalanche Fuji after the Etna upgrade.

### 🎯 Purpose

- **Educational**: Shows step-by-step how L1 subnet management works under the hood
- **Code-First**: Each concept is demonstrated with actual Go code
- **Manual Approach**: Deliberately avoids abstractions to clearly show the underlying processes

### ⚠️ Not for Production

This code intentionally does everything manually for educational purposes. For production use:

- **Developers**: Use [avalanche-cli SDK](https://github.com/ava-labs/avalanche-cli/tree/main/sdk) for programmatic integration
- **Users**: Use [avalanche-cli](https://github.com/ava-labs/avalanche-cli) for command-line management
- **This Guide**: Learn the concepts before implementing your solution

### 🔍 What's Inside

This guide walks through:
- Creating and managing L1 a PoA or PoS subnet
- Launching validators in Docker
- Adding and removing validators PoA (not PoS)

**TL;DR:**
```bash
# Create a PoA L1
./create.sh

# Print validators (may need ~30s for the P-Chain transaction to propagate)
go run . validators

# Print validator manager contract logs from node0
go run . logs 9650

# Add a validator
go run . add-poa-validator

# A docker run command will be displayed. Copy and paste it into your terminal, 
# and wait about 5 minutes for the new node to bootstrap.

go run . logs 9652 
# Use the AVALANCHEGO_HTTP_PORT from the docker run command instead of 9652
# This won't work until the node is fully bootstrapped, which takes about 5 minutes

# Remove a validator
go run . remove-poa-validator
# This will fail initially but print a list of nodes.
# Copy the NodeID of the node with weight 20, then run:

go run . remove-poa-validator NodeID-xxxxxxxxxxxxxxxxxx
# This removes the second node from the validator set, although it still runs.
```

**Requirements:**
- A fresh Docker installation (verify by running `docker compose ls` without any dashes).
- Go 1.22.10+.

Run `./create.sh` to create a new L1 on Devnet. Use `./cleanup.sh` to clean up afterward (this preserves your keys).

Use `go run . validators` to print the current validators.

Use `go run . logs 9650` to print contract logs from node0, and `go run . logs 9652` for node1, etc.

Below is an updated programming guide that follows the original style, maintaining code references, highlighting key conceptual steps, and including representative code snippets for each phase. With the updated file structure, we now reference `cmd/` directories.

> **Note:** These examples are simplified, linear demonstrations with hardcoded values, not intended for production use.

> **Note:** Creating both PoA and PoS networks is supported, but adding and removing validators is currently only supported on PoA.
---

### 1. 🔑 Generating Keys

**Source code:** [cmd/01_01_generate_keys.go](cmd/01_01_generate_keys.go)

Generates the validator manager private key if you don't have one yet.

Main method used:  
`secp256k1.NewPrivateKey()` from `github.com/ava-labs/avalanchego/utils/crypto/secp256k1`.

---

### 2. 💰 Transfer AVAX between C and P chains

**Source code:** [cmd/01_02_transfer_coins_C_to_P.go](cmd/01_02_transfer_coins_C_to_P.go)

- Checks your C-chain and P-chain balances.
- If insufficient funds on P-chain, exports from C-chain to P-chain.
- If still short, directs you to the [Fuji faucet](https://test.core.app/tools/testnet-faucet/?subnet=c&token=c).

Representative snippet:
```go
exportTx, err := cWallet.IssueExportTx(
    constants.PlatformChainID,
    []*secp256k1fx.TransferOutput{{
        Amt:          cChainBalance.Uint64() - 100*units.MilliAvax,
        OutputOwners: owner,
    }},
)
importTx, err := pWallet.IssueImportTx(cWallet.Builder().Context().BlockchainID, &owner)
```

> **Note:** We recommend transferring at least 4 AVAX to your P-chain address to avoid issues with the minimum balance. 

---

### 3. 🕸️ Creating subnet

**Source code:** [cmd/01_03_create_subnet.go](cmd/01_03_create_subnet.go)

A subnet is created, owned by a given address:

```go
owner := &secp256k1fx.OutputOwners{
    Locktime:  0,
    Threshold: 1,
    Addrs:     []ids.ShortID{subnetOwner},
}

createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
```

---

### 4. 📝 Generating Genesis

**Source code:** [cmd/01_04_generate_genesis.go](cmd/01_04_generate_genesis.go)

Generates a custom L1 genesis to include in the P-chain `createChain` transaction.

We define initial alloc and deploy initial contracts:

| Name                | Address                                      | Purpose                               |
|---------------------|----------------------------------------------|---------------------------------------|
| ProxyAdmin          | `0xC0FFEE1234567890aBcDEF1234567890AbCdEf34` | Admin contract for the proxy           |
| TransparentProxy    | `0xFEEDC0DE00000000000000000000000000000000` | Proxy contract delegating calls        |

We are using a precompiled Transparent proxy contract from OpenZeppelin, version 4.9. 

---

### 5. ⛓️ Creating chain

**Source code:** [cmd/01_05_create_chain.go](cmd/01_05_create_chain.go)

Creates a new chain for your subnet with the previously generated L1 genesis.

```go
createChainTx, err := pWallet.IssueCreateChainTx(
    subnetID,
    genesisBytes,
    constants.SubnetEVMID,
    nil,
    "My L1",
)
```

---

### 6. 🔮 Converting chain to Avalanche L1

**Source code:** [cmd/01_06_convert_to_L1.go](cmd/01_06_convert_to_L1.go)

Converts your chain into an Avalanche L1.  
You must provide bootstrap validators. Manager contract address is generated as the second deployed contract from the validator manager owner key.

```go
tx, err := wallet.P().IssueConvertSubnetToL1Tx(
    subnetID,
    chainID,
    managerAddress.Bytes(),
    avaGoBootstrapValidators,
    options...,
)
```

`avaGoBootstrapValidators` comes from local staker and signer credentials or via RPC `info.getNodeID`.

---

### 7. 🚀 Launching a validator node

**Source code:** [cmd/01_07_launch_node.go](cmd/01_07_launch_node.go)

Launches a local node using `docker-compose` and the configured environment. Enables EVM debugging and tracks the newly created subnet:

```bash
docker compose up -d --build
```

The node image: `containerman17/avalanchego-subnetevm:v1.12.0_v0.7.0`  
Contains a precompiled SubnetEVM and canonical container configuration options.

---

### 8. 🧱 Deploying the Validator Manager contract

**Source code:** [cmd/01_08_deploy_validator_manager.go](cmd/01_08_deploy_validator_manager.go)

Deploys your chosen validator manager contract (PoA or PoSNative) on the L1:

```go
newContractAddress, tx, _, err := poavalidatormanager.DeployPoAValidatorManager(opts, ethClient, 0)
```

Check the pre-defined addresses in `config.go` for Proxy, Admin, etc.

---

### 9. 🔌 Initializing the Validator Manager contract

**Source code:** [cmd/01_09_validator_manager_init.go](cmd/01_09_validator_manager_init.go)

```go
tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
    L1ID:                   subnetID,
    ChurnPeriodSeconds:     0,
    MaximumChurnPercentage: 20,
}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
```

---

### 10. 👥 Print validators (check P-chain state)

**Source code:** [cmd/01_10_print_validators.go](cmd/01_10_print_validators.go)

Queries the P-chain for current validators:

```go
validatorsPayload := map[string]interface{}{
    "jsonrpc": "2.0",
    "method":  "platform.getValidatorsAt",
    ...
}
validatorsResp, err := makeJSONRPCRequest(client, fujiPChainURL, validatorsPayload)
```

Prints them along with subnet info.

---

### 11. 📄 Print contract logs

**Source code:** [cmd/01_11_print_contract_logs.go](cmd/01_11_print_contract_logs.go)

Fetches and parses on-chain logs from the manager contract:

```go
logs, err := ethClient.FilterLogs(context.Background(), query)
...
if event, err := contract.ParseInitialized(vLog); err == nil {
    fmt.Printf("Initialized event...\n")
}
```

---

### 12. 🔮 Initialize validator set

**Source code:** [cmd/01_12_initialize_validator_set.go](cmd/01_12_initialize_validator_set.go)

Once the chain is converted to L1 and the manager is set up, initialize the validator set with warp messages:

```go
tx, _, err := contract.TxToMethodWithWarpMessage(
    evmChainURL,
    privateKey,
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

---

### Add PoA Validator to an existing L1

#### Step A1: 👾 Initialize registration

**Source code:** [cmd/02_01_add_validator_poa_step_1.go](cmd/02_01_add_validator_poa_step_1.go)

Collect node ID, BLS keys, and generate a warp message to initialize a new validator’s registration:

```go
tx, _, err := contract.TxToMethod(
    rpcURL,
    managerOwnerPrivateKey,
    managerAddress,
    big.NewInt(0),
    "initialize validator registration",
    validatorManagerSDK.ErrorSignatureToError,
    "initializeValidatorRegistration((bytes,bytes,uint64,(uint32,[address]),(uint32,[address])),uint64)",
    validatorRegistrationInput,
    weight,
)
```

#### Step A2: 📝 Register on P-chain

**Source code:** [cmd/02_02_add_validator_poa_step_2.go](cmd/02_02_add_validator_poa_step_2.go)

Use the signed warp message to register the validator on the P-chain:

```go
unsignedTx, err := wallet.P().Builder().NewRegisterL1ValidatorTx(
    1*units.Avax,
    proofOfPossession.ProofOfPossession,
    warpMessage.Bytes(),
)
wallet.P().IssueTx(&tx)
```

#### Step A3: 🏰 Complete registration

**Source code:** [cmd/02_03_add_validator_poa_step_3.go](cmd/02_03_add_validator_poa_step_3.go)

Finalize validator registration by calling `completeValidatorRegistration` with the warp message:

```go
tx, _, err := contract.TxToMethodWithWarpMessage(
    rpcURL,
    privateKey,
    managerAddress,
    subnetValidatorRegistrationSignedMessage,
    big.NewInt(0),
    "complete validator registration",
    validatorManagerSDK.ErrorSignatureToError,
    "completeValidatorRegistration(uint32)",
    uint32(0),
)
```

#### Step A4: 🚀 Start another node

**Source code:** [cmd/02_04_add_validator_poa_step_4.go](cmd/02_04_add_validator_poa_step_4.go)

Prints a `docker run` command to start another node with proper credentials and environment variables.

---

### Remove PoA Validator from existing L1

#### Step R1: Initialize removal

**Source code:** [cmd/03_05_remove_validator_step_1.go](cmd/03_05_remove_validator_step_1.go)

Initialize the validator removal process, sign a warp message to set their weight to zero:

```go
tx, _, err := contract.TxToMethod(
    nodeURL,
    hex.EncodeToString(privateKey.Bytes()),
    managerAddress,
    big.NewInt(0),
    "initializeEndValidation",
    validatormanager.ErrorSignatureToError,
    "initializeEndValidation(bytes32)",
    validationID,
)
```

#### Step R2: Adjust validator weight on P-chain

**Source code:** [cmd/03_05_remove_validator_step_2.go](cmd/03_05_remove_validator_step_2.go)

Issue a `setL1ValidatorWeight` transaction with the warp message to update weight on P-chain:

```go
unsignedTx, err := wallet.P().Builder().NewSetL1ValidatorWeightTx(message.Bytes())
wallet.P().IssueTx(&tx)
```

#### Step R3: Complete removal

**Source code:** [cmd/03_05_remove_validator_step_3.go](cmd/03_05_remove_validator_step_3.go)

Complete the removal by calling `completeEndValidation` with the finalized warp message:

```go
tx, _, err := contract.TxToMethodWithWarpMessage(
    rpcURL,
    privateKey,
    managerAddress,
    signedMessage,
    big.NewInt(0),
    "complete poa validator removal",
    validatorManagerSDK.ErrorSignatureToError,
    "completeEndValidation(uint32)",
    uint32(0),
)
```
