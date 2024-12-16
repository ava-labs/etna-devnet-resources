package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"

	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func init() {
	rootCmd.AddCommand(deployValidatorManagerCmd)
}

var deployValidatorManagerCmd = &cobra.Command{
	Use:   "deploy-validator-manager",
	Short: "Deploy the validator manager contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🚀 Deploying validator manager")
		privKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
		privKeyECDSA, err := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}

		ethClient, evmChainId, err := GetLocalEthClient("9650")
		if err != nil {
			return fmt.Errorf("failed to connect to client: %w", err)
		}

		myEthAddr := evm.PublicKeyToEthAddress(privKey.PublicKey())
		expectedContractAddress := MustDeriveContractAddress(myEthAddr, 1)

		deployedBytecode, err := ethClient.CodeAt(context.Background(), expectedContractAddress, nil)
		if err != nil {
			return fmt.Errorf("failed to get deployed bytecode: %w", err)
		}

		if len(deployedBytecode) > 0 {
			log.Printf("Validator manager already deployed at: %s\n", expectedContractAddress)
			return nil
		}

		opts, err := bind.NewKeyedTransactorWithChainID(privKeyECDSA, evmChainId)
		if err != nil {
			return fmt.Errorf("failed to create transactor: %w", err)
		}
		opts.GasLimit = 8000000
		opts.GasPrice = nil

		newContractAddress, tx, _, err := poavalidatormanager.DeployPoAValidatorManager(opts, ethClient, 0)
		if err != nil {
			return fmt.Errorf("failed to create contract instance: %w", err)
		}

		if newContractAddress != expectedContractAddress {
			return fmt.Errorf("expected contract address %s, got %s", expectedContractAddress, newContractAddress)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		_, err = bind.WaitMined(ctx, ethClient, tx)
		if err != nil {
			return fmt.Errorf("failed to wait for transaction confirmation: %w", err)
		}

		fmt.Printf("Validator manager deployed at: %s\n", tx.Hash().Hex())

		log.Println("Validator manager deployed")
		return nil
	},
}

func MustDeriveContractAddress(from common.Address, nonce uint64) common.Address {
	encoded, err := rlp.EncodeToBytes([]interface{}{from, nonce})
	if err != nil {
		panic(err)
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToAddress(hash[12:])
}
