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

## Download Avalanche CLI

In a separate terminal window:

```zsh
curl -sSfL https://raw.githubusercontent.com/ava-labs/avalanche-cli/main/scripts/install.sh | sh -s v1.8.0-rc0
```

Next:

`avalanche blockchain create <chainName> --evm --proof-of-authority`

Select `I want to use defaults for a production environment`
Choose your configs, and for ease of use, just use the `ewoq` key for everything.

Then:

`avalanche blockchain deploy <chainName> --etna-devnet --use-local-machine --avalanchego-path=<avalancheGoBuildPath>`

Use the `ewoq` key for everything.

When the blockchain deploy command is called, it will:
- Create 5 Avalanche Nodes on our local machine. 
- Add these 5 nodes as bootstrap validators in our sovereign L1
- Have these nodes track our L1
- Initialize Validator Manager Contract on our L1

By the end of your command, you would have a running sovereign L1 with Proof of Authority Validator Manager
Contract deployed into it!

To tear down our local Avalanche Nodes, run:

`avalanche node local destroy <nodeClusterName>`

nodeClusterName is in the form of <chainName>-local-node