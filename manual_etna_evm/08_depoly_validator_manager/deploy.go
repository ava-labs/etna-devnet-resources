package main

import (
	"fmt"
	"log"
	"mypkg/07_compile_validator_manager/bindings/povalidatormanager"
	"mypkg/helpers"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func main() {
	exists, err := helpers.TextFileExists("validator_manager_address")
	if err != nil {
		log.Fatalf("failed to check if validator manager address file exists: %s\n", err)
	}
	if exists {
		content, err := helpers.LoadText("validator_manager_address")
		if err != nil {
			log.Fatalf("failed to load validator manager address: %s\n", err)
		}
		log.Printf("✅ Validator manager already deployed at: %s\n", content)
		return
	}

	key, err := helpers.LoadValidatorManagerKeyECDSA()
	if err != nil {
		log.Fatalf("failed to load key from file: %s\n", err)
	}

	ethClient, evmChainId, err := helpers.GetLocalEthClient()
	if err != nil {
		log.Fatalf("failed to connect to client: %s\n", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(key, evmChainId)
	if err != nil {
		log.Fatalf("failed to create transactor: %s\n", err)
	}
	opts.GasLimit = 8000000 // Set a reasonable gas limit
	opts.GasPrice = nil     // Let the network determine the gas price

	addr, _, _, err := povalidatormanager.DeployPoAValidatorManager(opts, ethClient, uint8(0))
	if err != nil {
		log.Fatalf("failed to deploy contract: %s\n", err)
	}

	fmt.Printf("✅ Contract deployed at: %s\n", addr.Hex())

	if err := helpers.SaveText("validator_manager_address", addr.Hex()); err != nil {
		log.Fatalf("failed to save validator manager address: %s\n", err)
	}
}
