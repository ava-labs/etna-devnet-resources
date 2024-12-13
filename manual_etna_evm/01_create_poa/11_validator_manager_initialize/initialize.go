package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/core/types"

	nativestakingmanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/NativeTokenStakingManager"
	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
)

func main() {
	rewardCalculatorAddress := helpers.LoadEVMAddress(helpers.RewardCalculatorAddressPath)

	ecdsaKey := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)

	isInitialized := helpers.FileExists(helpers.IsValidatorManagerInitializedPath)
	if isInitialized {
		log.Println("✅ Validator manager was already initialized")
		return
	}

	subnetID := helpers.LoadId(helpers.SubnetIdPath)
	chainID := helpers.LoadId(helpers.ChainIdPath)

	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	key := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)

	ethClient, evmChainId, err := helpers.GetLocalEthClient("9650")
	if err != nil {
		log.Fatalf("failed to connect to client: %s\n", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(key, evmChainId)
	if err != nil {
		log.Fatalf("failed to create transactor: %s\n", err)
	}
	opts.GasLimit = 8000000
	opts.GasPrice = nil

	var tx *types.Transaction
	if helpers.GetDesiredContractName() == "PoAValidatorManager" {
		log.Printf("🔍 Initializing PoAValidatorManager\n")
		contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
		if err != nil {
			log.Fatalf("failed to deploy contract: %s\n", err)
		}

		tx, err = contract.Initialize(opts, poavalidatormanager.ValidatorManagerSettings{
			L1ID:                   subnetID,
			ChurnPeriodSeconds:     60,
			MaximumChurnPercentage: 20,
		}, crypto.PubkeyToAddress(ecdsaKey.PublicKey))
	} else if helpers.GetDesiredContractName() == "NativeTokenStakingManager" {
		log.Printf("🔍 Initializing NativeTokenStakingManager\n")
		contract, err := nativestakingmanager.NewNativeTokenStakingManager(managerAddress, ethClient)
		if err != nil {
			log.Fatalf("failed to deploy contract: %s\n", err)
		}

		tx, err = contract.Initialize(opts, nativestakingmanager.PoSValidatorManagerSettings{
			BaseSettings: nativestakingmanager.ValidatorManagerSettings{
				L1ID:                   subnetID,
				ChurnPeriodSeconds:     60,
				MaximumChurnPercentage: 20,
			},
			MinimumStakeAmount:       helpers.ApplyDefaultDenomination(1),
			MaximumStakeAmount:       helpers.ApplyDefaultDenomination(1000),
			MinimumStakeDuration:     100,
			MinimumDelegationFeeBips: 1,
			MaximumStakeMultiplier:   1,
			WeightToValueFactor:      big.NewInt(1),
			RewardCalculator:         rewardCalculatorAddress,
			UptimeBlockchainID:       chainID, //see https://github.com/ava-labs/icm-contracts/blob/87e7d53ff504c13ed702ac2fb3b34521488ebc5d/contracts/validator-manager/UptimeMessageSpec.md
		})
	}
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
