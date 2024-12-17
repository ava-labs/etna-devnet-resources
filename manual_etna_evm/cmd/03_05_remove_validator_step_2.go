package cmd

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func SetL1ValidatorWeight(
	message *warp.Message,
) (ids.ID, *txs.Tx, error) {
	key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return ids.Empty, nil, fmt.Errorf("failed to load manager key: %w", err)
	}

	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          config.RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return ids.Empty, nil, fmt.Errorf("failed to initialize wallet: %s", err)
	}

	unsignedTx, err := wallet.P().Builder().NewSetL1ValidatorWeightTx(
		message.Bytes(),
	)
	if err != nil {
		return ids.Empty, nil, fmt.Errorf("error building tx: %w", err)
	}

	tx := txs.Tx{Unsigned: unsignedTx}
	if err := wallet.P().Signer().Sign(context.Background(), &tx); err != nil {
		return ids.Empty, nil, fmt.Errorf("error signing tx: %w", err)
	}

	err = wallet.P().IssueTx(&tx)
	if err != nil {
		return ids.Empty, nil, fmt.Errorf("error issuing tx: %w", err)
	}

	return tx.ID(), &tx, nil
}
