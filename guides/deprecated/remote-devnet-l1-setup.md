# Creating a Remote DevNet L1

:::note
This guide has been deprecated as there is no longer any need to create run an independent Avalanche Network devnet to test the features that are currently supported by the Etna Devnet.
:::

This guide covers how to create your own DevNet and L1 chain connected to it, running in AWS.

It uses the [Avalanche CLI](https://github.com/ava-labs/avalanche-cli) to provision the AWS infrastructure and configure the nodes properly.

## Steps

1. Authenticate with the AWS account that you will use to create the infrastructure, and note the profile name to be used in `~/.aws/config`.
2. Clone the [`avalanche-cli`](https://github.com/ava-labs/avalanche-cli) repo, and checkout the ACP-77 branch.
3. Build the CLI by running `./scripts/build.sh`. For ease of use, you can move the `bin/avalanche` binary into a directory in your `PATH`, such as `/usr/local/bin`.
4. Run `avalanche node create <NODE_NAME> --use-ssh-agent --aws --aws-profile <AWS_PROFILE_NAME> --use-static-ip=false` to create your own personal DevNet running AvalancheGo.
   - This by default expects you to use a Yubikey to access the provisioned nodes.
   - Select `Devnet`
   - Use the custom AvalancheGo version of `v1.12.0-initial-poc.5`
   - 1 validator and 1 API node is sufficient for most use cases
5. Get the validators node ID and BLS key information
   - `ssh -A ubuntu@<VALIDATOR_IP_FROM_STEP_4>`
   - Call:
   ```bash
   curl -X POST --data '{
       "jsonrpc":"2.0",
       "id"     :1,
       "method" :"info.getNodeID"
   }' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info
   ```
6. Run `avalanche blockchain create <L1_NAME>` to create your L1 configuration to be created on the DevNet.
7. Deploy the new L1 on the DevNet by running `avalanche blockchain deploy <L1_NAME>`
   - Select `Devnet`
   - Use the cluster name that created in step 4
   - Select Yes to having set up your own nodes
   - Enter the node ID and BLS information gathered from the validator node in step 5
   - Use the `ewoq` key when asked which key should be used to pay for the P-Chain transactions

Once the devnet and node are no longer needed, destroy them using `avalanche node destroy <NODE_NAME>`.
