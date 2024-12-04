package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math/big"
	"mypkg/helpers"
	"time"

	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/pkg/utils"
	"github.com/ava-labs/avalanche-cli/pkg/ux"
	"github.com/ava-labs/avalanche-cli/pkg/validatormanager"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	validatorManagerSDK "github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/units"
	warp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	warpPayload "github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func AddValidator() error {
	signedMessage, validationID, err := InitValidatorRegistration(
		app,
		network,
		rpcURL,
		chainSpec,
		ownerPrivateKey,
		nodeID,
		blsInfo.PublicKey[:],
		expiry,
		remainingBalanceOwners,
		disableOwners,
		weight,
		extraAggregatorPeers,
		aggregatorLogLevel,
		pos,
		delegationFee,
		duration,
		big.NewInt(int64(stakeAmount)),
	)
	if err != nil {
		return err
	}
	log.Printf("ValidationID: %s", validationID)

	txID, _, err := deployer.RegisterL1Validator(balance, blsInfo, signedMessage)
	if err != nil {
		return err
	}
	log.Printf("RegisterL1ValidatorTx ID: %s", txID)

	if err := UpdatePChainHeight(
		"Waiting for P-Chain to update validator information ...",
	); err != nil {
		return err
	}

	if err := validatormanager.FinishValidatorRegistration(
		app,
		network,
		rpcURL,
		chainSpec,
		ownerPrivateKey,
		validationID,
		extraAggregatorPeers,
		aggregatorLogLevel,
	); err != nil {
		return err
	}

	log.Printf("  NodeID: %s", nodeID)
	log.Printf("  Network: %s", network.Name())
	// weight is inaccurate for PoS as it's fetched during registration
	if !pos {
		log.Printf("  Weight: %d", weight)
	}
	log.Printf("  Balance: %d", balance/units.Avax)
	ux.Logger.GreenCheckmarkToUser("Validator successfully added to the Subnet")

	return nil
}

func InitValidatorRegistration(
	rpcURL string,
	chainSpec contract.ChainSpec,
	ownerPrivateKey string,
	nodeID ids.NodeID,
	blsPublicKey []byte,
	expiry uint64,
	balanceOwners warpMessage.PChainOwner,
	disableOwners warpMessage.PChainOwner,
	weight uint64,
	aggregatorExtraPeerEndpoints []info.Peer,
	aggregatorLogLevelStr string,
	initWithPos bool,
	delegationFee uint16,
	stakeDuration time.Duration,
	stakeAmount *big.Int,
) (*warp.Message, ids.ID, error) {
	subnetID, err := helpers.LoadId("subnet")
	if err != nil {
		return nil, ids.Empty, err
	}

	chainID, err := helpers.LoadId("chain")
	if err != nil {
		return nil, ids.Empty, err
	}

	managerAddress := common.HexToAddress(validatorManagerSDK.ProxyContractAddress)

	managerAddress = common.HexToAddress(validatorManagerSDK.ProxyContractAddress)
	tx, _, err := PoAValidatorManagerInitializeValidatorRegistration(
		rpcURL,
		managerAddress,
		ownerPrivateKey,
		nodeID,
		blsPublicKey,
		expiry,
		balanceOwners,
		disableOwners,
		weight,
	)
	if err != nil {
		if !errors.Is(err, validatorManagerSDK.ErrNodeAlreadyRegistered) {
			return nil, ids.Empty, evm.TransactionError(tx, err, "failure initializing validator registration")
		}
		log.Printf("the validator registration was already initialized. Proceeding to the next step")
	}

	aggregatorLogLevel, err := logging.ToLevel(aggregatorLogLevelStr)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("invalid aggregator log level: %w", err)
	}

	log.Printf(fmt.Sprintf("Validator weight: %d", weight))
	return ValidatorManagerGetSubnetValidatorRegistrationMessage(
		network,
		aggregatorLogLevel,
		0,
		network.Kind == models.Local,
		aggregatorExtraPeerEndpoints,
		subnetID,
		chainID,
		managerAddress,
		nodeID,
		[48]byte(blsPublicKey),
		expiry,
		balanceOwners,
		disableOwners,
		weight,
	)
}

