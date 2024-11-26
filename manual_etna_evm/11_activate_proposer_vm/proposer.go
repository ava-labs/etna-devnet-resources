package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/helpers"

	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/subnet-evm/ethclient"
)

func main() {
	key, err := helpers.LoadValidatorManagerKeyECDSA()
	if err != nil {
		log.Fatalf("failed to load validator manager key: %s\n", err)
	}

	chainID, err := helpers.LoadId("chain")
	if err != nil {
		log.Fatalf("failed to load chain ID: %s\n", err)
	}

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("failed to connect to node0: %s\n", err)
	}

	evmChainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("failed to get chain ID: %s\n", err)
	}

	blockHeight, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("failed to get initial block height: %s\n", err)
	}
	fmt.Printf("Initial block height: %d\n", blockHeight)

	//FIXME: How to check if the fork is already activated? PRs are welcome!
	if blockHeight >= 3 {
		fmt.Printf("Block height is already greater than or equal to 3, skipping activation\n")
		return
	}

	if err := evm.IssueTxsToActivateProposerVMFork(client, evmChainID, key); err != nil {
		log.Fatalf("failed to activate proposer VM fork: %s\n", err)
	}

	blockHeight, err = client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("failed to get final block height: %s\n", err)
	}
	fmt.Printf("Final block height: %d\n", blockHeight)
	fmt.Println("✅ Successfully activated proposer VM fork")
}
