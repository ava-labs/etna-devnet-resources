package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(launchNodeCmd)
}

// original script:
// SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

// export CURRENT_UID=$(id -u)
// export CURRENT_GID=$(id -g)
// export AVALANCHEGO_TRACK_SUBNETS=$(cat "${SCRIPT_DIR}/../../data/subnet_id.txt" | tr -d '\n')

// export CHAIN_ID=$(cat "${SCRIPT_DIR}/../../data/chain_id.txt" | tr -d '\n')

// mkdir -p "${SCRIPT_DIR}/../../data/chains/${CHAIN_ID}"
// cp "${SCRIPT_DIR}/evm_debug_config.json" "${SCRIPT_DIR}/../../data/chains/${CHAIN_ID}/config.json"

// docker compose -f "${SCRIPT_DIR}/docker-compose.yml" down || true
// docker compose -f "${SCRIPT_DIR}/docker-compose.yml" up -d --build $1

var launchNodeCmd = &cobra.Command{
	Use:   "launch-node",
	Short: "Launch a node",
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🐳 Launching node")

		subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to load subnet ID: %w", err)
		}

		chainID, err := helpers.LoadId(helpers.ChainIdPath)
		if err != nil {
			return fmt.Errorf("failed to load chain ID: %w", err)
		}

		// Create chains directory if it doesn't exist
		chainsDir := fmt.Sprintf("./data/chains/%s", chainID)
		err = os.MkdirAll(chainsDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create chains directory: %w", err)
		}

		// Copy config file
		err = helpers.CopyFile("./cmd/node/evm_debug_config.json", fmt.Sprintf("%s/config.json", chainsDir))
		if err != nil {
			return fmt.Errorf("failed to copy config file: %w", err)
		}

		// Get current user and group IDs
		currentUID := exec.Command("id", "-u")
		uidOutput, err := currentUID.Output()
		if err != nil {
			return fmt.Errorf("failed to get current UID: %w", err)
		}

		currentGID := exec.Command("id", "-g")
		gidOutput, err := currentGID.Output()
		if err != nil {
			return fmt.Errorf("failed to get current GID: %w", err)
		}

		// Set environment variables
		env := []string{
			fmt.Sprintf("CURRENT_UID=%s", strings.TrimSpace(string(uidOutput))),
			fmt.Sprintf("CURRENT_GID=%s", strings.TrimSpace(string(gidOutput))),
			fmt.Sprintf("AVALANCHEGO_TRACK_SUBNETS=%s", subnetID),
		}

		// Change working directory for docker compose commands
		downCmd := exec.Command("docker", "compose", "down")
		downCmd.Dir = "./cmd/node"
		downCmd.Env = append(downCmd.Env, env...)
		output, err := downCmd.CombinedOutput()
		if err != nil {
			log.Printf("Warning: docker compose down failed: %v\n%s", err, output)
		} else {
			log.Printf("Docker compose down output:\n%s", output)
		}

		upCmd := exec.Command("docker", "compose", "up", "-d", "--build")
		upCmd.Dir = "./cmd/node"
		upCmd.Env = append(upCmd.Env, env...)
		output, err = upCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to start docker compose: %w\n%s", err, output)
		}
		log.Printf("Docker compose up output:\n%s", output)

		_, evmChainId, err := GetLocalEthClient("9650")
		if err != nil {
			return fmt.Errorf("failed to wait for chain to be available: %w", err)
		}

		fmt.Printf("✅ Subnet is healthy and responding\n")
		fmt.Printf("Chain ID (decimal): %d\n", evmChainId.Int64())
		fmt.Printf("To see logs, run: docker logs -f node0\n")

		return nil
	},
}

func GetLocalEthClient(port string) (ethclient.Client, *big.Int, error) {
	const maxAttempts = 100
	L1ChainId, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load chain ID: %w", err)
	}

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", port, L1ChainId)

	var client ethclient.Client
	var evmChainId *big.Int
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			log.Printf("Attempt %d/%d to connect to node (will sleep for %d seconds before retry)",
				i+1, maxAttempts, i)
		}

		client, err = ethclient.DialContext(context.Background(), nodeURL)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to node: %s", err)
			if i > 0 {
				fmt.Printf("Failed to connect: %s\n", err)
			}
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		evmChainId, err = client.ChainID(context.Background())
		if err != nil {
			lastErr = fmt.Errorf("failed to get chain ID: %s", err)
			if i > 0 {
				log.Printf("chain is not ready yet: %s (will sleep for %d seconds before retry)\n",
					strings.TrimSpace(string(lastErr.Error())), i)
			}
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		return client, evmChainId, nil
	}

	return nil, nil, fmt.Errorf("failed after %d attempts with error: %w", maxAttempts, lastErr)
}
