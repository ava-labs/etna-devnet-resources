// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"errors"
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	subnetIDFilePath := filepath.Join("data", "subnet.txt")

	// Try to load existing subnet ID first
	subnetIDBytes, err := os.ReadFile(subnetIDFilePath)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Subnet ID file does not exist, let's create a new subnet")
	} else if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	} else {
		log.Printf("📝 Existing subnet ID found: %s\n", ids.FromStringOrPanic(string(subnetIDBytes)))
		return
	}

	// If we get here, we need to create a new subnet
	key, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
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
		URI:          lib.RPC_URL,
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
	err = os.WriteFile(subnetIDFilePath, []byte(createSubnetTx.ID().String()), 0644)
	if err != nil {
		log.Printf("❌ Failed to save subnet ID to file: %s\n", err)
	}
}
