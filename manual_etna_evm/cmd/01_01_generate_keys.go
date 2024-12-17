package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/spf13/cobra"
)

var GenerateKeysCmd = &cobra.Command{
	Use:   "generate-keys",
	Short: "Generate keys for the L1",
	Long:  `Generate keys for the L1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🔑 Generating keys")
		return GenerateCredsIfNotExists(helpers.Node0KeysFolder)
	},
}

func PrintHeader(header string) {
	log.Printf("\n\n%s\n\n", header)
}

func GenerateCredsIfNotExists(folder string) error {
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	if err := generateStakerKey(folder); err != nil {
		return fmt.Errorf("❌ Failed to generate staker key: %s\n", err)
	}

	if err := generateSignerKey(folder); err != nil {
		return fmt.Errorf("❌ Failed to generate signer key: %s\n", err)
	}

	return nil
}

func generateStakerKey(folder string) error {
	exists, err := helpers.FileExists(folder + "staker.key")
	if err != nil {
		return fmt.Errorf("failed to check if staker key exists: %w", err)
	}
	if exists {
		log.Println("Staker key already exists, skipping...")
		return nil
	}

	cert, key, err := staking.NewCertAndKeyBytes()
	if err != nil {
		return fmt.Errorf("failed to generate staker key: %s\n", err)
	}

	err = helpers.SaveBytes(folder+"staker.key", key)
	if err != nil {
		return fmt.Errorf("failed to save staker key: %s\n", err)
	}

	err = helpers.SaveBytes(folder+"staker.crt", cert)
	if err != nil {
		return fmt.Errorf("failed to save staker cert: %s\n", err)
	}

	return nil
}

func generateSignerKey(folder string) error {
	exists, err := helpers.FileExists(folder + "signer.key")
	if err != nil {
		return fmt.Errorf("failed to check if signer key exists: %w", err)
	}
	if exists {
		log.Println("Signer key already exists, skipping...")
		return nil
	}

	key, err := bls.NewSecretKey()
	if err != nil {
		return fmt.Errorf("failed to generate secret key: %s\n", err)
	}

	err = helpers.SaveBLSKey(folder+"signer.key", key)
	if err != nil {
		return fmt.Errorf("failed to save signer key: %s\n", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(GenerateKeysCmd)
}
