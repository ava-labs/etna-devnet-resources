package l1

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

type CreateSubnetParams struct {
	PrivateKey *secp256k1.PrivateKey
	RpcURL     string
}

func CreateSubnet(params CreateSubnetParams) (ids.ID, error) {
	kc := secp256k1fx.NewKeychain(params.PrivateKey)
	subnetOwner := params.PrivateKey.Address()

	// Create context with 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          params.RpcURL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to initialize wallet: %s\n", err)
	}

	owner := &secp256k1fx.OutputOwners{
		Locktime:  0,
		Threshold: 1,
		Addrs:     []ids.ShortID{subnetOwner},
	}

	createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to issue create subnet transaction: %s\n", err)
	}

	return createSubnetTx.ID(), nil
}
