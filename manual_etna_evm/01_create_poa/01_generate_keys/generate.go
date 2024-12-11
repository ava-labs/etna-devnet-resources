package main

import (
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"
)

func main() {
	credshelper.GenerateCredsIfNotExists(helpers.Node0KeysFolder)
}
