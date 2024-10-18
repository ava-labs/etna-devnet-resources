# How to Deploy a Sovereign L1 on the Etna Devnet using the Avalanche Starter Kit

The Avalanche Starter Kit contains everything you need to get started quickly with Avalanche. Among other tools it contains Avalanche CLI. With that you can set up a local network, create a Avalanche L1, customize the Avalanche L1/VM configuration, and so on.

## Open the Avalanche Starter Kit Github Repository:

Open [Avalanche Starter Kit](https://github.com/ava-labs/avalanche-starter-kit/tree/acp-77), and ensure you are on the correct branch, `acp-77`.

## Create a Codespace

- Click the green Code button
- Switch to the Codespaces tab
- Click Create Codespace on main or if you already have a Codespace click the plus (+) button
- The Codespace will open in a new tab. Wait a few minutes until it's fully built.

All required packages, such as `avalanchego` and `avalanche-cli`, will come pre-installed with their correct versions.

## Create Blockchain

You can use the following command to create the blockchain:

```zsh
avalanche blockchain create <chainName> --evm --proof-of-authority
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
File /home/vscode/.avalanche-cli/subnets/stt2/chain.json successfully written
✓ Successfully created blockchain configuration
Run 'avalanche blockchain describe' to view all created addresses and what their roles are
```

## Deploy

You can deploy the blockchain and boot validator nodes using the following command, referencing the `avalanchego` location:

```zsh
avalanche blockchain deploy <chainName> --etna-devnet --use-local-machine --avalanchego-path=/usr/local/bin/avalanchego
```

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

## Make RPC-Endpoint publicly accessible

Since the Avalanche Network is running in a Github Codespace the localhost (127.0.0.1) will only be accessible from inside the Codespace.

Therefore, we need to make the RPC-Endpoint publicly accessible. Click on the little antenna icon in the bottom bar of the Codespace:

![](https://avalanche-academy-git-starter-kit-pause-and-resume-ava-labs.vercel.app/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Fports-open.ea371fea.png&w=828&q=75)

## Destroy Nodes

To tear down your local Avalanche Nodes, run:

```zsh
avalanche node local destroy <nodeClusterName>
```

`nodeClusterName` is in the form of <chainName>-local-node

## Restart Nodes

To restart your local Avalanche nodes after a shutdown, run:

```zsh
avalanche node local start <nodeClusterName> --etna-devnet --avalanchego-path=<avalancheGoBuildPath>
```