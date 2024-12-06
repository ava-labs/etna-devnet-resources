package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const maxAttempts = 10

func main() {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {

		lastErr = activateProposerVM()
		if lastErr == nil {
			log.Println("✅ Successfully activated proposer VM fork")
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

	err = IssueTxsToActivateProposerVMFork(client, evmChainID, key)
	if err != nil {
		return fmt.Errorf("failed to issue transactions to activate proposer VM fork: %w", err)
	}

	blockHeight, err = client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get initial block height: %w", err)
	}
	fmt.Printf("Block height after activation: %d\n", blockHeight)

	return nil

}

func IssueTxsToActivateProposerVMFork(
	client ethclient.Client,
	chainID *big.Int,
	privKey *ecdsa.PrivateKey,
) error {
	const (
		repeatsOnFailure    = 3
		sleepBetweenRepeats = 1 * time.Second
	)

	var errorList []error
	var err error
	for i := 0; i < repeatsOnFailure; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err = issueTxsToActivateProposerVMFork(client, ctx, chainID, privKey)
		if err == nil {
			break
		}
		err = fmt.Errorf(
			"failure issuing txs to activate proposer VM fork for client %#v: %w",
			client,
			err,
		)
		errorList = append(errorList, err)
		time.Sleep(sleepBetweenRepeats)
	}
	if err != nil {
		for _, indivError := range errorList {
			log.Printf("Error: %s", indivError)
		}
	}
	return err
}

func issueTxsToActivateProposerVMFork(
	client ethclient.Client,
	ctx context.Context,
	chainID *big.Int,
	fundedKey *ecdsa.PrivateKey,
) error {
	const numTriggerTxs = 2 // Number of txs needed to activate the proposer VM fork
	addr := crypto.PubkeyToAddress(fundedKey.PublicKey)
	gasPrice := big.NewInt(225_000_000_000) // 225 Gwei

	txSigner := types.LatestSignerForChainID(chainID)
	for i := 0; i < numTriggerTxs; i++ {
		prevBlockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		nonce, err := client.NonceAt(ctx, addr, nil)
		if err != nil {
			return err
		}
		tx := types.NewTransaction(
			nonce, addr, common.Big1, 21000, gasPrice, nil)
		triggerTx, err := types.SignTx(tx, txSigner, fundedKey)
		if err != nil {
			return err
		}
		if err := client.SendTransaction(ctx, triggerTx); err != nil {
			return err
		}
		if err := WaitForNewBlock(client, ctx, prevBlockNumber, 10*time.Second, time.Second); err != nil {
			return err
		}
	}
	return nil
}

func WaitForNewBlock(
	client ethclient.Client,
	ctx context.Context,
	prevBlockNumber uint64,
	totalDuration time.Duration,
	stepDuration time.Duration,
) error {
	steps := totalDuration / stepDuration
	for seconds := 0; seconds < int(steps); seconds++ {
		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		if blockNumber > prevBlockNumber {
			return nil
		}
		time.Sleep(stepDuration)
	}
	return fmt.Errorf("new block not produced in %f seconds", totalDuration.Seconds())
}
