# How to Add and Remove Validators from a PoA L1

In this guide, we will walk through how to add and remove validators from a Proof-of-Authority L1. If you have not yet created an L1, you must first complete the steps outlined in [this guide](/guides/deploy-sovereign-l1.md) before proceeding.

## Adding a New L1 Validator

We will first create a new Avalanche Node to be added as a new validator. Since we already have
a local AvalancheGo process running, we will create a new node in AWS / GCP.

```zsh
`./bin/avalanche node create <newClusterName> --custom-avalanchego-version=v1.12.0-initial-poc.6 --etna-devnet`
```

More info regarding `avalanche node create` command can be found in [our docs](https://docs.avax.network/tooling/create-avalanche-nodes/run-validators-aws).

Next, we will ssh into the created node to get the Node ID and BLS info.

To SSH into our node, run:

```zsh
`avalanche node ssh <newClusterName>`
```

The CLI will print a command that enables you to ssh into the remote node.

Next, to get the Node ID and BLS info, run:

```zsh
`curl -X POST --data '{
"jsonrpc":"2.0",
"id"     :1,
"method" :"info.getNodeID"
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info`
```

Note that `nodeClusterName` is in the form of <chainName>-local-node if it was created automatically with the `avalanche subnet deploy` command.

```zsh
`./bin/avalanche blockchain addValidator <chainName> --cluster <nodeClusterName>`
```

When prompted, enter the Node ID and BLS info we obtained earlier.

## Removing an L1 Validator

To remove an L1 validator from a PoA chain, just run:

```zsh
`./bin/avalanche blockchain removeValidator <chainName> --cluster <nodeClusterName>`
```
