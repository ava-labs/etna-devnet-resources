package main

import (
	"log"
	"mypkg/lib"
	"os"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

func main() {

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data directory: %s\n", err)
	}

	// Generate POA validator manager key if it doesn't exist
	if _, err := os.Stat(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH); err == nil {
		log.Println("POA validator manager keys were previously generated in ./data/ folder")
	} else {
		key, err := secp256k1.NewPrivateKey()
		if err != nil {
			log.Fatalf("failed to generate POA validator manager private key: %s\n", err)
		}

		err = lib.SaveKeyToFile(key, lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
		if err != nil {
			log.Fatalf("failed to save POA validator manager key to file: %s\n", err)
		}
		log.Println("✅ POA validator manager keys generated and saved in ./data/ folder")
	}

	// Generate teleporter deployer key if it doesn't exist
	if _, err := os.Stat(lib.TELEPORTER_DEPLOYER_KEY_PATH); err == nil {
		log.Println("Teleporter deployer keys were previously generated in ./data/ folder")
	} else {
		key, err := secp256k1.NewPrivateKey()
		if err != nil {
			log.Fatalf("failed to generate teleporter deployer private key: %s\n", err)
		}

		err = lib.SaveKeyToFile(key, lib.TELEPORTER_DEPLOYER_KEY_PATH)
		if err != nil {
			log.Fatalf("failed to save teleporter deployer key to file: %s\n", err)
		}
		log.Println("✅ Teleporter deployer keys generated and saved in ./data/ folder")
	}
}
