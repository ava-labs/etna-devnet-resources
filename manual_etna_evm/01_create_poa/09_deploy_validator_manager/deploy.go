package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"

	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func main() {
	if err := deployValidatorManager(); err != nil {
		log.Fatalf("failed to deploy validator manager: %s\n", err)
	}
}

func deployValidatorManager() error {
	privKey := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	privKeyECDSA := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)

	ethClient, evmChainId, err := helpers.GetLocalEthClient("9650")
	if err != nil {
		return fmt.Errorf("failed to connect to client: %w", err)
	}

	myEthAddr := evm.PublicKeyToEthAddress(privKey.PublicKey())
	expectedContractAddress := helpers.DeriveContractAddress(myEthAddr, 1)

	for i := 0; i <= 10; i++ {
		fmt.Printf("DEBUG: Expected contract address for nonce %d: %s\n", i, helpers.DeriveContractAddress(myEthAddr, uint64(i)))
	}

	deployedBytecode, err := ethClient.CodeAt(context.Background(), expectedContractAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get deployed bytecode: %w", err)
	}

	if len(deployedBytecode) > 0 {
		fmt.Printf("✅ Validator manager already deployed at: %s\n", expectedContractAddress)
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

	fmt.Printf("✅ Validator manager deployed at: %s\n", tx.Hash().Hex())

	log.Println("✅ Validator manager deployed")
	return nil
}
