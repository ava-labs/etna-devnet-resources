package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"mypkg/config"
	"mypkg/helpers"
	"strings"
	"time"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/pkg/utils"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	validatorManagerSDK "github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	warpPayload "github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	port := "9652"
	log.Printf("Adding validator on port %s\n", port)
	if err := AddValidator(fmt.Sprintf("http://127.0.0.1:%s", port)); err != nil {
		log.Fatalf("❌ Failed to add validator on port %s: %s\n", port, err)
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

	expiry, err := loadOrGenerateExpiry()
	if err != nil {
		return fmt.Errorf("failed to load or generate expiry: %s", err)
	}

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
		if strings.Contains(err.Error(), "node already registered") {
			log.Printf("reverted with an expected error: %s", err)
			log.Printf("✅ Node %s was already registered as validator previously\n", rpcURL)
		} else {
			return fmt.Errorf("failed to initialize validator registration: %s", err)
		}
	} else {
		log.Printf("✅ Validator registration initialized: %s\n", receipt.TxHash)
	}

	network := models.NewFujiNetwork()
	aggregatorLogLevel := logging.Level(logging.Info)
	aggregatorQuorumPercentage := uint64(0)
	aggregatorAllowPrivateIPs := true

	aggregatorExtraPeerEndpoints, err := blockchaincmd.ConvertURIToPeers([]string{"http://127.0.0.1:9650"})
	if err != nil {
		return fmt.Errorf("failed to get extra peers: %w", err)
	}

	subnetID, err := helpers.LoadId("subnet")
	if err != nil {
		return fmt.Errorf("failed to load subnet ID: %w", err)
	}

	blsPublicKey := [48]byte(blsInfo.PublicKey[:])
	weight := constants.NonBootstrapValidatorWeight

	signedMessage, validationID, err := ValidatorManagerGetSubnetValidatorRegistrationMessage(
		network,
		aggregatorLogLevel,
		aggregatorQuorumPercentage,
		aggregatorAllowPrivateIPs,
		aggregatorExtraPeerEndpoints,
		subnetID,
		chainID,
		managerAddress,
		nodeID,
		blsPublicKey,
		expiry,
		remainingBalanceOwners,
		disableOwners,
		uint64(weight),
	)
	if err != nil {
		return fmt.Errorf("failed to get subnet validator registration message: %s", err)
	}

	_ = signedMessage
	_ = validationID

	fmt.Printf("signedMessage: %s\n", signedMessage)
	fmt.Printf("validationID: %s\n", validationID)

	balance := 1 * units.Avax
	var proofOfPossession [96]byte
	copy(proofOfPossession[:], blsInfo.ProofOfPossession[:])
	if err := RegisterL1ValidatorOnPChain(key, balance, proofOfPossession, signedMessage); err != nil {
		return fmt.Errorf("failed to register L1 validator on P-chain: %s", err)
	}

	log.Printf("Waiting for P-chain to update validator information ...")
	time.Sleep(30 * time.Second)

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

	return nil
}

func RegisterL1ValidatorOnPChain(key *secp256k1.PrivateKey, balance uint64, proofOfPossession [96]byte, message *warp.Message) error {
	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          config.RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}

	unsignedTx, err := wallet.P().Builder().NewRegisterL1ValidatorTx(
		balance,
		proofOfPossession,
		message.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("error building tx: %w", err)
	}

	tx := txs.Tx{Unsigned: unsignedTx}
	if err := wallet.P().Signer().Sign(context.Background(), &tx); err != nil {
		return fmt.Errorf("error signing tx: %w", err)
	}

	err = wallet.P().IssueTx(&tx)
	if err != nil {
		return fmt.Errorf("error issuing tx: %w", err)
	}

	return nil
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
	// log.Printf("Initializing validator registration with:\n"+
	// 	"rpcURL: %s\n"+
	// 	"managerAddress: %s\n"+
	// 	"managerOwnerPrivateKey: %s\n"+
	// 	"nodeID: %s\n"+
	// 	"blsPublicKey: %x\n"+
	// 	"expiry: %d\n"+
	// 	"balanceOwners:\n"+
	// 	"\tThreshold: %d\n"+
	// 	"\tAddresses: %v\n"+
	// 	"disableOwners:\n"+
	// 	"\tThreshold: %d\n"+
	// 	"\tAddresses: %v\n"+
	// 	"weight: %d\n",
	// 	rpcURL,
	// 	managerAddress.String(),
	// 	managerOwnerPrivateKey,
	// 	nodeID.String(),
	// 	blsPublicKey,
	// 	expiry,
	// 	balanceOwners.Threshold,
	// 	balanceOwners.Addresses,
	// 	disableOwners.Threshold,
	// 	disableOwners.Addresses,
	// 	weight,
	// )

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

//		validationID := addressedCallPayload.ValidationID()
//		registerSubnetValidatorAddressedCall, err := warpPayload.NewAddressedCall(
//			managerAddress.Bytes(),
//			addressedCallPayload.Bytes(),
//		)
//		if err != nil {
//			return nil, ids.Empty, err
//		}
//		registerSubnetValidatorUnsignedMessage, err := warp.NewUnsignedMessage(
//			network.ID,
//			blockchainID,
//			registerSubnetValidatorAddressedCall.Bytes(),
//		)
//		if err != nil {
//			return nil, ids.Empty, err
//		}
//		signatureAggregator, err := interchain.NewSignatureAggregator(
//			network,
//			aggregatorLogLevel,
//			subnetID,
//			aggregatorQuorumPercentage,
//			aggregatorAllowPrivateIPs,
//			aggregatorExtraPeerEndpoints,
//		)
//		if err != nil {
//			return nil, ids.Empty, err
//		}
//		signedMessage, err := signatureAggregator.Sign(registerSubnetValidatorUnsignedMessage, nil)
//		return signedMessage, validationID, err
//	}
func loadOrGenerateExpiry() (uint64, error) {
	expiryFile := "validator_expiry"
	exists, err := helpers.TextFileExists(expiryFile)
	if err != nil {
		return 0, fmt.Errorf("failed to check if expiry file exists: %w", err)
	}

	if !exists {
		expiry := uint64(time.Now().Add(constants.DefaultValidationIDExpiryDuration).Unix())
		if err := helpers.SaveUint64(expiryFile, expiry); err != nil {
			return 0, fmt.Errorf("failed to save expiry: %w", err)
		}
		return expiry, nil
	}

	expiry, err := helpers.LoadUint64(expiryFile)
	if err != nil {
		return 0, fmt.Errorf("failed to load expiry: %w", err)
	}
	return expiry, nil
}
