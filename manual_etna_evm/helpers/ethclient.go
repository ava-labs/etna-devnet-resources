package helpers

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/subnet-evm/ethclient"
)

func GetLocalEthClient() (ethclient.Client, *big.Int, error) {
	const maxAttempts = 100
	L1ChainId, err := LoadId("chain")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load chain ID: %s", err)
	}

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", L1ChainId)

	var client ethclient.Client
	var evmChainId *big.Int
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			fmt.Printf("Attempt %d/%d to connect to node (will sleep for %d seconds before retry)\n",
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
				fmt.Printf("Failed to get chain ID: %s (will sleep for %d seconds before retry)\n",
					err, i)
			}
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		return client, evmChainId, nil
	}

	return nil, nil, fmt.Errorf("failed after %d attempts with error: %w", maxAttempts, lastErr)
}
