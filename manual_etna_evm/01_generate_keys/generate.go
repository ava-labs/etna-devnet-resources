package main

import (
	"log"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"
)

func main() {
	if err := credshelper.GenerateCredsIfNotExists(helpers.Node0KeysFolder); err != nil {
		log.Fatalf("Failed to generate credentials: %s\n", err)
	}
}
