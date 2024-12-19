package l1

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

type CreateChainParams struct {
	PrivateKey  *secp256k1.PrivateKey
	SubnetID    ids.ID
	GenesisData string
	RpcURL      string
	ChainName   string
}

func CreateChain(params CreateChainParams) (ids.ID, error) {
	// Create keychain from private key
	kc := secp256k1fx.NewKeychain(params.PrivateKey)

	// Create context with 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Initialize wallet with subnet ID
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          params.RpcURL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
		SubnetIDs:    []ids.ID{params.SubnetID},
	})
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to initialize wallet: %s", err)
	}

	// Issue create chain transaction
	createChainTx, err := wallet.P().IssueCreateChainTx(
		params.SubnetID,
		[]byte(params.GenesisData),
		constants.SubnetEVMID,
		nil,
		params.ChainName,
	)
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to issue create chain transaction: %s", err)
	}

	return createChainTx.ID(), nil
}
