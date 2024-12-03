package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"mypkg/helpers"
	"time"

	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/ethclient"
	goethereumcrypto "github.com/ethereum/go-ethereum/crypto"
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

	// Generate random destination address
	destKey, err := goethereumcrypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate random key: %s\n", err)
	}
	destAddr := goethereumcrypto.PubkeyToAddress(destKey.PublicKey)
	fmt.Printf("Generated destination address: %s\n", destAddr.Hex())

	node0URL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	client, err := ethclient.Dial(node0URL)
	if err != nil {
		log.Fatalf("failed to connect to node0: %s\n", err)
	}

	// Get the sender's address and nonce
	fromAddress := goethereumcrypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("failed to get nonce: %s\n", err)
	}

	// Send 1 AVAX
	value := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to get gas price: %s\n", err)
	}

	// Create and sign transaction
	tx := types.NewTransaction(nonce, destAddr, value, gasLimit, gasPrice, nil)
	chainIDInt, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get chain ID: %s\n", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainIDInt), key)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s\n", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("failed to send transaction: %s\n", err)
	}

	fmt.Printf("✅ Sent 1 AVAX from %s to %s in tx %s\n", fromAddress.Hex(), destAddr.Hex(), signedTx.Hash().Hex())

	time.Sleep(10 * time.Second)

	balance, err := client.BalanceAt(context.Background(), destAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %s\n", err)
	}

	if balance.Cmp(value) != 0 {
		log.Fatalf("❌ Balance is %s, expected %s\n", balance, value)
	}
	fmt.Printf("✅ Balance matches expected value: %s\n", balance)
}
