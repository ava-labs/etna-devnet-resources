package credshelper

import (
	"fmt"
	"log"

	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func GenerateCredsIfNotExists(folder string) error {
	if err := generateStakerKey(folder); err != nil {
		return err
	}

	if err := generateSignerKey(folder); err != nil {
		return err
	}

	return nil
}

func generateStakerKey(folder string) error {
	exists, err := helpers.FileExists(folder + "staker.key")
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

	helpers.SaveBytes(folder+"staker.key", key)
	helpers.SaveBytes(folder+"staker.crt", cert)

	return nil
}

func generateSignerKey(folder string) error {
	exists, err := helpers.FileExists(folder + "signer.key")
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

	helpers.SaveBLSKey(folder+"signer.key", key)

	return nil
}
