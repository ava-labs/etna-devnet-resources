package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

	rpcURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)
	fmt.Printf("Checking RPC endpoint: %s\n", rpcURL)

	const maxAttempts = 60
	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		client, err := evm.GetClient(rpcURL)
		if err != nil {
			fmt.Printf("🌱 Node is still booting up (try %d of %d) - waiting for RPC endpoint...\n", i+1, maxAttempts)

			lastLine, err := getLastLineOfDockerLogs("node1")
			if err != nil {
				log.Fatalf("❌ Failed to get last line of Docker logs: %s\n", err)
			}
			fmt.Printf("Last line of Docker logs: %s\n", lastLine)

			time.Sleep(2 * time.Second)
			continue
		}

		chainID, err := client.ChainID(ctx)
		if err != nil {
			fmt.Printf("🌱 Node is starting up (try %d of %d) - waiting for chain ID...\n", i+1, maxAttempts)
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Printf("✅ RPC endpoint is healthy (chain ID: %s)\n", chainID)
		return
	}

	log.Fatalf("❌ Node failed health check after %d attempts", maxAttempts)
}

func getLastLineOfDockerLogs(containerName string) (string, error) {
	logs, err := exec.Command("docker", "logs", containerName).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(logs), "\n")
	if len(lines) == 0 {
		return "", nil
	}
	// Get last non-empty line
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			return lines[i], nil
		}
	}
	return "", nil
}
