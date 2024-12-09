package main

import (
	"fmt"
	"log"

	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func main() {
	if err := generateManagerKey(); err != nil {
		log.Fatalf("failed to generate manager keys: %s\n", err)
	}

	if err := generateStakerKey(); err != nil {
		log.Fatalf("failed to generate staker keys: %s\n", err)
	}

	if err := generateSignerKey(); err != nil {
		log.Fatalf("failed to generate signer keys: %s\n", err)
	}
}

func generateManagerKey() error {
	exists, err := helpers.FileExists(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %s\n", err)
	}

	if exists {
		log.Println("Manager key already exists, skipping...")
		return nil
	}

	key, err := secp256k1.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s\n", err)
	}

	if err := helpers.SaveSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath, key); err != nil {
		return fmt.Errorf("failed to save private key: %s\n", err)
	}

	return nil
}

func generateStakerKey() error {
	exists, err := helpers.FileExists(helpers.Node0StakerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %s\n", err)
	}

	if exists {
		log.Println("Staker key already exists, skipping...")
		return nil
	}

	cert, key, err := staking.NewCertAndKeyBytes()
	if err != nil {
		return fmt.Errorf("failed to generate staker key: %s\n", err)
	}

	helpers.SaveBytes(helpers.Node0StakerKeyPath, key)
	helpers.SaveBytes(helpers.Node0StakerCertPath, cert)

	return nil
}

func generateSignerKey() error {
	exists, err := helpers.FileExists(helpers.Node0SignerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %s\n", err)
	}

	if exists {
		log.Println("Signer key already exists, skipping...")
		return nil
	}

	key, err := bls.NewSecretKey()
	if err != nil {
		return fmt.Errorf("failed to generate secret key: %s\n", err)
	}

	helpers.SaveBLSKey(helpers.Node0SignerKeyPath, key)

	return nil
}
