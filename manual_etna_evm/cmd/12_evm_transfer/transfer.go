package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"mypkg/lib"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	goethereumtypes "github.com/ethereum/go-ethereum/core/types"
	goethereumcrypto "github.com/ethereum/go-ethereum/crypto"
	goethereumethclient "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	chainIDBytes, err := os.ReadFile("data/chain.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	// Load validator manager key to use as source of funds
	keyHex, err := os.ReadFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("failed to read key file: %s\n", err)
	}
	keyBytes, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
	if err != nil {
		log.Fatalf("failed to decode key: %s\n", err)
	}

	key, err := goethereumcrypto.ToECDSA(keyBytes)
	if err != nil {
		log.Fatalf("failed to convert to ECDSA key: %s\n", err)
	}

	// Generate random destination address
	destKey, err := goethereumcrypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate random key: %s\n", err)
	}
	destAddr := goethereumcrypto.PubkeyToAddress(destKey.PublicKey)

	node0URL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	client, err := goethereumethclient.Dial(node0URL)
	if err != nil {
		log.Fatalf("failed to connect to node0: %s\n", err)
	}

	// Get the sender's address and nonce
	fromAddress := goethereumcrypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
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
	tx := goethereumtypes.NewTransaction(nonce, destAddr, value, gasLimit, gasPrice, nil)
	chainIDInt, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get chain ID: %s\n", err)
	}

	signedTx, err := goethereumtypes.SignTx(tx, goethereumtypes.NewEIP155Signer(chainIDInt), key)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s\n", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("failed to send transaction: %s\n", err)
	}

	fmt.Printf("✅ Sent 1 AVAX from %s to %s in tx %s\n", fromAddress.Hex(), destAddr.Hex(), signedTx.Hash().Hex())

	// Wait for transaction to be mined
	time.Sleep(5 * time.Second)

	// Check balance on all nodes
	for nodeNumber := 0; nodeNumber < lib.VALIDATORS_COUNT; nodeNumber++ {

		nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", fmt.Sprintf("%d", 9650+nodeNumber*2), chainID)
		client, err := goethereumethclient.Dial(nodeURL)
		if err != nil {
			log.Fatalf("failed to connect to node%d: %s\n", nodeNumber, err)
		}

		balance, err := client.BalanceAt(context.Background(), destAddr, nil)
		if err != nil {
			log.Fatalf("failed to get balance from node%d: %s\n", nodeNumber, err)
		}

		if balance.Cmp(value) != 0 {
			log.Fatalf("❌ Balance on node%d is %s, expected %s\n", nodeNumber, balance, value)
		}
		fmt.Printf("✅ Balance on node%d matches expected value\n", nodeNumber)
	}
}
