package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"mypkg/helpers"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/crypto"
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

	if blockHeight >= 3 {
		fmt.Printf("Block height is already greater than or equal to 3, skipping activation\n")
		return nil
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.NonceAt(context.Background(), address, nil)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	for i := 0; i < 2; i++ {
		tx := types.NewTransaction(
			nonce+uint64(i)+2,
			address,
			big.NewInt(1),
			21000,
			big.NewInt(225_000_000_000),
			nil,
		)

		signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(evmChainID), key)
		if err != nil {
			return fmt.Errorf("failed to sign transaction: %w", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return fmt.Errorf("failed to send transaction: %w", err)
		}

		fmt.Printf("Sent transaction %d: %s\n", i+1, signedTx.Hash().String())
	}

	time.Sleep(4 * time.Second)

	blockHeight, err = client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get final block height: %w", err)
	}
	fmt.Printf("Final block height: %d\n", blockHeight)
	return nil
}
