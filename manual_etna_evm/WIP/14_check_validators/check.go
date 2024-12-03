package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"mypkg/helpers"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	err := checkValidators()
	if err != nil {
		log.Fatalf("❌ Failed to check validators: %s\n", err)
	}
}

func checkValidators() error {
	validatorManagerAddresshex, err := helpers.LoadText("validator_manager_address")
	if err != nil {
		log.Fatalf("failed to load validator manager address: %s\n", err)
	}
	validatorManagerAddress := common.HexToAddress(validatorManagerAddresshex)

	ethClient, _, err := helpers.GetLocalEthClient()
	if err != nil {
		return fmt.Errorf("failed to get local eth client: %w", err)
	}

	// contract, err := povalidatormanager.NewPoAValidatorManager(validatorManagerAddress, ethClient)
	// if err != nil {
	// 	log.Fatalf("failed to create contract instance: %s\n", err)
	// }

	// Get the latest block number
	latestBlock, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	fmt.Printf("\nChecking transactions for address: %s\n", validatorManagerAddress.Hex())

	// Iterate through all blocks
	for blockNum := uint64(0); blockNum <= latestBlock; blockNum++ {
		block, err := ethClient.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNum))
		if err != nil {
			return fmt.Errorf("failed to get block %d: %w", blockNum, err)
		}

		for _, tx := range block.Transactions() {
			if tx.To() != nil && *tx.To() == validatorManagerAddress {
				receipt, err := ethClient.TransactionReceipt(context.Background(), tx.Hash())
				if err != nil {
					return fmt.Errorf("failed to get receipt for tx %s: %w", tx.Hash().Hex(), err)
				}

				fmt.Printf("\nTransaction Hash: %s\n", tx.Hash().Hex())
				fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
				fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
				fmt.Printf("Status: %d\n", receipt.Status)
				fmt.Printf("Input Data: %x\n", tx.Data())

				// Print logs with more detail
				for i, log := range receipt.Logs {
					fmt.Printf("\nLog #%d:\n", i)
					fmt.Printf("  Address: %s\n", log.Address.Hex())
					fmt.Printf("  Block Number: %d\n", log.BlockNumber)
					fmt.Printf("  Transaction Index: %d\n", log.TxIndex)
					fmt.Printf("  Log Index: %d\n", log.Index)
					fmt.Printf("  Removed: %v\n", log.Removed)
					for j, topic := range log.Topics {
						fmt.Printf("  Topic[%d]: %s\n", j, topic.Hex())
					}
					fmt.Printf("  Data: %x\n", log.Data)
				}
				fmt.Println("------------------------")
			}
		}
	}

	return nil
}
