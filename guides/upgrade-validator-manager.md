# Upgrading ValidatorManager


## Context
Both the PoA and PoS options offered in the CLI `acp-77-pos` branch are deployed behind a `TransparentProxy` contract included in genesis at the address:
```bash
0xC0FFEE1234567890aBcDEF1234567890AbCdEf34
```

The implementation for either PoA or PoS ValidatorManager (depending on what you choose in CLI) is included in genesis at the address:
```bash
0x5F584C2D56B4c356e7d82EC6129349393dc5df17
```

The admin of this `TransparentProxy` is a `ProxyAdmin` contract included in genesis at the address:
```bash
0xFEEDBEEF0000000000000000000000000000000A
```

The owner of the `ProxyAdmin` contract is decided when you select a `ValidatorManager` owner during
```
avalanche blockchain create
```

This address has control over upgrading the `ValidatorManager` implementation referenced by the `TransparentProxy`.

### More Info
The `TransparentProxy` and `ProxyAdmin` are from [OpenZeppelin v4.9](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/release-v4.9/contracts/proxy/transparent/TransparentUpgradeableProxy.sol)



## Set environment variables
Set RPC
```bash
export LOCAL_RPC=http://127.0.0.1:9650/ext/bc/6MiZgo2yGyetKKn8pPKaDo6hBjvLpS4N7FyJDRwHHNPSuVkbG/rpc
```
Set PKEY (this is key to ewoq in CLI)
```bash
export PKEY=56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027
```


## Upgrading PoA -> PoS

### 1. Deploy `NativeTokenStakingManager`

```bash
forge create contracts/validator-manager/NativeTokenStakingManager.sol:NativeTokenStakingManager --constructor-args 0 --private-key $PKEY --rpc-url=$LOCAL_RPC
```

```bash
Deployer: 0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC
Deployed to: 0x4Ac1d98D9cEF99EC6546dEd4Bd550b0b287aaD6D
Transaction hash: 0x130aa259436614d193cabacbe2adb115c8b0f83e6e744293ab57f5c3aa3fa292
```

Point the `TransparentProxy`'s implementation to the new address through a call to `ProxyAdmin`

```bash
cast send 0xFEEDBEEF0000000000000000000000000000000A "upgrade(address,address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0x4Ac1d98D9cEF99EC6546dEd4Bd550b0b287aaD6D --private-key $PKEY --rpc-url=$LOCAL_RPC
```

### 2. Initialize the new NativeTokenStakingManager

```bash
avalanche contract initPosManager poa
```

```bash
✔ Devnet
✔ Get Devnet RPC endpoint from an existing node cluster (created from avalanche node create or avalanche devnet wiz)
✔ poa-local-node
RPC Endpoint: http://127.0.0.1:9650/ext/bc/6MiZgo2yGyetKKn8pPKaDo6hBjvLpS4N7FyJDRwHHNPSuVkbG/rpc
✔ Enter the minimum stake amount: 1█
Enter the minimum stake amount: 1
Enter the maximum stake amount: 100
Enter the minimum stake duration (in seconds): 10
Enter the minimum stake duration (in seconds): 10
Enter the minimum delegation fee: 100
Enter the maximum stake multiplier: 2
Enter the weight to value factor: 1
Error: failure initializing validators set on pos manager: validators set already initialized (txHash=0x7ff2860be4db6c6d760653864e96fdbcd4699b2f4d3fcfd783d8bc95d6a4d5c1)
```
**The error is ok.**
Under the hood of the CLI we called `"initialize(((bytes32,uint64,uint8),uint256,uint256,uint64,uint16,uint8,uint256,address))"` then fetched the warp message and signed it + submitted it.

The part that failed is `"initializeValidatorSet((bytes32,bytes32,address,[(bytes,bytes,uint64)]),uint32)"` which is expected, as the validator set was already initialized on PoA manager.


## More Info

For more info on how to perform upgrades please reference `ValidatorManager` contracts:
https://github.com/ava-labs/teleporter/tree/main/contracts/validator-manager#convert-poa-to-pos

It is important to note this PoS upgrade won't be completely functional in distributing rewards as there is no `RewardCalculator` precompile on default settings for PoA network.

To include the `RewardCalculator` either make sure it is enabled in genesis, or perform a network upgrade:

reference: https://docs.avax.network/virtual-machines/evm-customization/precompile-overview

### Helpers

Check the current Proxy implementation through sslot
```bash
cast storage 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc --rpc-url=$LOCAL_RPC
```

Check the current Proxy admin through sslot
```bash
cast storage 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103 --rpc-url=$LOCAL_RPC
```

Check Proxy's implementation through ProxyAdmin
```bash
cast call 0xFEEDBEEF0000000000000000000000000000000A "getProxyAdmin(address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 --rpc-url=$LOCAL_RPC
```

Check Proxy's admin through ProxyAdmin

```bash
cast call 0xFEEDBEEF0000000000000000000000000000000A "getProxyAdmin(address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 --rpc-url=$LOCAL_RPC
```

Check owner of ProxyAdmin
```bash
cast call 0xFEEDBEEF0000000000000000000000000000000A "owner()" --rpc-url=$LOCAL_RPC
```

