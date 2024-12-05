package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/helpers"

	poavalidatormanager "github.com/ava-labs/teleporter/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	err := checkValidators()
	if err != nil {
		log.Fatalf("❌ Failed to check validators: %s\n", err)
	}
}

func checkValidators() error {
	ethClient, _, err := helpers.GetLocalEthClient()
	if err != nil {
		return fmt.Errorf("failed to get eth client: %w", err)
	}

	// Add network check
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}
	log.Printf("Connected to network with Chain ID: %s", chainID.String())

	initializeValidatorSetTxHashText, err := helpers.LoadText("initialize_validator_set_tx")
	if err != nil {
		return fmt.Errorf("failed to load initialize validator set tx: %w", err)
	}
	log.Printf("Loaded transaction hash: %s", initializeValidatorSetTxHashText)

	initializeValidatorSetTxHash := common.HexToHash(initializeValidatorSetTxHashText)

	// Get transaction details
	tx, _, err := ethClient.TransactionByHash(context.Background(), initializeValidatorSetTxHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Print receiving contract address (To address from transaction)
	log.Printf("Receiving Contract Address: %s\n", tx.To().Hex())

	validatorManagerCaller, err := poavalidatormanager.NewPoAValidatorManagerCaller(*tx.To(), ethClient)
	if err != nil {
		return fmt.Errorf("failed to create validator manager caller: %w", err)
	}

	registeredValidators, err := validatorManagerCaller.RegisteredValidators(nil, []byte{})
	if err != nil {
		return fmt.Errorf("failed to get registered validators: %w", err)
	}

	log.Printf("registeredValidators: %v\n", registeredValidators)

	return nil
}
