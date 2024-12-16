package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"crypto/rand"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const maxAttempts = 10

func main() {
	key := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)
	chainID := helpers.LoadId(helpers.ChainIdPath)

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("failed to connect to node0: %w", err)
	}

	evmChainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("failed to get chain ID: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	addr := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.NonceAt(ctx, addr, nil)
	if err != nil {
		log.Fatalf("failed to get nonce: %w", err)
	}

	gasPrice := big.NewInt(225_000_000_000)      // 225 Gwei
	value := big.NewInt(100_000_000_000_000_000) // 0.1 AVAX in wei

	// Generate random address
	var randomAddr common.Address
	if _, err := rand.Read(randomAddr[:]); err != nil {
		log.Fatalf("failed to generate random address: %w", err)
	}

	// Check balance of random address
	balance, err := client.BalanceAt(ctx, randomAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance of random address: %w", err)
	}
	log.Printf("Random address %s has balance of %s wei", randomAddr.Hex(), balance.String())

	tx := types.NewTransaction(
		nonce,
		randomAddr, // Use random address instead of hardcoded one
		value,
		21000,
		gasPrice,
		nil,
	)

	txSigner := types.LatestSignerForChainID(evmChainID)
	signedTx, err := types.SignTx(tx, txSigner, key)
	if err != nil {
		log.Fatalf("failed to sign transaction: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to send transaction: %w", err)
	}

	// Wait a bit for the transaction to be processed
	time.Sleep(2 * time.Second)

	// Check balance of random address after sending
	balanceAfter, err := client.BalanceAt(ctx, randomAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance of random address after transfer: %w", err)
	}

	// Verify the balance matches what we sent
	if balanceAfter.Cmp(value) != 0 {
		log.Fatalf("unexpected balance after transfer: got %s, want %s", balanceAfter.String(), value.String())
	}

	log.Printf("✅ Successfully sent %s wei to random address %s", value.String(), randomAddr.Hex())
}
