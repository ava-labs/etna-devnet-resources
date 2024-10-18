# How to Deploy a Sovereign L1 on the Etna Devnet

Use the CLI to create, deploy, and convert your L1 tracked by a locally run Node.

Warning: this flow is in active development. None of the following should be used in or with production-related infrastructure.

In this guide, we will be creating a sovereign L1 with locally run Avalanche Nodes as its bootstrap validators.
At the end of this guide, we will also go through adding and removing validators in our sovereign L1

## Build Etna-enabled AvalancheGo

```zsh
mkdir -p $GOPATH/src/github.com/ava-labs
cd $GOPATH/src/github.com/ava-labs
git clone https://github.com/ava-labs/avalanchego.git
cd $GOPATH/src/github.com/ava-labs/avalanchego
git checkout v1.12.0-initial-poc.5
./scripts/build.sh
```

Take note of path of AvalancheGo build as we will use it later on.

## Build Etna-enabled Avalanche CLI

In a separate terminal window:

```zsh
git clone https://github.com/ava-labs/avalanche-cli.git
cd avalanche-cli
git checkout acp-77
./scripts/build.sh
```

## Create Blockchain

You can use the following command to create the blockchain:

```zsh
./bin/avalanche blockchain create <chainName> --evm --proof-of-authority
```

Select `I want to use defaults for a production environment`
Choose your configs, and for ease of use, just use the `ewoq` key for everything.

```zsh
✔ Get address from an existing stored key (created from avalanche key create or avalanche key import)
✔ ewoq
✓ Validator Manager Contract owner address 0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC
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

If it's a public network and you're using an ewoq key, you will receive the following error. This is for security reasons, to avoid attacks.

```zsh
Error: can't airdrop to default address on public networks, please edit the genesis by calling `avalanche subnet create <chainName> --force`
```

When you try to generate the blockchain again with --force flag, new keys named `subnet_<chainName>_airdrop` will be generated for you, which you can use.

```zsh
✔ Get address from an existing stored key (created from avalanche key create or avalanche key import)
Use the arrow keys to navigate: ↓ ↑ → ← 
? Which stored key should be used enable as controller of ValidatorManager contract?: 
    ewoq
  ▸ subnet_<chainName>_airdrop
    cli-awm-relayer
    cli-teleporter-deployer
```

When the blockchain deploy command is called, it will:
- Create 5 Avalanche Nodes on your local machine. 
- Add these 5 nodes as bootstrap validators in your sovereign L1
- Have these nodes track your L1
- Initialize Validator Manager Contract on your L1

By the end of your command, you would have a running sovereign L1 with Proof of Authority Validator Manager
Contract deployed into it!

## Destroy Nodes

To tear down your local Avalanche Nodes, run:

```zsh
./bin/avalanche node local destroy <nodeClusterName>
```

`nodeClusterName` is in the form of <chainName>-local-node

## Restart Nodes

To restart your local Avalanche nodes after a shutdown, run:

```zsh
./bin/avalanche node local start <nodeClusterName> --etna-devnet --avalanchego-path=<avalancheGoBuildPath>
```

## Adding a new validator

We will first create a new Avalanche Node to be added as a new validator. Since we already have
a local AvalancheGo process running, we will create a new node in AWS / GCP.

```zsh
`./bin/avalanche node create <newClusterName> --custom-avalanchego-version=v1.12.0-initial-poc.5 --etna-devnet`
```

More info regarding `avalanchenode create` command can be found at [our docs](https://docs.avax.network/tooling/create-avalanche-nodes/run-validators-aws).

Next we will ssh into the created node to get the Node ID and BLS info.

To SSH into our node, run:

```zsh
`avalanche node ssh <newClusterName>`
```

run the printed out command to ssh into the remote node

Next, to get the Node ID and BLS info, run:

```zsh
`curl -X POST --data '{
"jsonrpc":"2.0",
"id"     :1,
"method" :"info.getNodeID"
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info`
```

Now we will be adding this node as a validator in our sovereign L1. Note that `nodeClusterName` is the
in the form of <chainName>-local-node if it was created automatically through `avalanche subnet deploy`
command.

```zsh
`./bin/avalanche blockchain addValidator <chainName> --cluster <nodeClusterName>`
```

Enter the Node ID and BLS info we obtained

## Removing validator

To remove a validator, run: 

```zsh
`./bin/avalanche blockchain removeValidator <chainName> --cluster <nodeClusterName>`
```
