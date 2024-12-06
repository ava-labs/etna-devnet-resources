package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/helpers"
	"time"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"

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

	addr, tx, validatorManagerCaller, err := poavalidatormanager.DeployPoAValidatorManager(opts, ethClient, uint8(1))
	if err != nil {
		log.Fatalf("failed to deploy contract: %s\n", err)
	}

	fmt.Printf("✅ Contract deployed at: %s\n", addr.Hex())

	fmt.Printf("✅ Waiting for transaction %s to be mined...\n", tx.Hash().Hex())
	time.Sleep(5 * time.Second)

	receipt, err := ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatalf("failed to get transaction receipt: %s\n", err)
	}
	fmt.Printf("✅ Transaction receipt: %+v\n", receipt)

	if err := helpers.SaveText("validator_manager_address", addr.Hex()); err != nil {
		log.Fatalf("failed to save validator manager address: %s\n", err)
	}

	nodeID, _, err := helpers.GetNodeInfoRetry("http://127.0.0.1:9650")
	if err != nil {
		log.Fatalf("❌ Failed to get node info: %s\n", err)
	}

	registeredValidators, err := validatorManagerCaller.RegisteredValidators(nil, nodeID[:])
	if err != nil {
		log.Fatalf("failed to get registered validators: %s\n", err)
	}

	log.Printf("registeredValidators: %v\n", registeredValidators)

	// Get deployed bytecode from chain
	deployedCode, err := ethClient.CodeAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatalf("failed to get deployed code: %s\n", err)
	}

	fmt.Printf("\n✅ Deployed contract bytecode:\n0x%x\n", deployedCode)
}
