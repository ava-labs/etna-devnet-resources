# How to Deploy a PoS L1 on the Etna Devnet

Use the CLI to create, deploy, and convert your L1 tracked by a locally run Node.

Warning: this flow is in active development. None of the following should be used in or with production-related infrastructure.

In this guide, we will be creating a sovereign L1 with locally run Avalanche Nodes as its bootstrap validators with the Native Token PoS ValidatorManager.

## Build Etna-enabled AvalancheGo

```zsh
mkdir -p $GOPATH/src/github.com/ava-labs
cd $GOPATH/src/github.com/ava-labs
git clone https://github.com/ava-labs/avalanchego.git
cd $GOPATH/src/github.com/ava-labs/avalanchego
git checkout v1.11.13
./scripts/build.sh
```

Take note of path of AvalancheGo build as we will use it later on.

## Build Etna-enabled Avalanche CLI

In a separate terminal window:

```zsh
git clone https://github.com/ava-labs/avalanche-cli.git
cd avalanche-cli
git checkout acp-77-pos
./scripts/build.sh
```

## Create Blockchain

You can use the following command to create the blockchain:

```zsh
./bin/avalanche blockchain create <chainName> --evm --proof-of-stake
```

```
Enter reward basis points for PoS Reward Calculator: 100
```

We include an [ExampleRewardCalculator](https://github.com/ava-labs/teleporter/blob/main/contracts/validator-manager/ExampleRewardCalculator.sol) in the genesis with storage for rewardBasisPoints set to this parameter

You can include the `--reward-basis-points` flag instead to skip this prompt.

Select `I want to use defaults for a production environment`
Choose your configs, and for ease of use, just use the `ewoq` key for everything.

```zsh
✔ I want to use defaults for a production environment
Chain ID: <chainId>
Token Symbol: <symbol>
prefunding address 0x187b4F2412825D8d359308195f0026D4932a3Cf0 with balance 1000000000000000000000000
File /home/vscode/.avalanche-cli/subnets/<chainName>/chain.json successfully written
✓ Successfully created blockchain configuration
Run 'avalanche blockchain describe' to view all created addresses and what their roles are
```

## Deploy

You can deploy the blockchain and boot validator nodes using the following command, referencing the `avalanchego` location:

```zsh
./bin/avalanche blockchain deploy <chainName> --etna-devnet --use-local-machine --avalanchego-path=<avalancheGoBuildPath>
```

If you installed avalanchego with the workflow defined in the [Build Etna-enabled AvalancheGo](#build-etna-enabled-avalanchego) section, or according to [this tutorial](https://docs.avax.network/nodes/run-a-node/manually), the `avalancheGoBuildPath` should be `$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego`.

When the blockchain deploy command is called, it will:

- Create 5 Avalanche Nodes on your local machine.
- Add these 5 nodes as bootstrap validators in your sovereign L1
- Have these nodes track your L1
- Initialize Native Token PoS Validator Manager Contract on your L1

For the prompting of Native Token Staking Manger choose these default values to avoid error:

```
Enter the minimum stake amount (1 = 1 NATIVE TOKEN): 1
Enter the maximum stake amount (1 = 1 NATIVE TOKEN): 1000
Enter the minimum stake duration (in seconds): 10
Enter the minimum delegation fee (in bips): 1
Enter the maximum stake multiplier: 1
Enter the weight to value factor: 1
```

By the end of your command, you should have a running sovereign L1 with a Proof of Stake Native Token Validator Manager
Contract deployed into it!

## Notice

The initial validator set is still treated as a Proof of Authority network, in order for validators to receive rewards they must first be cycled (removed from validator set and added again).

## Add Validator

```bash
./avalanche blockchain addValidator pos
```

Follow the prompting,
select **Yes** when asked if network is PoS

```bash
✔ Enter the amount of tokens to stake (in OWEN): 3
✔ Enter the delegation fee (in bips): 100
✔ Enter the stake duration (in seconds): 100

Validator weight: 3
ValidationID: nJTTNdofYhZttdjH234BiGKHeSpi9L5wtLmyxhkoL4iyaJWJY
RegisterSubnetValidatorTX fee: 0.000001546 AVAX
RegisterSubnetValidatorTx ID: XVsUBdXSE3ixqkwjoDP6dYDqNMrnMLjKGmhKFveTjgmxZVb8T
Waiting for P-Chain to update validator information ... 100% [===============]
  NodeID: NodeID-C3Cbvw4wDtGuWRCJh2v3LttdiLdNWYwox
  Network: Cluster pos-local-node
  Weight: 3
  Balance: 10
✓ Validator successfully added to the Subnet
```

## Destroy Nodes

To tear down your local Avalanche Nodes, run:

```zsh
./bin/avalanche node local destroy <nodeClusterName>
```

`nodeClusterName` is in the form of <chainName>-local-node

or

```zsh
killall avalanche
```

## Restart Nodes

To restart your local Avalanche nodes after a shutdown, run:

```zsh
./bin/avalanche node local start <nodeClusterName> --etna-devnet --avalanchego-path=<avalancheGoBuildPath> --pos
```
