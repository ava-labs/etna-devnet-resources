package helpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/ethclient"
)

func GetLocalEthClient() (ethclient.Client, *big.Int, error) {
	L1ChainId, err := LoadId("chain")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load chain ID: %s", err)
	}

	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", L1ChainId)

	client, err := ethclient.DialContext(context.Background(), nodeURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to node: %s", err)
	}

	evmChainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain ID: %s", err)
	}

	return client, evmChainId, nil
}
