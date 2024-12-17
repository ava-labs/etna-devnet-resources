package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/spf13/cobra"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func init() {
	rootCmd.AddCommand(validatorManagerInitCmd)

	validatorManagerInitCmd.Flags().StringVar(&validatorType, "validator-type", "", "Type of validator manager to deploy (poa or pos-native)")
	validatorManagerInitCmd.MarkFlagRequired("validator-type")
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

		var receipt *types.Receipt
		var tx *types.Transaction

		if validatorType == "poa" {
			receipt, tx, err = initializeValidatorManagerPoA(validatorType, managerAddress, ethClient, subnetID, opts, ecdsaKey.PublicKey)
			if err != nil {
				return fmt.Errorf("failed to initialize validator manager: %w", err)
			}
		} else if validatorType == "pos-native-token-staking" {
			return fmt.Errorf("not implemented yet")
		}

		PrintLogs(receipt.Logs)

		fmt.Printf("Validator manager initialized at: %s\n", tx.Hash().Hex())

		return nil
	},
}

func initializeValidatorManagerPoA(validatorManagerType string, managerAddress common.Address, ethClient ethclient.Client, subnetID ids.ID, opts *bind.TransactOpts, ecdsaPubKey ecdsa.PublicKey) (*types.Receipt, *types.Transaction, error) {
	logs, err := ethClient.FilterLogs(context.Background(), interfaces.FilterQuery{
		Addresses: []common.Address{managerAddress},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get contract logs: %w", err)
	}

	// Replace sleep with transaction wait
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create contract instance: %w", err)
	}
	for _, vLog := range logs {
		if _, err := contract.ParseInitialized(vLog); err == nil {
			log.Printf("Validator manager was already initialized")
			PrintLogs([]*types.Log{&vLog})
			return nil, nil, nil
		}
	}
	log.Println("Validator manager was not initialized, initializing...")

	tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
		L1ID:                   subnetID,
		ChurnPeriodSeconds:     0,
		MaximumChurnPercentage: 20,
	}, crypto.PubkeyToAddress(ecdsaPubKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize validator manager: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}

	return receipt, tx, nil
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
