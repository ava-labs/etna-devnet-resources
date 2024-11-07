package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"
	"time"

	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/avalanchego/ids"
)

func checkNodeHealth(nodeNumber int, rpcURL string) error {
	const maxAttempts = 60
	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		client, err := evm.GetClient(rpcURL)
		if err != nil {
			fmt.Printf("🌱 Node%d is still booting up (try %d of %d) - waiting for RPC endpoint...\n", nodeNumber, i+1, maxAttempts)
			time.Sleep(2 * time.Second)
			continue
		}

		chainID, err := client.ChainID(ctx)
		if err != nil {
			fmt.Printf("🌱 Node%d is starting up (try %d of %d) - waiting for chain ID...\n", nodeNumber, i+1, maxAttempts)
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Printf("🌱 Node%d RPC endpoint is healthy (chain ID: %s)\n", nodeNumber, chainID)
		return nil
	}

	return fmt.Errorf("node%d failed health check after 20 attempts", nodeNumber)
}

func main() {
	chainIDBytes, err := os.ReadFile("data/chain.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	var lastError error
	for nodeNumber := 0; nodeNumber < lib.VALIDATORS_COUNT; nodeNumber++ {
		configBytes, err := os.ReadFile(filepath.Join("data", "configs", fmt.Sprintf("config-node%d.json", nodeNumber)))
		if err != nil {
			log.Fatalf("❌ Failed to read config file: %s\n", err)
		}

		nodeConfig := lib.NodeConfig{}
		err = json.Unmarshal(configBytes, &nodeConfig)
		if err != nil {
			log.Fatalf("❌ Failed to unmarshal config: %s\n", err)
		}

		rpcURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", nodeConfig.PublicIP, nodeConfig.HTTPPort, chainID)
		fmt.Printf("Checking RPC endpoint for node%d: %s\n", nodeNumber, rpcURL)

		if err := checkNodeHealth(nodeNumber, rpcURL); err != nil {
			lastError = err
			fmt.Printf("❌ %v\n", err)
		}
	}

	if lastError != nil {
		log.Fatalf("❌ Some nodes are unhealthy")
	}

	fmt.Println("✅ All nodes are healthy!")
}
