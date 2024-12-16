package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/spf13/cobra"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func init() {
	rootCmd.AddCommand(validatorManagerInitCmd)
}

var validatorManagerInitCmd = &cobra.Command{
	Use:   "validator-manager-init",
	Short: "Initialize the validator manager contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🔌 Initializing validator manager (EVM transaction)")

		ecdsaKey, err := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}

		managerAddress := common.HexToAddress(config.ProxyContractAddress)
		ethClient, evmChainId, err := GetLocalEthClient("9650")
		if err != nil {
			return fmt.Errorf("failed to connect to client: %w", err)
		}

		// Check for Initialized event in logs
		query := interfaces.FilterQuery{
			Addresses: []common.Address{managerAddress},
		}
		logs, err := ethClient.FilterLogs(context.Background(), query)
		if err != nil {
			return fmt.Errorf("failed to get contract logs: %w", err)
		}

		contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
		if err != nil {
			return fmt.Errorf("failed to create contract instance: %w", err)
		}
		for _, vLog := range logs {
			if _, err := contract.ParseInitialized(vLog); err == nil {
				log.Printf("Validator manager was already initialized")
				PrintLogs([]*types.Log{&vLog})
				return nil
			}
		}
		log.Println("Validator manager was not initialized, initializing...")

		subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to load subnet ID: %w", err)
		}

		key, err := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}

		opts, err := bind.NewKeyedTransactorWithChainID(key, evmChainId)
		if err != nil {
			return fmt.Errorf("failed to create transactor: %w", err)
		}
		opts.GasLimit = 8000000
		opts.GasPrice = nil

		tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
			L1ID:                   subnetID,
			ChurnPeriodSeconds:     0,
			MaximumChurnPercentage: 20,
		}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
		if err != nil {
			return fmt.Errorf("failed to initialize validator manager: %w", err)
		}

		// Replace sleep with transaction wait
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		receipt, err := bind.WaitMined(ctx, ethClient, tx)
		if err != nil {
			return fmt.Errorf("failed to wait for transaction confirmation: %w", err)
		}

		PrintLogs(receipt.Logs)

		fmt.Printf("Validator manager initialized at: %s\n", tx.Hash().Hex())

		return nil
	},
}

func PrintLogs(logs []*types.Log) {
	log.Println("Transaction logs:")
	for _, logEntry := range logs {
		logJSON, err := json.MarshalIndent(logEntry, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal log to JSON: %v", err)
			continue
		}
		log.Printf("Log: %s", logJSON)
	}
}
