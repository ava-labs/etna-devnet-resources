// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
)

func main() {
	chainIDFilePath := filepath.Join("data", "chain.txt")
	chainIDBytes, err := os.ReadFile(chainIDFilePath)
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	key, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
	}
	kc := secp256k1fx.NewKeychain(key)

	subnetIDBytes, err := os.ReadFile("data/subnet.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}
	subnetID := ids.FromStringOrPanic(string(subnetIDBytes))
}
