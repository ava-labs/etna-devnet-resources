package main

import (
	"fmt"
	"log"
	"mypkg/lib"
	"os"

	"github.com/ava-labs/avalanchego/ids"
)

func main() {
	subnetIDBytes, err := os.ReadFile("data/subnet.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}
	subnetID := ids.FromStringOrPanic(string(subnetIDBytes))

	if err := lib.FillNodeConfigs(subnetID.String()); err != nil {
		log.Fatalf("❌ Failed to fill node configs: %s", err)
	}
	fmt.Println("✅ Successfully updated node configs")
}
