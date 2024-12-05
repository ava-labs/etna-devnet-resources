package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/config"
	"mypkg/helpers"

	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
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
	blsInfoJSON, err := helpers.LoadText("add_validator_bls_json")
	if err != nil {
		return fmt.Errorf("loading add_validator_bls_json: %w", err)
	}

	blsInfo := signer.ProofOfPossession{}
	err = blsInfo.UnmarshalJSON([]byte(blsInfoJSON))
	if err != nil {
		return fmt.Errorf("unmarshaling BLS info: %w", err)
	}

	warpMessageBytes, err := helpers.LoadHex("add_validator_warp_message")
	if err != nil {
		return fmt.Errorf("loading add_validator_warp_message: %w", err)
	}

	balance := 1 * units.Avax

	key, err := helpers.LoadValidatorManagerKey()
	if err != nil {
		return fmt.Errorf("failed to load key from file: %s", err)
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
		balance,
		blsInfo.ProofOfPossession,
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
