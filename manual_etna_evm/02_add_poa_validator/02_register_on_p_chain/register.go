package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	log.Printf("Attempting to register L1 validator on P-chain...")

	for i := 0; i < 3; i++ {
		log.Printf("Attempt %d/3", i+1)
		if err := RegisterL1ValidatorOnPChain(); err != nil {
			log.Printf("Attempt %d failed: %s", i+1, err)
			if i < 2 {
				log.Printf("Retrying...")
				continue
			}
			log.Fatalf("❌ All attempts to register L1 validator failed")
		}
		log.Printf("✅ Successfully registered L1 validator on P-chain")
		break
	}
}

func RegisterL1ValidatorOnPChain() error {
	subnetID := helpers.LoadId(helpers.SubnetIdPath)

	warpMessageBytes := helpers.LoadHex(helpers.AddValidatorWarpMessagePath)

	_, proofOfPossession := credshelper.NodeInfoFromCreds(helpers.AddValidatorKeysFolder)

	key := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)

	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), config.RPC_URL, kc, kc, primary.WalletConfig{
		SubnetIDs: []ids.ID{subnetID},
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}

	unsignedTx, err := wallet.P().Builder().NewRegisterL1ValidatorTx(
		1*units.Avax,
		proofOfPossession.ProofOfPossession,
		warpMessageBytes,
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
