package main

import (
	"fmt"
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/ava-labs/avalanchego/ids"
)

//go:embed chain.json
var defaultEVMSubnetChainJson []byte

func main() {
	if err := lib.FillNodeConfigs(""); err != nil {
		log.Fatalf("❌ Failed to fill node configs: %s", err)
	}
	fmt.Println("✅ Successfully created configs")

	chainIDFilePath := filepath.Join("data", "chain.txt")
	chainIDBytes, err := os.ReadFile(chainIDFilePath)
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}

	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	if err := os.MkdirAll(fmt.Sprintf("data/chains/%s", chainID), 0755); err != nil {
		log.Fatalf("❌ Failed to create chains directory: %s\n", err)
	}

	if err := os.WriteFile(fmt.Sprintf("data/chains/%s/config.json", chainID), defaultEVMSubnetChainJson, 0644); err != nil {
		log.Fatalf("❌ Failed to write chain.json: %s\n", err)
	}
}
