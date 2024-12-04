package main

import (
	"context"
	"fmt"
	"log"
	"mypkg/config"
	"mypkg/helpers"

	"github.com/ava-labs/subnet-evm/interfaces"
	poavalidatormanager "github.com/ava-labs/teleporter/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	ethClient, _, err := helpers.GetLocalEthClient()
	if err != nil {
		log.Fatalf("failed to connect to client: %s\n", err)
	}

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		log.Fatalf("failed to deploy contract: %s\n", err)
	}

	// Get all logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{managerAddress},
	}

	logs, err := ethClient.FilterLogs(context.Background(), (interfaces.FilterQuery)(query))
	if err != nil {
		log.Fatal(err)
	}

	// Print all logs
	for _, vLog := range logs {
		fmt.Println("------------------------")

		fmt.Printf("Log TxHash: %s\n", vLog.TxHash.Hex())

		if event, err := contract.PoAValidatorManagerFilterer.ParseInitialValidatorCreated(vLog); err == nil {
			fmt.Printf("InitialValidatorCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodCreated(vLog); err == nil {
			fmt.Printf("ValidationPeriodCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  RegisterValidationMessageID: %x\n", event.RegisterValidationMessageID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  RegistrationExpiry: %d\n", event.RegistrationExpiry)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodEnded(vLog); err == nil {
			fmt.Printf("ValidationPeriodEnded:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Status: %d\n", event.Status)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodRegistered(vLog); err == nil {
			fmt.Printf("ValidationPeriodRegistered:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  Timestamp: %s\n", event.Timestamp.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorRemovalInitialized(vLog); err == nil {
			fmt.Printf("ValidatorRemovalInitialized:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  EndTime: %s\n", event.EndTime.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorWeightUpdate(vLog); err == nil {
			fmt.Printf("ValidatorWeightUpdate:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Nonce: %d\n", event.Nonce)
			fmt.Printf("  ValidatorWeight: %d\n", event.ValidatorWeight)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseInitialized(vLog); err == nil {
			fmt.Printf("Initialized:\n")
			fmt.Printf("  Version: %d\n", event.Version)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseOwnershipTransferred(vLog); err == nil {
			fmt.Printf("OwnershipTransferred:\n")
			fmt.Printf("  Previous Owner: %s\n", event.PreviousOwner.Hex())
			fmt.Printf("  New Owner: %s\n", event.NewOwner.Hex())
			continue
		}

		log.Printf("Failed to parse log: unknown event type\n")
	}
}
