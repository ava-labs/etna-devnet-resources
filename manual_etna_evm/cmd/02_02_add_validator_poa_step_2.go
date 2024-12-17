package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func RegisterL1ValidatorOnPChain(warpMessage *warp.Message, credsFolder string) error {
	_, proofOfPossession, err := NodeInfoFromCreds(credsFolder)
	if err != nil {
		return fmt.Errorf("failed to get node info from creds: %w", err)
	}

	key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load manager key: %w", err)
	}

	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          config.RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}

	unsignedTx, err := wallet.P().Builder().NewRegisterL1ValidatorTx(
		1*units.Avax,
		proofOfPossession.ProofOfPossession,
		warpMessage.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("error building tx: %w", err)
	}

	tx := txs.Tx{Unsigned: unsignedTx}
	if err := wallet.P().Signer().Sign(context.Background(), &tx); err != nil {
		return fmt.Errorf("error signing tx: %w", err)
	}

	err = wallet.P().IssueTx(&tx)
	if err != nil {
		return fmt.Errorf("error issuing tx: %w", err)
	}

	return nil
}
