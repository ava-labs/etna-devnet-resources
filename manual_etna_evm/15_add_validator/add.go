package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"mypkg/config"
	"mypkg/helpers"
	"strings"
	"time"

	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/utils"
	validatorManagerSDK "github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/ids"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ports := []string{"9654", "9652"}
	for _, port := range ports {
		if err := AddValidator(fmt.Sprintf("http://127.0.0.1:%s", port)); err != nil {
			if strings.Contains(err.Error(), "node already registered") {
				log.Printf("✅ Node on port %s was already registered as validator\n", port)
				continue
			}
			log.Fatalf("❌ Failed to add validator on port %s: %s\n", port, err)
		}
	}
}

func AddValidator(rpcURL string) error {
	chainID, err := helpers.LoadId("chain")
	if err != nil {
		return fmt.Errorf("failed to load chain ID: %w", err)
	}

	evmChainURL := fmt.Sprintf("http://127.0.0.1:9650/ext/bc/%s/rpc", chainID)

	nodeID, blsInfo, err := helpers.GetNodeInfoRetry(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to get node info: %s", err)
	}

	expiry := uint64(time.Now().Add(constants.DefaultValidationIDExpiryDuration).Unix())

	key, err := helpers.LoadValidatorManagerKey()
	if err != nil {
		return fmt.Errorf("failed to load key from file: %s", err)
	}

	managerKey, err := helpers.LoadValidatorManagerKey()
	if err != nil {
		return fmt.Errorf("failed to load key from file: %s", err)
	}

	pChainAddr := key.Address()

	remainingBalanceOwners := warpMessage.PChainOwner{
		Threshold: 1,
		Addresses: []ids.ShortID{pChainAddr},
	}
	disableOwners := remainingBalanceOwners

	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	_, receipt, err := PoAValidatorManagerInitializeValidatorRegistration(
		evmChainURL,
		managerAddress,
		hex.EncodeToString(managerKey.Bytes()),
		nodeID,
		blsInfo.PublicKey[:],
		expiry,
		remainingBalanceOwners,
		disableOwners,
		constants.NonBootstrapValidatorWeight,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize validator registration: %s", err)
	}

	log.Printf("✅ Validator registration initialized: %s\n", receipt.TxHash)

	return nil
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
	log.Printf("Initializing validator registration with:\n"+
		"rpcURL: %s\n"+
		"managerAddress: %s\n"+
		"managerOwnerPrivateKey: %s\n"+
		"nodeID: %s\n"+
		"blsPublicKey: %x\n"+
		"expiry: %d\n"+
		"balanceOwners:\n"+
		"\tThreshold: %d\n"+
		"\tAddresses: %v\n"+
		"disableOwners:\n"+
		"\tThreshold: %d\n"+
		"\tAddresses: %v\n"+
		"weight: %d\n",
		rpcURL,
		managerAddress.String(),
		managerOwnerPrivateKey,
		nodeID.String(),
		blsPublicKey,
		expiry,
		balanceOwners.Threshold,
		balanceOwners.Addresses,
		disableOwners.Threshold,
		disableOwners.Addresses,
		weight,
	)

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

	//TODO:
	// return ValidatorManagerGetSubnetValidatorRegistrationMessage(
	// 	network,
	// 	aggregatorLogLevel,
	// 	0,
	// 	network.Kind == models.Local,
	// 	aggregatorExtraPeerEndpoints,
	// 	subnetID,
	// 	blockchainID,
	// 	managerAddress,
	// 	nodeID,
	// 	[48]byte(blsPublicKey),
	// 	expiry,
	// 	balanceOwners,
	// 	disableOwners,
	// 	weight,
	// )

	//TODO:
	//txID, _, err := deployer.RegisterL1Validator(balance, blsInfo, signedMessage)

	//TODO:
	// if err := UpdatePChainHeight(

	//TODO:

	// if err := validatormanager.FinishValidatorRegistration(
	// 	app,
	// 	network,
	// 	rpcURL,
	// 	chainSpec,
	// 	ownerPrivateKey,
	// 	validationID,
	// 	extraAggregatorPeers,
	// 	aggregatorLogLevel,
	// ); err != nil {
	// 	return err
	// }
}

// func AddValidator() error {
// 	signedMessage, validationID, err := InitValidatorRegistration(
// 		app,
// 		network,
// 		rpcURL,
// 		chainSpec,
// 		ownerPrivateKey,
// 		nodeID,
// 		blsInfo.PublicKey[:],
// 		expiry,
// 		remainingBalanceOwners,
// 		disableOwners,
// 		weight,
// 		extraAggregatorPeers,
// 		aggregatorLogLevel,
// 		pos,
// 		delegationFee,
// 		duration,
// 		big.NewInt(int64(stakeAmount)),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("ValidationID: %s", validationID)

// 	txID, _, err := deployer.RegisterL1Validator(balance, blsInfo, signedMessage)
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("RegisterL1ValidatorTx ID: %s", txID)

// 	if err := UpdatePChainHeight(
// 		"Waiting for P-Chain to update validator information ...",
// 	); err != nil {
// 		return err
// 	}

