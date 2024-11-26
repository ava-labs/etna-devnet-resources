package main

import (
	"fmt"
	"log"
	"mypkg/helpers"
	"os"

	_ "embed"
)

//go:embed chain.json
var defaultEVMSubnetChainJson []byte

//go:embed upgrades.json
var upgradesJSON []byte

func main() {
	chainID, err := helpers.LoadId("chain")
	if err != nil {
		log.Fatalf("❌ Failed to load chain ID: %s\n", err)
	}

	if err := os.MkdirAll(fmt.Sprintf("data/chains/%s", chainID), 0755); err != nil {
		log.Fatalf("❌ Failed to create chains directory: %s\n", err)
	}

	if err := os.WriteFile(fmt.Sprintf("data/chains/%s/config.json", chainID), defaultEVMSubnetChainJson, 0644); err != nil {
		log.Fatalf("❌ Failed to write chain.json: %s\n", err)
	}

	if err := os.WriteFile("data/upgrade.json", upgradesJSON, 0644); err != nil {
		log.Fatalf("❌ Failed to write upgrades.json: %s\n", err)
	}
}
