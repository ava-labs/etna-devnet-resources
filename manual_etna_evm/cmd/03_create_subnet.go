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
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func init() {
	rootCmd.AddCommand(CreateSubnetCmd)
}

var CreateSubnetCmd = &cobra.Command{
	Use:   "create-subnet",
	Short: "Create a subnet",
	Long:  `Create a subnet`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🧱 Creating subnet")

		exists, err := helpers.FileExists(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to check if subnet ID file exists: %w", err)
		}
		if exists {
			log.Println("Subnet already exists, exiting")
			return nil
		}

		// If we get here, we need to create a new subnet
		key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load validator manager owner key: %w", err)
		}

		kc := secp256k1fx.NewKeychain(key)
		subnetOwner := key.Address()

		ctx := context.Background()

		// MakeWallet fetches the available UTXOs owned by [kc] on the network that
		// [uri] is hosting.
		walletSyncStartTime := time.Now()
		wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
			URI:          config.RPC_URL,
			AVAXKeychain: kc,
			EthKeychain:  kc,
		})
		if err != nil {
			log.Fatalf("❌ Failed to initialize wallet: %s\n", err)
		}
		log.Printf("Synced wallet in %s\n", time.Since(walletSyncStartTime))

		// Pull out useful constants to use when issuing transactions.
		owner := &secp256k1fx.OutputOwners{
			Locktime:  0,
			Threshold: 1,
			Addrs:     []ids.ShortID{subnetOwner},
		}

		createSubnetStartTime := time.Now()
		createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
		if err != nil {
			log.Fatalf("❌ Failed to issue create subnet transaction: %s\n", err)
		}
		log.Printf("✅ Created new subnet %s in %s\n", createSubnetTx.ID(), time.Since(createSubnetStartTime))

		// Save the subnet ID to file
		err = helpers.SaveId(helpers.SubnetIdPath, createSubnetTx.ID())
		if err != nil {
			return fmt.Errorf("failed to save subnet ID: %w", err)
		}
		return nil
	},
}
