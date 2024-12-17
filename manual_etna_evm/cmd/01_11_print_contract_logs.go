package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printContractLogsCmd)
}

var printContractLogsCmd = &cobra.Command{
	Use:   "print-contract-logs",
	Short: "Print contract logs",
	RunE: func(cmd *cobra.Command, args []string) error {

		var port string
		if len(args) >= 1 {
			port = args[0]
		} else {
			port = "9650"
		}

		PrintHeader(fmt.Sprintf("🧱 Printing contract logs from localhost:%s", port))

		if err := printEVMContractLogs(port); err != nil {
			return fmt.Errorf("failed to print EVM contract logs: %w", err)
		}

		return nil
	},
}

func printEVMContractLogs(port string) error {
	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	ethClient, _, err := GetLocalEthClient(port)
	if err != nil {
		return fmt.Errorf("failed to connect to client: %s\n", err)
	}

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %s\n", err)
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

		// Try parsing each event type
		if event, err := contract.PoAValidatorManagerFilterer.ParseInitialValidatorCreated(vLog); err == nil {
			fmt.Printf("InitialValidatorCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  Weight: %d\n", event.Weight)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodCreated(vLog); err == nil {
			fmt.Printf("ValidationPeriodCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  RegisterValidationMessageID: %x\n", event.RegisterValidationMessageID)
			fmt.Printf("  Weight: %d\n", event.Weight)
			fmt.Printf("  RegistrationExpiry: %d\n", event.RegistrationExpiry)
			continue
		}

		// Add these new event parsers
		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorWeightUpdate(vLog); err == nil {
			fmt.Printf("ValidatorWeightUpdate:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Nonce: %d\n", event.Nonce)
			fmt.Printf("  Weight: %d\n", event.Weight)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodRegistered(vLog); err == nil {
			fmt.Printf("ValidationPeriodRegistered:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Weight: %d\n", event.Weight)
			fmt.Printf("  Timestamp: %d\n", event.Timestamp)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodEnded(vLog); err == nil {
			fmt.Printf("ValidationPeriodEnded:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Status: %d\n", event.Status)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorRemovalInitialized(vLog); err == nil {
			fmt.Printf("ValidatorRemovalInitialized:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			fmt.Printf("  Weight: %d\n", event.Weight)
			fmt.Printf("  EndTime: %d\n", event.EndTime)
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

		log.Printf("❗ Failed to parse log: unknown event type\n")
		fmt.Printf("  Address: %s\n", vLog.Address.Hex())
		fmt.Printf("  Topics: %v\n", vLog.Topics)
		fmt.Printf("  Data: %x\n", vLog.Data)
	}
	return nil
}
