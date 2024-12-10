// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	exists, err := helpers.FileExists(helpers.ChainIdPath)
	if err != nil {
		log.Fatalf("❌ Failed to check if chain ID exists: %s\n", err)
	}
	if exists {
		log.Println("✅ Chain already exists, exiting")
		return
	}

	key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
	}
	kc := secp256k1fx.NewKeychain(key)

	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}

	log.Printf("Using vmID: %s\n", constants.SubnetEVMID)

	genesisString, err := helpers.LoadText(helpers.L1GenesisPath)
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
		[]byte(genesisString),
		constants.SubnetEVMID,
		nil,
		"My L1",
	)
	if err != nil {
		log.Fatalf("❌ Failed to issue create chain transaction: %s\n", err)
	}
	log.Printf("✅ Created new chain %s in %s\n", createChainTx.ID(), time.Since(createChainStartTime))

	// Save the chain ID to file
	err = helpers.SaveId(helpers.ChainIdPath, createChainTx.ID())
	if err != nil {
		log.Printf("❌ Failed to save chain ID to file: %s\n", err)
	}

	log.Println("✅ Saved chain ID to file")
}