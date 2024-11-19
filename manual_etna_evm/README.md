## Create L1 manually with Go

This repo helps you create and update L1 on Avalanche after Etna upgrade. But heads up - you probably want to use `avalanche-cli` instead (it's much easier!).

Quick start:
- Run `./run.sh` to start a new L1 on Devnet
- Run `./cleanup.sh` to clean up (keeps your keys)

**What this does**: Creates a single-validator L1 chain, uses Devnet and gets AVAX automatically from the Ewoq key

TODO (PRs welcome):
- Add and remove validators

Example log below:
```bash
vscode ➜ .../Projects/ava-labs/etna-devnet-resources/manual_etna_evm (manual-checrry-pick) $ 
echo -e "\n🔑 Generating keys\n"
go run ./cmd/01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./cmd/02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./cmd/03_create_subnet/

echo -e "\n🧱 Generating genesis\n"
go run ./cmd/04_gen_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./cmd/05_create_chain/

echo -e "\n🏗️  Setting up node configs\n"
go run ./cmd/06_node_configs/

echo -e "\n🚀 Launching nodes\n"
export CURRENT_UID=$(id -u)
export CURRENT_GID=$(id -g)
docker compose -f ./cmd/07_launch_nodes/docker-compose.yml up -d --build

echo -e "\n🔄 Converting chain\n"
go run ./cmd/08_convert_chain/

echo -e "\n🔄 Updating node configs\n"
go run ./cmd/12_evm_transfer/st coins\n"docker-compose.yml up -d

🔑 Generating keys

2024/11/19 08:11:33 ✅ POA validator manager keys generated and saved in ./data/ folder

💰 Checking balance

2024/11/19 08:11:43 fetched state in 4.796712772s
2024/11/19 08:11:43 P-chain balance: 0
2024/11/19 08:11:43 P-chain balance insufficient on address Gd7pXyUgfvNjiBDrKXimbQRatus81yEW9: 0 < 1100000000
2024/11/19 08:11:43 Balance on c-chain at address 0xfF5a39ca8679a4Fe4304238683001322AE2d15B6: 0
2024/11/19 08:11:43 Balance 0 is less than minimum balance: 1100000000
2024/11/19 08:11:44 Transaction sent: 0x953a9a52f8fcfecc619b74c6c391921bd311242bfbc184553be15bb95e8a9941
2024/11/19 08:11:44 Transferring 1100000000 from C-chain to P-chain
constants.PlatformChainID 11111111111111111111111111111111LpoYY
2024/11/19 08:11:50 ✅ Issued export u6K5QXJ9pQ2YoaTvFwNHYCsRpyBd1dqipGTs9SNY3exfZeXMi
2024/11/19 08:11:50 ✅ Issued import Gc9JTJMzYeJ8CZnmcLG6dV7KNWgSpPFxCQgRzVodo8c2HpBRt
2024/11/19 08:11:55 fetched state in 4.631953649s
2024/11/19 08:11:55 ✅ Final P-chain balance: 1199993330 (greater than minimum 1100000000)

🕸️  Creating subnet

2024/11/19 08:11:57 Subnet ID file does not exist, let's create a new subnet
2024/11/19 08:12:02 Synced wallet in 4.964308294s
2024/11/19 08:12:02 ✅ Created new subnet Wh3CaXkg2mukJMdEray66E9JQ7zGLzymWCu9wDX28PRyBmLKd in 404.847662ms

🧱 Generating genesis

2024/11/19 08:12:04 ✅ Successfully wrote genesis to data/L1-genesis.json

⛓️  Creating chain

2024/11/19 08:12:05 🔍 Chain ID file does not exist, let's create a new chain
2024/11/19 08:12:05 Using vmID: srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
2024/11/19 08:12:10 synced wallet in 5.141260108s
2024/11/19 08:12:11 ✅ Created new chain hKr97zf4cRYEfb2fcgTaLAXWL3H1Mo6Pn7GiAh8HyZwi4hRZ2 in 627.062709ms
2024/11/19 08:12:11 Saved chain ID to file data/chain.txt

🏗️  Setting up node configs

✅ Successfully created configs

🚀 Launching nodes

[+] Building 1.8s (22/22) FINISHED                                                                                                                   docker:default
 => [node0 internal] load build definition from Dockerfile                                                                                                     0.0s
 => => transferring dockerfile: 1.11kB                                                                                                                         0.0s
 => [node0 internal] load metadata for docker.io/library/debian:bookworm-slim                                                                                  1.6s
 => [node0 internal] load metadata for docker.io/library/golang:1.22-bookworm                                                                                  1.6s
 => [node0 auth] library/golang:pull token for registry-1.docker.io                                                                                            0.0s
 => [node0 auth] library/debian:pull token for registry-1.docker.io                                                                                            0.0s
 => [node0 internal] load .dockerignore                                                                                                                        0.0s
 => => transferring context: 2B                                                                                                                                0.0s
 => [node0 subnet-evm-builder 1/6] FROM docker.io/library/golang:1.22-bookworm@sha256:475ff60e52faaf037be2e7a1bc2ea5ea4aaa3396274af3def6545124a18b99b4         0.0s
 => [node0 stage-2 1/4] FROM docker.io/library/debian:bookworm-slim@sha256:ca3372ce30b03a591ec573ea975ad8b0ecaf0eb17a354416741f8001bbcae33d                    0.0s
 => CACHED [node0 stage-2 2/4] RUN groupadd -r nobody || true                                                                                                  0.0s
 => CACHED [node0 subnet-evm-builder 2/6] WORKDIR /app                                                                                                         0.0s
 => CACHED [node0 avalanchego-builder 3/6] RUN git clone https://github.com/ava-labs/avalanchego.git                                                           0.0s
 => CACHED [node0 avalanchego-builder 4/6] WORKDIR /app/avalanchego                                                                                            0.0s
 => CACHED [node0 avalanchego-builder 5/6] RUN git checkout v1.12.0-initial-poc.6                                                                              0.0s
 => CACHED [node0 avalanchego-builder 6/6] RUN ./scripts/build.sh                                                                                              0.0s
 => CACHED [node0 stage-2 3/4] COPY --from=avalanchego-builder /app/avalanchego/build/avalanchego /usr/local/bin/avalanchego                                   0.0s
 => CACHED [node0 subnet-evm-builder 3/6] RUN git clone https://github.com/ava-labs/subnet-evm.git                                                             0.0s
 => CACHED [node0 subnet-evm-builder 4/6] WORKDIR /app/subnet-evm                                                                                              0.0s
 => CACHED [node0 subnet-evm-builder 5/6] RUN git checkout v0.6.11                                                                                             0.0s
 => CACHED [node0 subnet-evm-builder 6/6] RUN go build -v -o /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy ./plugin                                   0.0s
 => CACHED [node0 stage-2 4/4] COPY --from=subnet-evm-builder /app/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTW  0.0s
 => [node0] exporting to image                                                                                                                                 0.0s
 => => exporting layers                                                                                                                                        0.0s
 => => writing image sha256:0556dd4528411822505a083f2972d36ef809df256582ed2ec20039e50f813aa4                                                                   0.0s
 => => naming to docker.io/library/07_launch_nodes-node0                                                                                                       0.0s
 => [node0] resolving provenance for metadata file                                                                                                             0.0s
[+] Running 1/1
 ✔ Container node0  Started                                                                                                                                    0.2s 

🔄 Converting chain

Using changeOwnerAddress: P-custom14d0gm448df670hde98uvgh0vfxe79nkxu4a8e8
Getting node info from http://127.0.0.1:9650
✅ Convert subnet tx ID: nSPEner38dZtu7Ar3PTDr6DE9mVtuDfkqxXmADb71JSWEqBHx

🔄 Updating node configs

✅ Successfully updated node configs

🚀 Stopping nodes

[+] Running 1/0
 ✔ Container node0  Removed                                                                                                                                    0.1s 

🚀 Starting nodes again with a new subnet

[+] Running 1/1
 ✔ Container node0  Started                                                                                                                                    0.1s 

🏥 Checking subnet health

Checking RPC endpoint for node0: http://127.0.0.1:9652/ext/bc/hKr97zf4cRYEfb2fcgTaLAXWL3H1Mo6Pn7GiAh8HyZwi4hRZ2/rpc
🌱 Node0 is starting up (try 1 of 60) - waiting for chain ID...
🌱 Node0 is starting up (try 2 of 60) - waiting for chain ID...
🌱 Node0 RPC endpoint is healthy (chain ID: 12345)
✅ All nodes are healthy!

💸 Sending some test coins

✅ Sent 1 AVAX from 0xfF5a39ca8679a4Fe4304238683001322AE2d15B6 to 0xfA0382fBE355c7A3A1cd1c10FCE02342d9b68B7A in tx 0x5d10b2580c56910671adb5177dec57184aa0dcba518910ecd09e87591dfdf4fd
✅ Balance on node0 matches expected value

echo -e "\n🔄 Initializing PoA\n"
go run ./cmd/13_init_poa/

🔄 Initializing PoA

{"level":"debug","timestamp":"2024-11-19T08:13:20.997Z","logger":"p2p-network","caller":"dialer/dialer.go:52","msg":"creating dialer","throttleRPS":50,"dialTimeout":30000000000}
{"level":"info","timestamp":"2024-11-19T08:13:20.998Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:348","msg":"Signing subnet not found, requesting from PChain","blockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:21.001Z","logger":"init-aggregator","caller":"cache/cache.go:51","msg":"cache miss","msgID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk"}
{"level":"debug","timestamp":"2024-11-19T08:13:21.002Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:219","msg":"Aggregator collecting signatures from peers.","attempt":1,"sourceBlockchainID":"11111111111111111111111111111111LpoYY","signingSubnetID":"Wh3CaXkg2mukJMdEray66E9JQ7zGLzymWCu9wDX28PRyBmLKd","validatorSetSize":1,"signatureMapSize":0,"responsesExpected":1}
{"level":"debug","timestamp":"2024-11-19T08:13:21.002Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:239","msg":"Added node ID to query.","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:21.002Z","logger":"p2p-network","caller":"peers/external_handler.go:132","msg":"Registering request ID","requestID":3887330499}
{"level":"debug","timestamp":"2024-11-19T08:13:21.002Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:259","msg":"Sent signature request to network","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sentTo":[],"sourceBlockchainID":"11111111111111111111111111111111LpoYY","sourceSubnetID":"11111111111111111111111111111111LpoYY","signingSubnetID":"Wh3CaXkg2mukJMdEray66E9JQ7zGLzymWCu9wDX28PRyBmLKd"}
{"level":"warn","timestamp":"2024-11-19T08:13:21.002Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:269","msg":"Failed to make async request to node","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd"}
{"level":"debug","timestamp":"2024-11-19T08:13:21.006Z","logger":"p2p-network","caller":"peers/external_handler.go:104","msg":"Connected","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd","version":"avalanchego/1.11.12","subnetID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.003Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:219","msg":"Aggregator collecting signatures from peers.","attempt":2,"sourceBlockchainID":"11111111111111111111111111111111LpoYY","signingSubnetID":"Wh3CaXkg2mukJMdEray66E9JQ7zGLzymWCu9wDX28PRyBmLKd","validatorSetSize":1,"signatureMapSize":0,"responsesExpected":1}
{"level":"debug","timestamp":"2024-11-19T08:13:23.003Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:239","msg":"Added node ID to query.","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.003Z","logger":"p2p-network","caller":"peers/external_handler.go:132","msg":"Registering request ID","requestID":3887330499}
{"level":"debug","timestamp":"2024-11-19T08:13:23.003Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:259","msg":"Sent signature request to network","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sentTo":["NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd"],"sourceBlockchainID":"11111111111111111111111111111111LpoYY","sourceSubnetID":"11111111111111111111111111111111LpoYY","signingSubnetID":"Wh3CaXkg2mukJMdEray66E9JQ7zGLzymWCu9wDX28PRyBmLKd"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.006Z","logger":"p2p-network","caller":"peers/external_handler.go:87","msg":"Handling app response","op":"app_response","from":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.006Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:282","msg":"Processing response from node","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.010Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:406","msg":"Got valid signature response","nodeID":"NodeID-A5GAN71mt57Mg9omtJC4vZmPB5qhSE4wd","stakeWeight":100,"warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
{"level":"debug","timestamp":"2024-11-19T08:13:23.010Z","logger":"init-aggregator","caller":"cache/cache.go:51","msg":"cache miss","msgID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk"}
{"level":"info","timestamp":"2024-11-19T08:13:23.010Z","logger":"init-aggregator","caller":"aggregator/aggregator.go:311","msg":"Created signed message.","warpMessageID":"2Ef6kNFW4E2JVdcG9m6h6w8PQGtCXqGjdX67xxNe5h8i3Pxwbk","signatureWeight":100,"sourceBlockchainID":"11111111111111111111111111111111LpoYY"}
✅ Successfully initialized Proof of Authority
```
