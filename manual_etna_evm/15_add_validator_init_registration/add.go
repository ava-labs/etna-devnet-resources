package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/pkg/utils"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	validatorManagerSDK "github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	warpPayload "github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
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
	warpMessageExists, err := helpers.TextFileExists("add_validator_warp_message")
	if err != nil {
		return fmt.Errorf("failed to check if warp message exists: %w", err)
	}

	if warpMessageExists {
		log.Printf("✅ Warp message already exists, skipping initialization\n")
		return nil
	}

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

	err = helpers.SaveHex("add_validator_warp_message", signedMessage.Bytes())
	if err != nil {
		return fmt.Errorf("saving validator warp message: %w", err)
	}

	BLSInfoJSON, err := blsInfo.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling BLS: %w", err)
	}

	err = helpers.SaveText("add_validator_bls_json", string(BLSInfoJSON))
	if err != nil {
		return fmt.Errorf("saving BLS JSON: %w", err)
	}

	fmt.Printf("validationID: %s\n", validationID)

	err = helpers.SaveId("add_validator_validation_id", validationID)
	if err != nil {
		return fmt.Errorf("saving validation ID: %w", err)
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

func loadOrGenerateExpiry() (uint64, error) {
	expiryFile := "add_validator_expiry"
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
