package main

import (
	"log"
	"os"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("❌ Failed to execute command: %w\n", err)
	} else {
		if len(os.Args) > 1 {
			log.Printf("✅ Successfully executed command: %s\n", os.Args[1:])
		}
	}
}
