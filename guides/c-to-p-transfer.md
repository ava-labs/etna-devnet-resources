# Send AVAX from C to P-Chain

If you want to transfer AVAX between X,P, and C-Chain you can use the examples from [avalanchejs](https://github.com/meaghanfitzgerald/avalanchejs) to send import/export txns.

## How To

1. Create a set of keys. You can do this with the [CLI](https://docs.avax.network/tooling/avalanche-cli#key-create).

2. Get the private and public addresses made with this key with the CLI command:

```zsh
avalanche key list --devnet --endpoint https://etna.avax-dev.network --keys <KEYNAME>
```

3. Get C-Chain Devnet AVAX from the [Faucet](https://core.app/tools/testnet-faucet/?subnet=cdevnet&token=cdevnet).

4. Clone this fork of the [Avalanchejs](https://github.com/meaghanfitzgerald/avalanchejs) repo.

5. Rename `example.env` file to `.env`, and populate with your keys. All non-EVM addresses should be prepended with `P-custom`. Save the file.

6.Run `yarn install` in your terminal.

7.Run `yarn run export-c` in your terminal.

9. Run `yarn run import-p` in your terminal.

10. Check your balance with the CLI.

```zsh
avalanche key list --devnet --endpoint https://etna.avax-dev.network --keys <KEYNAME>
```

If you require more AVAX than the faucet provides, reach out to our team on the relevant Telegram or Slack channel.
