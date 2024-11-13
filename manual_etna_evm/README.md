For this prototype, we will use Avalanche Snap to sign transactions, or even a browser-generated private key with the hope that the Core team will integrate the functionality needed later. We can supply a signer interface and its implementation using Snap, Private Key, and, in the end, Core wallet.

1. Check C-Chain balance and show a link to a faucet
2. Transfer C-Chain funds to P-Chain using Core wallet
3. Call a create subnet transaction (P-Chain)
4. Construct basic genesis by replacing a couple of fields
5. Create subnet transaction (P-Chain)
6. Generate Node BLS keys
7. Convert chain, using those BLS keys
8. PoAValidatorManagerInitialize on EVM
9. GetPChainSubnetConversionWarpMessage
10. InitializeValidatorsSet

Questions:
- There are 2 transactions, 2 WARP messages, what else? (Asked Sarp)
