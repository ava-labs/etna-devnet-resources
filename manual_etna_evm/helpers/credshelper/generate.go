package credshelper

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func GenerateCredsIfNotExists(folder string) {
	if err := generateStakerKey(folder); err != nil {
		log.Fatalf("❌ Failed to generate staker key: %s\n", err)
	}

	if err := generateSignerKey(folder); err != nil {
		log.Fatalf("❌ Failed to generate signer key: %s\n", err)
	}
}
func GetCredsBase64(folder string) (string, string, string) {
	stakerKey := helpers.LoadBytes(folder + "staker.key")
	stakerCert := helpers.LoadBytes(folder + "staker.crt")
	signerKey := helpers.LoadBytes(folder + "signer.key")

	return base64.StdEncoding.EncodeToString(stakerKey),
		base64.StdEncoding.EncodeToString(stakerCert),
		base64.StdEncoding.EncodeToString(signerKey)
}

func generateStakerKey(folder string) error {
	exists := helpers.FileExists(folder + "staker.key")
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
	exists := helpers.FileExists(folder + "signer.key")
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
