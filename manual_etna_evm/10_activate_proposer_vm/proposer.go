package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/helpers"
	"time"

	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/subnet-evm/ethclient"
)

const maxAttempts = 10

func main() {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		lastErr = activateProposerVM()
		if lastErr == nil {
			fmt.Println("✅ Successfully activated proposer VM fork")
			return
		}
		fmt.Printf("Attempt %d/%d of activating proposerVM failed: %s\n", i+1, maxAttempts, lastErr)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	log.Fatalf("failed to activate proposer VM after %d attempts: %s\n", maxAttempts, lastErr)
}

func activateProposerVM() error {
	key, err := helpers.LoadValidatorManagerKeyECDSA()
	if err != nil {
		return fmt.Errorf("failed to load validator manager key: %w", err)
	}

	chainID, err := helpers.LoadId("chain")
	if err != nil {
		return fmt.Errorf("failed to load chain ID: %w", err)
	}

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return fmt.Errorf("failed to connect to node0: %w", err)
	}

	evmChainID, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	blockHeight, err := client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get initial block height: %w", err)
	}
	fmt.Printf("Initial block height: %d\n", blockHeight)

	//FIXME: How to check if the fork is already activated? PRs are welcome!
	if blockHeight >= 3 {
		fmt.Printf("Block height is already greater than or equal to 3, skipping activation\n")
		return nil
	}

	if err := evm.IssueTxsToActivateProposerVMFork(client, evmChainID, key); err != nil {
		return fmt.Errorf("failed to activate proposer VM fork: %w", err)
	}

	blockHeight, err = client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get final block height: %w", err)
	}
	fmt.Printf("Final block height: %d\n", blockHeight)
	fmt.Println("✅ Successfully activated proposer VM fork")
	return nil
}