// step 1 of flow for adding a new validator
func PoAValidatorManagerInitializeValidatorRegistration(
	rpcURL string,
	managerAddress common.Address,
	managerOwnerPrivateKey string,
	nodeID ids.NodeID,
	blsPublicKey []byte,
	expiry uint64,
	balanceOwners warpMessage.PChainOwner,
	disableOwners warpMessage.PChainOwner,
	weight uint64,
) (*types.Transaction, *types.Receipt, error) {
	type PChainOwner struct {
		Threshold uint32
		Addresses []common.Address
	}
	type ValidatorRegistrationInput struct {
		NodeID                []byte
		BlsPublicKey          []byte
		RegistrationExpiry    uint64
		RemainingBalanceOwner PChainOwner
		DisableOwner          PChainOwner
	}
	balanceOwnersAux := PChainOwner{
		Threshold: balanceOwners.Threshold,
		Addresses: utils.Map(balanceOwners.Addresses, func(addr ids.ShortID) common.Address {
			return common.BytesToAddress(addr[:])
		}),
	}
	disableOwnersAux := PChainOwner{
		Threshold: disableOwners.Threshold,
		Addresses: utils.Map(disableOwners.Addresses, func(addr ids.ShortID) common.Address {
			return common.BytesToAddress(addr[:])
		}),
	}
	validatorRegistrationInput := ValidatorRegistrationInput{
		NodeID:                nodeID[:],
		BlsPublicKey:          blsPublicKey,
		RegistrationExpiry:    expiry,
		RemainingBalanceOwner: balanceOwnersAux,
		DisableOwner:          disableOwnersAux,
	}
	return contract.TxToMethod(
		rpcURL,
		managerOwnerPrivateKey,
		managerAddress,
		big.NewInt(0),
		"initialize validator registration",
		validatorManagerSDK.ErrorSignatureToError,
		"initializeValidatorRegistration((bytes,bytes,uint64,(uint32,[address]),(uint32,[address])),uint64)",
		validatorRegistrationInput,
		weight,
	)
}

func ValidatorManagerGetSubnetValidatorRegistrationMessage(
	network models.Network,
	aggregatorLogLevel logging.Level,
	aggregatorQuorumPercentage uint64,
	aggregatorAllowPrivateIPs bool,
	aggregatorExtraPeerEndpoints []info.Peer,
	subnetID ids.ID,
	blockchainID ids.ID,
	managerAddress common.Address,
	nodeID ids.NodeID,
	blsPublicKey [48]byte,
	expiry uint64,
	balanceOwners warpMessage.PChainOwner,
	disableOwners warpMessage.PChainOwner,
	weight uint64,
) (*warp.Message, ids.ID, error) {
	addressedCallPayload, err := warpMessage.NewRegisterL1Validator(
		subnetID,
		nodeID,
		blsPublicKey,
		expiry,
		balanceOwners,
		disableOwners,
		weight,
	)
	if err != nil {
		return nil, ids.Empty, err
	}

	validationID := addressedCallPayload.ValidationID()
	registerSubnetValidatorAddressedCall, err := warpPayload.NewAddressedCall(
		managerAddress.Bytes(),
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		return nil, ids.Empty, err
	}
	registerSubnetValidatorUnsignedMessage, err := warp.NewUnsignedMessage(
		network.ID,
		blockchainID,
		registerSubnetValidatorAddressedCall.Bytes(),
	)
	if err != nil {
		return nil, ids.Empty, err
	}
	signatureAggregator, err := interchain.NewSignatureAggregator(
		network,
		aggregatorLogLevel,
		subnetID,
		aggregatorQuorumPercentage,
		aggregatorAllowPrivateIPs,
		aggregatorExtraPeerEndpoints,
	)
	if err != nil {
		return nil, ids.Empty, err
	}
	signedMessage, err := signatureAggregator.Sign(registerSubnetValidatorUnsignedMessage, nil)
	return signedMessage, validationID, err
}
