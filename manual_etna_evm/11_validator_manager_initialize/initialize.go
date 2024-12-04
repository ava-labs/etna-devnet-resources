package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mypkg/config"
	"mypkg/helpers"
	"time"

	poavalidatormanager "github.com/ava-labs/teleporter/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func main() {
	ecdsaKey, err := helpers.LoadValidatorManagerKeyECDSA()
	if err != nil {
		log.Fatalf("failed to load validator manager key: %w", err)
	}

	isInitialized, err := helpers.TextFileExists("validator_manager_initialized")
	if err != nil {
		log.Fatalf("failed to check if validator manager initialized file exists: %s\n", err)
	}
	if isInitialized {
		log.Println("✅ Validator manager was already initialized")
		return
	}

	subnetID, err := helpers.LoadId("subnet")
	if err != nil {
		log.Fatalf("failed to load subnet ID: %s\n", err)
	}

	managerAddress := common.HexToAddress(config.ProxyContractAddress)

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

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		log.Fatalf("failed to deploy contract: %s\n", err)
	}

	tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
		SubnetID:               subnetID,
		ChurnPeriodSeconds:     0,
		MaximumChurnPercentage: 20,
	}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
	if err != nil {
		log.Fatalf("failed to initialize validator manager: %s\n", err)
	}

	if err := helpers.SaveText("validator_manager_initialized", "true"); err != nil {
		log.Fatalf("failed to save validator manager initialized file: %s\n", err)
	}

	// Replace sleep with transaction wait
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction confirmation: %s\n", err)
	}

	receiptJSON, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal receipt to JSON: %s\n", err)
	}
	fmt.Printf("✅ Transaction receipt: %s\n", receiptJSON)

	fmt.Printf("✅ Validator manager initialized at: %s\n", tx.Hash().Hex())

	log.Println("✅ Validator manager initialized")
}
