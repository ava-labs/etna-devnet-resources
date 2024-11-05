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

func main() {
	chainIDBytes, err := os.ReadFile("data/chain.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

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

		client, err := evm.GetClient(rpcURL)
		if err != nil {
			log.Fatalf("❌ Node%d failed to create client: %s\n", nodeNumber, err)
		}

		chainID, err := client.ChainID(ctx)
		if err != nil {
			log.Fatalf("❌ Node%d RPC endpoint error: %s\n", nodeNumber, err)
		}
		fmt.Printf("✅ Node%d RPC endpoint is healthy (chain ID: %s)\n", nodeNumber, chainID)
	}

	fmt.Println("✅ All nodes are healthy!")
}
