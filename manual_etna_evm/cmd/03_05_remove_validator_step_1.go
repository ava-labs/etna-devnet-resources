package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/pkg/ux"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	"github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	warp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	warpPayload "github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeValidatorStep1Cmd)
}

var removeValidatorStep1Cmd = &cobra.Command{
	Use:   "remove-validator",
	Short: "Remove validator",
	RunE: func(cmd *cobra.Command, args []string) error {
		subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to load subnet ID: %w", err)
		}

		validatorsResp, err := callPChainValidatorsAt("https://api.avax-test.network/ext/P", subnetID.String())
		if err != nil {
			return fmt.Errorf("failed to get validators: %w", err)
		}

		if len(args) != 1 {
			fmt.Println("Existing validators:")
			for nodeID, details := range validatorsResp.Validators {
				fmt.Printf("Node ID: %s, Public Key: %s, Weight: %s\n", nodeID, details.PublicKey, details.Weight)
			}

			return errors.New("expected NodeID as argument")
		}

		nodeID, err := ids.NodeIDFromString(args[0])
		if err != nil {
			log.Fatalf("failed to parse node ID: %s", err)
		}

		signedMessage, validationID, err := InitValidatorRemoval(nodeID)
		if err != nil {
			return fmt.Errorf("failed to initialize validator removal: %w", err)
		}

		log.Printf("Signed message: %x\n", signedMessage.Bytes())
		log.Printf("Validation ID: %s\n", validationID.String())

		// Check if nodeID exists in validatorsResp
		_, exists := validatorsResp.Validators[nodeID.String()]
		if !exists {
			log.Printf("NodeID %s not found in current validators, skipping weight update", nodeID.String())
		} else {
			_, _, err = SetL1ValidatorWeight(signedMessage)
			if err != nil {
				return fmt.Errorf("failed to set L1 validator weight: %w", err)
			}

			log.Println("Waiting for 30 seconds before proceeding to the next step...")
			time.Sleep(30 * time.Second)
		}

		if err := FinishValidatorRemoval(validationID); err != nil {
			return fmt.Errorf("failed to finish validator removal: %w", err)
		}

		return nil
	},
}

func InitValidatorRemoval(nodeId ids.NodeID) (*warp.Message, ids.ID, error) {
	chainID, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to load chain ID: %w", err)
	}
	nodeURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)

	privateKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to load private key: %w", err)
	}

	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	validationID, err := GetRegisteredValidator(nodeURL, managerAddress, nodeId)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to get registered validator: %w", err)
	}

	tx, _, err := contract.TxToMethod(
		nodeURL,
		hex.EncodeToString(privateKey.Bytes()),
		managerAddress,
		big.NewInt(0),
		"POA validator removal initialization",
		validatormanager.ErrorSignatureToError,
		"initializeEndValidation(bytes32)",
		validationID,
	)
	if err != nil {
		if !errors.Is(err, validatormanager.ErrInvalidValidatorStatus) {
			return nil, ids.Empty, evm.TransactionError(tx, err, "failure initializing validator removal")
		}
		ux.Logger.PrintToUser("the validator removal process was already initialized. Proceeding to the next step")
	}

	network := models.NewFujiNetwork()
	aggregatorLogLevel := logging.Level(logging.Info)
	aggregatorQuorumPercentage := uint64(0)
	aggregatorAllowPrivateIPs := true
	aggregatorExtraPeerEndpoints, err := blockchaincmd.ConvertURIToPeers([]string{"http://127.0.0.1:9650"})
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to get extra peers: %w", err)
	}

	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to load subnet ID: %w", err)
	}

	blockchainID, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to load blockchain ID: %w", err)
	}

	nonce := uint64(1)
	weight := uint64(0)

	signedMsg, err := GetSubnetValidatorWeightMessage(
		network,
		aggregatorLogLevel,
		aggregatorQuorumPercentage,
		aggregatorAllowPrivateIPs,
		aggregatorExtraPeerEndpoints,
		subnetID,
		blockchainID,
		managerAddress,
		validationID,
		nonce,
		weight,
	)
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("failed to get subnet validator weight message: %w", err)
	}

	return signedMsg, validationID, nil
}

func GetRegisteredValidator(
	rpcURL string,
	managerAddress common.Address,
	nodeID ids.NodeID,
) (ids.ID, error) {
	out, err := contract.CallToMethod(
		rpcURL,
		managerAddress,
		"registeredValidators(bytes)->(bytes32)",
		nodeID[:],
	)
	if err != nil {
		return ids.Empty, err
	}
	validatorID, b := out[0].([32]byte)
	if !b {
		return ids.Empty, fmt.Errorf("error at registeredValidators call, expected [32]byte, got %T", out[0])
	}
	return validatorID, nil
}

func GetSubnetValidatorWeightMessage(
	network models.Network,
	aggregatorLogLevel logging.Level,
	aggregatorQuorumPercentage uint64,
	aggregatorAllowPrivateIPs bool,
	aggregatorExtraPeerEndpoints []info.Peer,
	subnetID ids.ID,
	blockchainID ids.ID,
	managerAddress common.Address,
	validationID ids.ID,
	nonce uint64,
	weight uint64,
) (*warp.Message, error) {
	addressedCallPayload, err := warpMessage.NewL1ValidatorWeight(
		validationID,
		nonce,
		weight,
	)
	if err != nil {
		return nil, err
	}
	addressedCall, err := warpPayload.NewAddressedCall(
		managerAddress.Bytes(),
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		return nil, err
	}
	unsignedMessage, err := warp.NewUnsignedMessage(
		network.ID,
		blockchainID,
		addressedCall.Bytes(),
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return signatureAggregator.Sign(unsignedMessage, nil)
}