// 	if err := validatormanager.FinishValidatorRegistration(
// 		app,
// 		network,
// 		rpcURL,
// 		chainSpec,
// 		ownerPrivateKey,
// 		validationID,
// 		extraAggregatorPeers,
// 		aggregatorLogLevel,
// 	); err != nil {
// 		return err
// 	}

// 	log.Printf("  NodeID: %s", nodeID)
// 	log.Printf("  Network: %s", network.Name())
// 	// weight is inaccurate for PoS as it's fetched during registration
// 	if !pos {
// 		log.Printf("  Weight: %d", weight)
// 	}
// 	log.Printf("  Balance: %d", balance/units.Avax)
// 	ux.Logger.GreenCheckmarkToUser("Validator successfully added to the Subnet")

// 	return nil
// }

// func InitValidatorRegistration(
// 	rpcURL string,
// 	chainSpec contract.ChainSpec,
// 	ownerPrivateKey string,
// 	nodeID ids.NodeID,
// 	blsPublicKey []byte,
// 	expiry uint64,
// 	balanceOwners warpMessage.PChainOwner,
// 	disableOwners warpMessage.PChainOwner,
// 	weight uint64,
// 	aggregatorExtraPeerEndpoints []info.Peer,
// 	aggregatorLogLevelStr string,
// 	initWithPos bool,
// 	delegationFee uint16,
// 	stakeDuration time.Duration,
// 	stakeAmount *big.Int,
// ) (*warp.Message, ids.ID, error) {
// 	subnetID, err := helpers.LoadId("subnet")
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}

// 	chainID, err := helpers.LoadId("chain")
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}

// 	managerAddress := common.HexToAddress(validatorManagerSDK.ProxyContractAddress)

// 	managerAddress = common.HexToAddress(validatorManagerSDK.ProxyContractAddress)
// 	tx, _, err := PoAValidatorManagerInitializeValidatorRegistration(
// 		rpcURL,
// 		managerAddress,
// 		ownerPrivateKey,
// 		nodeID,
// 		blsPublicKey,
// 		expiry,
// 		balanceOwners,
// 		disableOwners,
// 		weight,
// 	)
// 	if err != nil {
// 		if !errors.Is(err, validatorManagerSDK.ErrNodeAlreadyRegistered) {
// 			return nil, ids.Empty, evm.TransactionError(tx, err, "failure initializing validator registration")
// 		}
// 		log.Printf("the validator registration was already initialized. Proceeding to the next step")
// 	}

// 	aggregatorLogLevel, err := logging.ToLevel(aggregatorLogLevelStr)
// 	if err != nil {
// 		return nil, ids.Empty, fmt.Errorf("invalid aggregator log level: %w", err)
// 	}

// 	log.Printf(fmt.Sprintf("Validator weight: %d", weight))
// 	return ValidatorManagerGetSubnetValidatorRegistrationMessage(
// 		network,
// 		aggregatorLogLevel,
// 		0,
// 		network.Kind == models.Local,
// 		aggregatorExtraPeerEndpoints,
// 		subnetID,
// 		chainID,
// 		managerAddress,
// 		nodeID,
// 		[48]byte(blsPublicKey),
// 		expiry,
// 		balanceOwners,
// 		disableOwners,
// 		weight,
// 	)
// }

// func ValidatorManagerGetSubnetValidatorRegistrationMessage(
// 	network models.Network,
// 	aggregatorLogLevel logging.Level,
// 	aggregatorQuorumPercentage uint64,
// 	aggregatorAllowPrivateIPs bool,
// 	aggregatorExtraPeerEndpoints []info.Peer,
// 	subnetID ids.ID,
// 	blockchainID ids.ID,
// 	managerAddress common.Address,
// 	nodeID ids.NodeID,
// 	blsPublicKey [48]byte,
// 	expiry uint64,
// 	balanceOwners warpMessage.PChainOwner,
// 	disableOwners warpMessage.PChainOwner,
// 	weight uint64,
// ) (*warp.Message, ids.ID, error) {
// 	addressedCallPayload, err := warpMessage.NewRegisterL1Validator(
// 		subnetID,
// 		nodeID,
// 		blsPublicKey,
// 		expiry,
// 		balanceOwners,
// 		disableOwners,
// 		weight,
// 	)
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}

// 	validationID := addressedCallPayload.ValidationID()
// 	registerSubnetValidatorAddressedCall, err := warpPayload.NewAddressedCall(
// 		managerAddress.Bytes(),
// 		addressedCallPayload.Bytes(),
// 	)
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}
// 	registerSubnetValidatorUnsignedMessage, err := warp.NewUnsignedMessage(
// 		network.ID,
// 		blockchainID,
// 		registerSubnetValidatorAddressedCall.Bytes(),
// 	)
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}
// 	signatureAggregator, err := interchain.NewSignatureAggregator(
// 		network,
// 		aggregatorLogLevel,
// 		subnetID,
// 		aggregatorQuorumPercentage,
// 		aggregatorAllowPrivateIPs,
// 		aggregatorExtraPeerEndpoints,
// 	)
// 	if err != nil {
// 		return nil, ids.Empty, err
// 	}
// 	signedMessage, err := signatureAggregator.Sign(registerSubnetValidatorUnsignedMessage, nil)
// 	return signedMessage, validationID, err
// }
