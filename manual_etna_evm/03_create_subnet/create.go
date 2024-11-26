// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"mypkg/config"
	"mypkg/helpers"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	exists, err := helpers.IdFileExists("subnet")
	if err != nil {
		log.Fatalf("❌ Failed to check if subnet ID exists: %s\n", err)
	}
	if exists {
		log.Println("Subnet already exists, exiting")
		return
	}

	// If we get here, we need to create a new subnet
	key, err := helpers.LoadValidatorManagerKey()
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
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
	err = helpers.SaveId("subnet", createSubnetTx.ID())
	if err != nil {
		log.Printf("❌ Failed to save subnet ID to file: %s\n", err)
	}
}
