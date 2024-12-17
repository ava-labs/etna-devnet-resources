package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/spf13/cobra"

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

func init() {
	rootCmd.AddCommand(AddValidatorCmd)
}

var AddValidatorCmd = &cobra.Command{
	Use:   "add-validator",
	Short: "Add a validator to the validator set",
	RunE: func(cmd *cobra.Command, args []string) error {
		credsFolder, err := generateAddValidatorFolder()
		if err != nil {
			return fmt.Errorf("failed to generate add validator folder: %w", err)
		}

		err = GenerateCredsIfNotExists(credsFolder)
		if err != nil {
			return fmt.Errorf("failed to generate creds: %w", err)
		}

		warpMessage, validationID, expiry, err := InitValidatorRegistration(credsFolder)
		if err != nil {
			return fmt.Errorf("failed to initialize validator registration: %w", err)
		}

		log.Printf("Validator registration initialized: %x\n", warpMessage.Bytes())
		log.Printf("Validation ID: %s\n", validationID)
		log.Printf("Expiry: %d\n", expiry)

		pChainRegistrationCompleted := false
		for i := 0; i < 5; i++ {
			log.Printf("Attempting to register L1 validator on P-chain (attempt %d/5)...", i+1)
			err = RegisterL1ValidatorOnPChain(warpMessage, credsFolder)
			if err != nil {
				log.Printf("Attempt %d failed: %s", i+1, err)
				if i < 4 {
					log.Printf("Waiting 10 seconds before retrying...")
					time.Sleep(10 * time.Second)
					continue
				}
				return fmt.Errorf("all attempts to register L1 validator failed: %w", err)
			}
			pChainRegistrationCompleted = true
			log.Printf("Successfully registered L1 validator on P-chain")
			break
		}

		if !pChainRegistrationCompleted {
			return fmt.Errorf("failed to register L1 validator on P-chain")
		}

		err = AddValidatorCompleteRegistration(validationID)
		if err != nil {
			return fmt.Errorf("failed to complete validator registration: %w", err)
		}

		return nil
	},
}

func generateAddValidatorFolder() (string, error) {
	for i := 0; i < 100; i++ {
		folderName := fmt.Sprintf("data/add_validator_%d", i)
		exists, err := helpers.FileExists(folderName)
		if err != nil {
			return "", fmt.Errorf("failed to check if folder exists: %w", err)
		}
		if exists {
			continue
		}
		err = os.MkdirAll(folderName, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create folder: %w", err)
		}
		return folderName, nil
	}
	return "", fmt.Errorf("failed to generate add validator folder")
}

func InitValidatorRegistration(credsFolder string) (*warp.Message, ids.ID, uint64, error) {
	nodeID, proofOfPossession, err := NodeInfoFromCreds(credsFolder)
	if err != nil {
		return nil, ids.Empty, 0, fmt.Errorf("failed to get node info from creds: %w", err)
	}

	chainID, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return nil, ids.Empty, 0, fmt.Errorf("failed to load chain ID: %w", err)
	}
	evmChainURL := fmt.Sprintf("http://127.0.0.1:9650/ext/bc/%s/rpc", chainID)

	expiry := uint64(time.Now().Add(constants.DefaultValidationIDExpiryDuration).Unix())

	managerKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return nil, ids.Empty, 0, fmt.Errorf("failed to load manager key: %w", err)
	}

	pChainAddr := managerKey.Address()

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
		proofOfPossession.PublicKey[:],
		expiry,
		remainingBalanceOwners,
		disableOwners,
		constants.NonBootstrapValidatorWeight,
	)
	if err != nil {
		if strings.Contains(err.Error(), "node already registered") {
			log.Printf("reverted with an expected error: %s", err)
			log.Printf("✅ Node %s was already registered as validator previously\n", nodeID)
		} else {
			log.Fatalf("failed to initialize validator registration: %s", err)
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
		return nil, ids.Empty, 0, fmt.Errorf("failed to get extra peers: %w", err)
	}

	blsPublicKey := [48]byte(proofOfPossession.PublicKey[:])
	weight := constants.NonBootstrapValidatorWeight

	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return nil, ids.Empty, 0, fmt.Errorf("failed to load subnet ID: %w", err)
	}

	warpMessage, validationID, err := ValidatorManagerGetSubnetValidatorRegistrationMessage(
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
		return nil, ids.Empty, 0, fmt.Errorf("failed to get subnet validator registration message: %w", err)
	}

	return warpMessage, validationID, uint64(expiry), nil
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
