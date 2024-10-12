# Etna DevNet Resources

The Etna DevNet is a temporary Avalanche network instance that was created for the purpose of testing and integrating with the changes introduced in the Etna upgrade prior to their activation on the Fuji testnet. The network may be wiped or reset at any time as new AvalancheGo versions become ready for testing. The network will be deprecated following the activation of the Etna upgrade on the Fuji testnet.

This README and repository is a collection of various resources to help in interacting and integrating with the DevNet.

## Common Credentials

```text
RPC: https://etna.avax-dev.network/ext/bc/C/rpc
Name: Etna C-Chain
Chain ID: 43117
Token: AVAX
Faucet: https://core.app/tools/testnet-faucet/?subnet=cdevnet&token=cdevnet
Glacier Dev APIs: https://glacier-api-dev.avax.network/api
NetworkID: 76
```

[Devnet Explorer](https://2ffd1590.etna-83w.pages.dev/)

## Public RPCs

The Etna DevNet has a public RPC endpoint available at **`https://etna.avax-dev.network`**.

This endpoint supports all of the common extensions for the primary network chains supported for Fuji and Mainnet. Examples include:

```zsh
curl --location 'https://etna.avax-dev.network/ext/info' \
--header 'Content-Type: application/json' \
--data '{
    "jsonrpc":"2.0",
    "id"     :1,
    "method" :"info.getNetworkID"
}'
```

```zsh
curl --location 'https://etna.avax-dev.network/ext/bc/C/rpc' \
--header 'Content-Type: application/json' \
--data '{
    "jsonrpc": "2.0",
    "method": "eth_chainId",
    "params": [],
    "id": 1
}'
```

You can find the Postman Collection and an Example Environment in the [resources](resources) file.

## Running AvalancheGo connected to Etna DevNet

To run an AvalancheGo node connected to the Etna DevNet:

1. Ensure that you are using an AvalancheGo build from the `v1.12.0-initial-poc.5` tag or later.

```zsh
➜  avalanchego git:(master) git pull
➜  avalanchego git:(master) git checkout v1.12.0-initial-poc.5
➜  avalanchego git:(v1.12.0-initial-poc.5) ./scripts/build.sh
```

For more info on the setup required to run a Node, see this [tutorial](https://docs.avax.network/nodes/run-a-node/manually).

2. Specify the `network-id`, `bootstrap-ids`, `bootstrap-ips`, `genesis-file-content`, and `upgrade-file-content` below. These can also be provided via configuration files if desired.

```zsh
./build/avalanchego \
    --network-id="network-76" \
    --bootstrap-ids="NodeID-8LbTmmGsDC991SbD8Nkx88VULT3XYzYXC,NodeID-bojBKDrpt81bYhxYKQfLw89V7CpoH2m7,NodeID-WrLWMK5sJ4dBUAsx1dP2FUyTqrYwbFA1,NodeID-DDhXtFm6Q9tCq2yiFRmcSMKvHgUgh8yQC,NodeID-QDYnWDQd6g4cQ5H6yiWNqSmfRMBqEH9AG" \
    --bootstrap-ips="52.201.126.172:9651,34.233.248.130:9651,107.21.11.213:9651,35.170.144.5:9651,98.82.41.186:9651" \
    --upgrade-file-content="ewogICAgImFwcmljb3RQaGFzZTFUaW1lIjogIjIwMjAtMTItMDVUMDU6MDA6MDBaIiwKICAgICJhcHJpY290UGhhc2UyVGltZSI6ICIyMDIwLTEyLTA1VDA1OjAwOjAwWiIsCiAgICAiYXByaWNvdFBoYXNlM1RpbWUiOiAiMjAyMC0xMi0wNVQwNTowMDowMFoiLAogICAgImFwcmljb3RQaGFzZTRUaW1lIjogIjIwMjAtMTItMDVUMDU6MDA6MDBaIiwKICAgICJhcHJpY290UGhhc2U0TWluUENoYWluSGVpZ2h0IjogMCwKICAgICJhcHJpY290UGhhc2U1VGltZSI6ICIyMDIwLTEyLTA1VDA1OjAwOjAwWiIsCiAgICAiYXByaWNvdFBoYXNlUHJlNlRpbWUiOiAiMjAyMC0xMi0wNVQwNTowMDowMFoiLAogICAgImFwcmljb3RQaGFzZTZUaW1lIjogIjIwMjAtMTItMDVUMDU6MDA6MDBaIiwKICAgICJhcHJpY290UGhhc2VQb3N0NlRpbWUiOiAiMjAyMC0xMi0wNVQwNTowMDowMFoiLAogICAgImJhbmZmVGltZSI6ICIyMDIwLTEyLTA1VDA1OjAwOjAwWiIsCiAgICAiY29ydGluYVRpbWUiOiAiMjAyMC0xMi0wNVQwNTowMDowMFoiLAogICAgImNvcnRpbmFYQ2hhaW5TdG9wVmVydGV4SUQiOiAiMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTFMcG9ZWSIsCiAgICAiZHVyYW5nb1RpbWUiOiAiMjAyMC0xMi0wNVQwNTowMDowMFoiLAogICAgImV0bmFUaW1lIjogIjIwMjQtMTAtMDlUMjA6MDA6MDBaIgp9Cg==" \
    --genesis-file-content="ewogICJuZXR3b3JrSUQiOiA3NiwKICAiYWxsb2NhdGlvbnMiOiBbCiAgICB7CiAgICAgICJldGhBZGRyIjogIjB4QzcxQTYxYTgxNWU0OWQxNkM0MjU0ODJBMzQyYTM2N0NENDJFMzhhNiIsCiAgICAgICJhdmF4QWRkciI6ICJYLWN1c3RvbTF2NnZ1d3hqZ3IwNDNzZzBudXVocTcwazZ2Z251bGU2OTJmdm5yOSIsCiAgICAgICJpbml0aWFsQW1vdW50IjogNTAwMDAwMDAwMDAwMDAwMDAwLAogICAgICAidW5sb2NrU2NoZWR1bGUiOiBbCiAgICAgICAgewogICAgICAgICAgImFtb3VudCI6IDEwMDAwMDAwMDAwMDAwMDAwMCwKICAgICAgICAgICJsb2NrdGltZSI6IDE2MzM4MjQwMDAKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICJhbW91bnQiOiAxMDAwMDAwMDAwMDAwMDAwMDAsCiAgICAgICAgICAibG9ja3RpbWUiOiAxNjMzODI1MDAwCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAiYW1vdW50IjogMTAwMDAwMDAwMDAwMDAwMDAwLAogICAgICAgICAgImxvY2t0aW1lIjogMTYzMzgyNjAwMAogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgImFtb3VudCI6IDEwMDAwMDAwMDAwMDAwMDAwMCwKICAgICAgICAgICJsb2NrdGltZSI6IDE2MzM4MjcwMDAKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICJhbW91bnQiOiAxMDAwMDAwMDAwMDAwMDAwMDAsCiAgICAgICAgICAibG9ja3RpbWUiOiAxNjMzODI4MDAwCiAgICAgICAgfQogICAgICBdCiAgICB9CiAgXSwKICAic3RhcnRUaW1lIjogMTcyNTMwMDAwMCwKICAiaW5pdGlhbFN0YWtlRHVyYXRpb24iOiAzMTUzMDAwMCwKICAiaW5pdGlhbFN0YWtlRHVyYXRpb25PZmZzZXQiOiA1NDAwLAogICJpbml0aWFsU3Rha2VkRnVuZHMiOiBbCiAgICAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiCiAgXSwKICAiaW5pdGlhbFN0YWtlcnMiOiBbCiAgICB7CiAgICAgICJub2RlSUQiOiAiTm9kZUlELWdwWFdCRXhRU1pYcUpQUXQ2TDZNbnZlVWZncjdISjRxIiwKICAgICAgInJld2FyZEFkZHJlc3MiOiAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiLAogICAgICAiZGVsZWdhdGlvbkZlZSI6IDYyNTAwLAogICAgICAic2lnbmVyIjogewogICAgICAgICJwdWJsaWNLZXkiOiAiMHhhMTRkNjdmMDk3ZDdlNjUxNDY5NmZkODMwODA3OTRiNmI1Y2E2NjQwMDFmMmVkZTRmZDZmMDFkYTQ5MzNkYjg3NWZmMDI4ZmVjNDJiMjlmYzU1MjQ5NDFlMGYyMDgzMGYiLAogICAgICAgICJwcm9vZk9mUG9zc2Vzc2lvbiI6ICIweDgyMzUyYWUxZTAxMDM4MTczZTkyZTA4OGJkMzRjMmJlZTljYzRiMzRkZjVjNWU4YmQyNzczY2VmOTIzOGVlZjg3MGMyZjkzZmE4OTYwNzMzMmNjYmI4NGFhNjY2MDhjNzA2YjdjMmYxMjdiOGI4MGM0NjFjMDRiYmM2MDgyYWZiZmZlMjIwYWFjNzlmNjY1MzNlYTdjNjNmMDQ1MWQ3ZDMyNDU2MzY5ZGQzMzVjOTcxMDkzOGVlNDExMWQwOGQ3OSIKICAgICAgfQogICAgfSwKICAgIHsKICAgICAgIm5vZGVJRCI6ICJOb2RlSUQtNzhpYldwanRaejVaR1Q2RXlURWR1OFZLbWJvVUhUdUdUIiwKICAgICAgInJld2FyZEFkZHJlc3MiOiAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiLAogICAgICAiZGVsZWdhdGlvbkZlZSI6IDYyNTAwLAogICAgICAic2lnbmVyIjogewogICAgICAgICJwdWJsaWNLZXkiOiAiMHg4MzI3ZGJlMWJhNDExYzI3MDYzN2IwODBhODQ3MWZiNDFlZWI4YTliMzkxN2FmMDcyNzUwMWVmOGJkYWE5MDFkMDYzNzgwYmQ3MDJmMzBmNDU4YTYxZjNkNDI5N2RjOTgiLAogICAgICAgICJwcm9vZk9mUG9zc2Vzc2lvbiI6ICIweGE5YzAzOWI1NzY1YWIwNjhiZDYzMmJiY2RjOWJjMmE1M2YyOWUyYzU2YjMzZTMwZDczMmEyM2Q4YzQzMGQ1M2VmNDdlYmNjZmFhNWNmY2VkZDhmMDQxYzJjMTM0OGYwYjBlYWM0MTMxOTJiNzU0NGQyODRmODJkMWZhMGY3NGY5OGQ1ODA1OTA1MzYzYjgxODZlZmRlZjZlNzcxODJmYjFlNzE0N2Y4NTExZTkwMGQxOTVkYjA2ZGE2YTIyZjBhMCIKICAgICAgfQogICAgfSwKICAgIHsKICAgICAgIm5vZGVJRCI6ICJOb2RlSUQtTDRDWThCNXVWU0RlNGNuTjFCcGVEc0hhY01wNHE0cThxIiwKICAgICAgInJld2FyZEFkZHJlc3MiOiAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiLAogICAgICAiZGVsZWdhdGlvbkZlZSI6IDYyNTAwLAogICAgICAic2lnbmVyIjogewogICAgICAgICJwdWJsaWNLZXkiOiAiMHhhOThjNjQ2YThjODYyZWMxNTMyNmU0Y2ZlMmEwZjY2YThmYjdjZjU1NTc2NWY4M2ZmMzIwYTFhNzYyNjgyMjhmM2M4YjI2MmQxZGU0MDA4ZTBiYTQ5YTg5Y2ZhYmZiOTUiLAogICAgICAgICJwcm9vZk9mUG9zc2Vzc2lvbiI6ICIweDk1YzgxNmE0ZDI5MmE0N2M0ZDk5MzRlMzU4NjIxZDA3ODVmZDI5MjBhMWEzMDRjZWJjOWI3ZTQ3NDE3ZTNmZmY3ODBjZmNkZGY3Y2ExMTc3YjQ1YmJmYWMyZjk5Nzg1ODE3NjFkOWRkZDU1ZWM2MTQyZDkyOTk4ZWVhZGJhZmU4Y2Q3NjUxMDU2ZmJiNzlhZmVhNjQzZjBjZDIwZmY0ZjYzODlkZGQ5MWVlMmRiNDU3OTQzOGE2OTA4NjA5YjRjMSIKICAgICAgfQogICAgfSwKICAgIHsKICAgICAgIm5vZGVJRCI6ICJOb2RlSUQtUDVRR0g0RVhkZHJjeU5BemtxeVpLSFhnRXBWWDZIRXhMIiwKICAgICAgInJld2FyZEFkZHJlc3MiOiAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiLAogICAgICAiZGVsZWdhdGlvbkZlZSI6IDYyNTAwLAogICAgICAic2lnbmVyIjogewogICAgICAgICJwdWJsaWNLZXkiOiAiMHhiMGQzNWNjZjcwYTZkODRlMmJjYTFkYzE2NmE0YzMzMjRkN2VkZDg2ZTg3OWFkZDJiYTY1MTFjOGVmNmJmZDg5YTE1NTM0ZTY3NDY3Y2NkOWM5MjExNTM0YjMzMjk1YTEiLAogICAgICAgICJwcm9vZk9mUG9zc2Vzc2lvbiI6ICIweGE0YjQ4MGE5YTA3YjRhOTc2MzBkZDlkNmUyYmY0ODM0YTNlNjcwNDQ4YjU1NzVlM2JhNzJhMDNlMDZlYzUwOWVjODU5ODQwYTExMDRiMThmMGNkNTQ2OTZlNmY5OGFkYjBlOWY1MjYwMTYxMzMyZmUzMmE1MGNiMWE5ODA2YjFiNTAyNTAzNzczMWVhNzdjNjQxZDYwN2ZkMDU4NGNlMjlkNzk1NGY1ZThmNTIzYzEzYTJlNTczMjUxMTIyN2Y1MCIKICAgICAgfQogICAgfSwKICAgIHsKICAgICAgIm5vZGVJRCI6ICJOb2RlSUQtN2VSdm5mczJhMlB2clBIVXVDUlJwUFZBb1ZqYld4YUZHIiwKICAgICAgInJld2FyZEFkZHJlc3MiOiAiWC1jdXN0b20xdjZ2dXd4amdyMDQzc2cwbnV1aHE3MGs2dmdudWxlNjkyZnZucjkiLAogICAgICAiZGVsZWdhdGlvbkZlZSI6IDYyNTAwLAogICAgICAic2lnbmVyIjogewogICAgICAgICJwdWJsaWNLZXkiOiAiMHhhNWZkZTYwNDJjNmUwZWU0ODJmNDYzNDZkZjA0NjAwMGNkNTdkZDg3OGQzNjYzN2E1YTYyYWRlYzA3YTUxMTRjZGVlYTA5NGE4NWY0ZjcyYjQ2NjQ1Zjk0ZTkwNzY2OTIiLAogICAgICAgICJwcm9vZk9mUG9zc2Vzc2lvbiI6ICIweDkxNDQ0MjMwYzVjZWI4ZTUxNjQyMTM5ZTE4NDJiNTZmODU2Mzg2NTM3NmI2ZjQyMDViZWNhNGRjMGJiMGJjNGIzMTRiY2UxZTE5ZTNiNTQyYTM5NDFlY2U1MWFlMjA1ZTAzYTA5NDgyNGZlZTI4ZjlmNzAyZWQzMTA3NTZmMDYzN2JmMTY2MzcxNjU2ZTFjM2ViOTAwMWRmODlmNGNkY2NjNzM0MTAyNDJhNmQ4NzVlYjYzNjNkMTJiY2U0MDMxNiIKICAgICAgfQogICAgfQogIF0sCiAgImNDaGFpbkdlbmVzaXMiOiAie1wiY29uZmlnXCI6e1wiY2hhaW5JZFwiOjQzMTE3LFwiaG9tZXN0ZWFkQmxvY2tcIjowLFwiZGFvRm9ya0Jsb2NrXCI6MCxcImRhb0ZvcmtTdXBwb3J0XCI6dHJ1ZSxcImVpcDE1MEJsb2NrXCI6MCxcImVpcDE1MEhhc2hcIjpcIjB4MjA4Njc5OWFlZWJlYWUxMzVjMjQ2YzY1MDIxYzgyYjRlMTVhMmM0NTEzNDA5OTNhYWNmZDI3NTE4ODY1MTRmMFwiLFwiZWlwMTU1QmxvY2tcIjowLFwiZWlwMTU4QmxvY2tcIjowLFwiYnl6YW50aXVtQmxvY2tcIjowLFwiY29uc3RhbnRpbm9wbGVCbG9ja1wiOjAsXCJwZXRlcnNidXJnQmxvY2tcIjowLFwiaXN0YW5idWxCbG9ja1wiOjAsXCJtdWlyR2xhY2llckJsb2NrXCI6MH0sXCJub25jZVwiOlwiMHgwXCIsXCJ0aW1lc3RhbXBcIjpcIjB4MFwiLFwiZXh0cmFEYXRhXCI6XCIweDAwXCIsXCJnYXNMaW1pdFwiOlwiMHg1ZjVlMTAwXCIsXCJkaWZmaWN1bHR5XCI6XCIweDBcIixcIm1peEhhc2hcIjpcIjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMFwiLFwiY29pbmJhc2VcIjpcIjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMFwiLFwiYWxsb2NcIjp7XCIwMTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwXCI6e1wiY29kZVwiOlwiMHg3MzAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAzMDE0NjA4MDYwNDA1MjYwMDQzNjEwNjAzZDU3NjAwMDM1NjBlMDFjODA2MzFlMDEwNDM5MTQ2MDQyNTc4MDYzYjY1MTBiYjMxNDYwNmU1NzViNjAwMDgwZmQ1YjYwNWM2MDA0ODAzNjAzNjAyMDgxMTAxNTYwNTY1NzYwMDA4MGZkNWI1MDM1NjBiMTU2NWI2MDQwODA1MTkxODI1MjUxOTA4MTkwMDM2MDIwMDE5MGYzNWI4MTgwMTU2MDc5NTc2MDAwODBmZDViNTA2MGFmNjAwNDgwMzYwMzYwODA4MTEwMTU2MDhlNTc2MDAwODBmZDViNTA2MDAxNjAwMTYwYTAxYjAzODEzNTE2OTA2MDIwODEwMTM1OTA2MDQwODEwMTM1OTA2MDYwMDEzNTYwYjY1NjViMDA1YjMwY2Q5MDU2NWI4MzYwMDE2MDAxNjBhMDFiMDMxNjgxODM2MTA4ZmM4NjkwODExNTAyOTA2MDQwNTE2MDAwNjA0MDUxODA4MzAzODE4ODg4ODc4YzhhY2Y5NTUwNTA1MDUwNTA1MDE1ODAxNTYwZjQ1NzNkNjAwMDgwM2UzZDYwMDBmZDViNTA1MDUwNTA1MDU2ZmVhMjY0Njk3MDY2NzM1ODIyMTIyMDFlZWJjZTk3MGZlM2Y1Y2I5NmJmOGFjNmJhNWY1YzEzM2ZjMjkwOGFlM2RjZDUxMDgyY2ZlZThmNTgzNDI5ZDA2NDczNmY2YzYzNDMwMDA2MGEwMDMzXCIsXCJiYWxhbmNlXCI6XCIweDBcIn0sXCIweDY0M0YyNDU0NDMwRTIxODc1MGI1ZTY1MzNkOUMwZTBEZDUwQjhkNjhcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweGY5QkZBNEM0NWE4ZDgzMGE1OTFCMzM3NDMyMGZkOENDRjNGRDc1RDRcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweEQ5ZDRmMTZhNzFFMjNlRGY4ZTJGMmExRWJlY2Q0NkIwMzE3N2EyMmNcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweDJhMTc4MzE0MjViYzZEMjAwODREMTUyNmIxMDAxQzQ1MUVENEM0QTdcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweDdjNUE4NjM5RjFlODZGMTM0ZjFFNDIzOTQyOWY3NTZBMTQ0MWUzMjJcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweGZEREVmNWNiMEQwOUU0ODNkQkFCNTg3QkE5NTg2NTdCNzlBNDJFNThcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweEI0Y0E2QzEyMUQ2Mjg3YWY3YWM3Y2I2MkFlMzNkMmIwNTRiOUZDNDRcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn0sXCIweEM3MUE2MWE4MTVlNDlkMTZDNDI1NDgyQTM0MmEzNjdDRDQyRTM4YTZcIjp7XCJiYWxhbmNlXCI6XCIweDE0MzFFMEZBRTZENzIxN0NBQTAwMDAwMDBcIn19LFwibnVtYmVyXCI6XCIweDBcIixcImdhc1VzZWRcIjpcIjB4MFwiLFwicGFyZW50SGFzaFwiOlwiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwXCJ9IiwKICAibWVzc2FnZSI6ICJFdG5hIGhlcmUgd2UgY29tZSIKfQo="
```

## Getting DevNet AVAX

You can get 5 AVAX on the DevNet C-Chain from the public faucet [here](https://core.app/tools/testnet-faucet/?subnet=cdevnet&token=cdevnet).

If you require more, please reach out to our team in the relevant Telegram/Slack channel.

## Testing New Transaction Types

Ensure you are locally running an AvalancheGo node. 

The AvalancheGo [wallet](https://github.com/ava-labs/avalanchego/tree/v1.12.0-initial-poc.5/wallet) folder contains example scripts written in Golang for testing common workflows. 

Inside `wallet/subnet/primary/examples` [folder](https://github.com/ava-labs/avalanchego/tree/v1.12.0-initial-poc.5/wallet/subnet/primary/examples), you will find scripts you can run locally. To test creating a new L1 on Etna Devnet, you can run:

`go run wallet/subnet/primary/examples/create-subnet/main.go`

`go run wallet/subnet/primary/examples/create-chain/main.go`

`go run wallet/subnet/primary/examples/convert-subnet/main.go`

Some of the values in the scripts (such as `subnetID`) are hard-coded, and will need to be adjusted based on the output of each consecutive transaction in order to function properly. 
