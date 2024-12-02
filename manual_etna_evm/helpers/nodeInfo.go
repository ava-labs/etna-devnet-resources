package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
)

// Naively retries getting node info from the node until it succeeds
func GetNodeInfoRetry(endpoint string) (nodeID ids.NodeID, proofOfPossession *signer.ProofOfPossession, err error) {
	infoClient := info.NewClient(endpoint)
	fmt.Printf("Getting node info from %s\n", endpoint)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		nodeID, proofOfPossession, err = infoClient.GetNodeID(ctx)
		if err == nil {
			return
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return ids.NodeID{}, nil, fmt.Errorf("failed to get node info after 10 retries")
}
