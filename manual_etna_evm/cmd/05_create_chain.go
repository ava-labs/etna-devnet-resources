package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/spf13/cobra"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func init() {
	rootCmd.AddCommand(CreateChainCmd)
}

var CreateChainCmd = &cobra.Command{
	Use:   "create-chain",
	Short: "Create a chain",
	Long:  `Create a chain`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🧱 Creating chain")

		exists, err := helpers.FileExists(helpers.ChainIdPath)
		if err != nil {
			return fmt.Errorf("failed to check if chain ID file exists: %w", err)
		}
		if exists {
			log.Println("Chain already exists, exiting")
			return nil
		}

		key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load validator manager owner key: %w", err)
		}
		kc := secp256k1fx.NewKeychain(key)

		subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to load subnet ID: %w", err)
		}

		log.Printf("Using vmID: %s\n", constants.SubnetEVMID)

		genesisString, err := helpers.LoadText(helpers.L1GenesisPath)
		if err != nil {
			return fmt.Errorf("failed to load genesis: %w", err)
		}

		ctx := context.Background()

		// MakeWallet fetches the available UTXOs owned by [kc] on the network that
		// [uri] is hosting and registers [subnetID].
		walletSyncStartTime := time.Now()
		wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
			URI:          config.RPC_URL,
			AVAXKeychain: kc,
			EthKeychain:  kc,
			SubnetIDs:    []ids.ID{subnetID},
		})
		if err != nil {
			return fmt.Errorf("failed to initialize wallet: %w", err)
		}
		log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

		// Get the P-chain wallet
		pWallet := wallet.P()

		createChainStartTime := time.Now()
		createChainTx, err := pWallet.IssueCreateChainTx(
			subnetID,
			[]byte(genesisString),
			constants.SubnetEVMID,
			nil,
			"My L1",
		)
		if err != nil {
			return fmt.Errorf("failed to issue create chain transaction: %w", err)
		}
		log.Printf("Created new chain %s in %s\n", createChainTx.ID(), time.Since(createChainStartTime))

		// Save the chain ID to file
		err = helpers.SaveId(helpers.ChainIdPath, createChainTx.ID())
		if err != nil {
			return fmt.Errorf("failed to save chain ID: %w", err)
		}

		log.Println("Saved chain ID to file")
		return nil
	},
}
