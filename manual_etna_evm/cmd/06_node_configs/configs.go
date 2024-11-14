package main

import (
	"fmt"
	"log"
	"mypkg/lib"

	_ "embed"
)

func main() {
	if err := lib.FillNodeConfigs(""); err != nil {
		log.Fatalf("❌ Failed to fill node configs: %s", err)
	}
	fmt.Println("✅ Successfully created configs")
}
