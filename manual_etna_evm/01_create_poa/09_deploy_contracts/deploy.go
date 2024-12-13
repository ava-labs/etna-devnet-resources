package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"

	examplerewardcalculator "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/ExampleRewardCalculator"
	nativestakingmanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/NativeTokenStakingManager"
	poavalidatormanager "github.com/ava-labs/icm-contracts/abi-bindings/go/validator-manager/PoAValidatorManager"
)

type ICMInitializable int

const (
	Allowed ICMInitializable = iota
	Disallowed
)

func main() {
	for i := 0; i < 3; i++ {
		err := deployContracts()
		if err == nil {
			return
		}
		log.Printf("Attempt %d failed: %v\n", i+1, err)
		if i < 2 {
			log.Println("Retrying...")
		}
	}
	log.Fatal("Failed to deploy contracts after 3 attempts")
}

func deployContracts() error {
	if helpers.FileExists(helpers.ValidatorManagerAddressPath) {
		fmt.Println("Validator manager already deployed")
		return nil
	}

	ethClient, evmChainId, err := helpers.GetLocalEthClient("9650")
	if err != nil {
		return fmt.Errorf("failed to connect to client: %s", err)
	}

	key := helpers.LoadSecp256k1PrivateKeyECDSA(helpers.ValidatorManagerOwnerKeyPath)
	opts, err := bind.NewKeyedTransactorWithChainID(key, evmChainId)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %s", err)
	}
	opts.GasLimit = 8000000
	opts.GasPrice = nil

	var tx *types.Transaction
	var validatorManagerAddress common.Address
	var rewardCalculatorAddress common.Address

	if helpers.GetDesiredContractName() == "PoAValidatorManager" {
		log.Printf("🔍 Deploying PoAValidatorManager\n")
		//FIXME: not sure it should be Allowed or Disallowed
		validatorManagerAddress, tx, _, err = poavalidatormanager.DeployPoAValidatorManager(opts, ethClient, uint8(Allowed))
	} else if helpers.GetDesiredContractName() == "NativeTokenStakingManager" {
		log.Printf("🔍 Deploying NativeTokenStakingManager\n")
		validatorManagerAddress, tx, _, err = nativestakingmanager.DeployNativeTokenStakingManager(opts, ethClient, uint8(Allowed))

		rewardBasisPoints := uint64(100)
		rewardCalculatorAddress, _, _, err = examplerewardcalculator.DeployExampleRewardCalculator(opts, ethClient, rewardBasisPoints)
	} else {
		return fmt.Errorf("invalid contract name: %s", helpers.GetDesiredContractName())
	}

	if err != nil {
		return fmt.Errorf("failed to deploy contract: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction confirmation: %s", err)
	}

	helpers.SaveEVMAddress(helpers.ValidatorManagerAddressPath, validatorManagerAddress)
	helpers.SaveEVMAddress(helpers.RewardCalculatorAddressPath, rewardCalculatorAddress)

	log.Printf("✅ Contract deployed at: %s\n", validatorManagerAddress.Hex())
	log.Printf("✅ Transaction hash: %s\n", receipt.TxHash.Hex())

	return nil
}
