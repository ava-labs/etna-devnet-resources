package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func main() {
	ecdsaKey := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)

	isInitialized := helpers.FileExists(helpers.IsValidatorManagerInitializedPath)
	if isInitialized {
		log.Println("✅ Validator manager was already initialized")
		return
	}

	subnetID := helpers.LoadId(helpers.SubnetIdPath)

	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	key := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)

	ethClient, evmChainId, err := helpers.GetLocalEthClient()
	if err != nil {
		log.Fatalf("failed to connect to client: %s\n", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(key, evmChainId)
	if err != nil {
		log.Fatalf("failed to create transactor: %s\n", err)
	}
	opts.GasLimit = 8000000
	opts.GasPrice = nil

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		log.Fatalf("failed to deploy contract: %s\n", err)
	}

	tx, err := contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
		L1ID:                   subnetID,
		ChurnPeriodSeconds:     0,
		MaximumChurnPercentage: 20,
	}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
	if err != nil {
		log.Fatalf("failed to initialize validator manager: %s\n", err)
	}

	helpers.SaveText(helpers.IsValidatorManagerInitializedPath, "true")

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
