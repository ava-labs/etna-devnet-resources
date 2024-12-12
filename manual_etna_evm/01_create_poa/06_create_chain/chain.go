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
	exists := helpers.FileExists(helpers.ChainIdPath)
	if exists {
		log.Println("✅ Chain already exists, exiting")
		return
	}

	key := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	kc := secp256k1fx.NewKeychain(key)

	subnetID := helpers.LoadId(helpers.SubnetIdPath)

	log.Printf("Using vmID: %s\n", constants.SubnetEVMID)

	genesisString := helpers.LoadText(helpers.L1GenesisPath)

	ctx := context.Background()

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting and registers [subnetID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, config.RPC_URL, kc, kc, primary.WalletConfig{
		SubnetIDs: []ids.ID{subnetID},
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
	helpers.SaveId(helpers.ChainIdPath, createChainTx.ID())

	log.Println("✅ Saved chain ID to file")
}
