// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"mypkg/config"
	"mypkg/helpers"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	exists, err := helpers.IdFileExists("chain")
	if err != nil {
		log.Fatalf("❌ Failed to check if chain ID exists: %s\n", err)
	}
	if exists {
		log.Println("✅ Chain already exists, exiting")
		return
	}

	key, err := helpers.LoadValidatorManagerKey()
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
	}
	kc := secp256k1fx.NewKeychain(key)

	subnetIDBytes, err := os.ReadFile("data/subnet.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}
	subnetID := ids.FromStringOrPanic(string(subnetIDBytes))

	vmID := constants.SubnetEVMID
	name := "Step by step subnet"

	log.Printf("Using vmID: %s\n", vmID)

	genesisBytes, err := os.ReadFile("data/L1-genesis.json")
	if err != nil {
		log.Fatalf("❌ Failed to read genesis: %s\n", err)
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
		log.Fatalf("❌ Failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()

	createChainStartTime := time.Now()
	createChainTx, err := pWallet.IssueCreateChainTx(
		subnetID,
		genesisBytes,
		vmID,
		nil,
		name,
	)
	if err != nil {
		log.Fatalf("❌ Failed to issue create chain transaction: %s\n", err)
	}
	log.Printf("✅ Created new chain %s in %s\n", createChainTx.ID(), time.Since(createChainStartTime))

	// Save the chain ID to file
	err = helpers.SaveId("chain", createChainTx.ID())
	if err != nil {
		log.Printf("❌ Failed to save chain ID to file: %s\n", err)
	}

	log.Println("✅ Saved chain ID to file")
}
