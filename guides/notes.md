### Set RPC
export LOCAL_RPC=http://127.0.0.1:9650/ext/bc/6MiZgo2yGyetKKn8pPKaDo6hBjvLpS4N7FyJDRwHHNPSuVkbG/rpc

### deploy erc20

forge create contracts/mocks/ExampleERC20.sol:ExampleERC20 --private-key 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 --rpc-url=$LOCAL_RPC




## Main

### deploy native token staker

(I get server returned an error response: error code -32000: max code size exceeded for ERC20StakingManager) 

forge create contracts/validator-manager/NativeTokenStakingManager.sol:NativeTokenStakingManager --constructor-args 0 --private-key 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 --rpc-url=$LOCAL_RPC

```
Deployer: 0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC
Deployed to: 0x4Ac1d98D9cEF99EC6546dEd4Bd550b0b287aaD6D
Transaction hash: 0x130aa259436614d193cabacbe2adb115c8b0f83e6e744293ab57f5c3aa3fa292
```

### upgrade proxy

cast send 0xFEEDBEEF0000000000000000000000000000000A "upgrade(address,address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0x4Ac1d98D9cEF99EC6546dEd4Bd550b0b287aaD6D --private-key 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 --rpc-url=$LOCAL_RPC

### 
```
avalanche contract initPosManager poa
```
```
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
The error is ok.
Under the hood we called `"initialize(((bytes32,uint64,uint8),uint256,uint256,uint64,uint16,uint8,uint256,address))"` then fetched the warp message and signed it.
The part that failed is `"initializeValidatorSet((bytes32,bytes32,address,[(bytes,bytes,uint64)]),uint32)"` which is expected, as the validator set was already initialized on PoA manager.




## Helpers

### read sslot for proxy impl address
cast storage 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc --rpc-url=$LOCAL_RPC

### read sslot for proxy admin on proxy
cast storage 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103 --rpc-url=$LOCAL_RPC

### proxy admin ownership check
cast call 0xFEEDBEEF0000000000000000000000000000000A "owner()" --rpc-url=$LOCAL_RPC

## transfer ownership of proxyadmin
cast send 0xFEEDBEEF0000000000000000000000000000000A "transferOwnership(address)" 0x8db97c7cece249c2b98bdc0226cc4c2a57bf52fc --private-key 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 --rpc-url=$LOCAL_RPC




### get proxy impl from admin

cast call 0xFEEDBEEF0000000000000000000000000000000A "getProxyImplementation(address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 --rpc-url=$LOCAL_RPC 

### get proxy admin from proxy admin

cast call 0xFEEDBEEF0000000000000000000000000000000A "getProxyAdmin(address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 --rpc-url=$LOCAL_RPC

### change proxy admin address (new contract)
cast send 0xFEEDBEEF0000000000000000000000000000000A "changeProxyAdmin(address,address)" 0xC0FFEE1234567890aBcDEF1234567890AbCdEf34 0xFEEDBEEF0000000000000000000000000000000A --private-key 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 --rpc-url=$LOCAL_RPC