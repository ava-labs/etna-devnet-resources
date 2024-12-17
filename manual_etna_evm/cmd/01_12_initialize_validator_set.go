package cmd

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	"github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ethereum/go-ethereum/common"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initializeValidatorSetCmd)
}

var initializeValidatorSetCmd = &cobra.Command{
	Use:   "initialize-validator-set",
	Short: "Initialize the validator set",
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🧱 Initializing validator set")

		const maxAttempts = 10
		const retryDelay = 10 * time.Second

		var lastErr error
		for i := 0; i < maxAttempts; i++ {
			if i > 0 {
				fmt.Printf("Attempt %d/%d (will sleep for %v before retry)\n", i+1, maxAttempts, retryDelay)
				time.Sleep(retryDelay)
			}

			err := initializeValidatorSet()
			if err == nil {
				return nil
			}
			lastErr = err
			fmt.Printf("Failed to initialize validator set: %s, retrying...\n", err)
		}

		return fmt.Errorf("Failed to initialize validator set after %d attempts: %s\n", maxAttempts, lastErr)
	},
}

func initializeValidatorSet() error {
	alreadyInitialized, err := helpers.FileExists(helpers.InitializeValidatorSetTxPath)
	if err != nil {
		return fmt.Errorf("failed to check if validator set is already initialized: %w", err)
	}
	if alreadyInitialized {
		log.Println("✅ Validator set is already initialized")
		return nil
	}

	managerAddress := goethereumcommon.HexToAddress(config.ProxyContractAddress)

	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return fmt.Errorf("failed to load subnet ID: %w", err)
	}

	subnetConversionIDFromFile, err := helpers.LoadId(helpers.ConversionIdPath)
	if err != nil {
		return fmt.Errorf("failed to load subnet conversion ID: %w", err)
	}

	_ = subnetConversionIDFromFile

	nodeID, proofOfPossession, err := helpers.GetNodeInfoRetry(fmt.Sprintf("http://%s:%s", "127.0.0.1", "9650"))
	if err != nil {
		return fmt.Errorf("failed to get node info: %w", err)
	}

	chainID, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return fmt.Errorf("failed to load chain ID: %w", err)
	}

	validators := []message.SubnetToL1ConverstionValidatorData{}
	validators = append(validators, message.SubnetToL1ConverstionValidatorData{
		NodeID:       nodeID[:],
		BLSPublicKey: proofOfPossession.PublicKey,
		Weight:       constants.BootstrapValidatorWeight,
	})

	subnetConversionData := message.SubnetToL1ConversionData{
		SubnetID:       subnetID,
		ManagerChainID: chainID,
		ManagerAddress: managerAddress.Bytes(),
		Validators:     validators,
	}
	subnetConversionID, err := message.SubnetToL1ConversionID(subnetConversionData)
	if err != nil {
		return fmt.Errorf("failed to create subnet conversion ID: %w", err)
	}

	addressedCallPayload, err := message.NewSubnetToL1Conversion(subnetConversionID)
	if err != nil {
		return fmt.Errorf("failed to create addressed call payload: %w", err)
	}

	subnetConversionAddressedCall, err := payload.NewAddressedCall(
		nil,
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("failed to create addressed call payload: %w", err)
	}

	network := models.NewFujiNetwork()

	subnetConversionUnsignedMessage, err := warp.NewUnsignedMessage(
		network.ID,
		avagoconstants.PlatformChainID,
		subnetConversionAddressedCall.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("failed to create unsigned message: %w", err)
	}

	peers, err := blockchaincmd.ConvertURIToPeers([]string{"http://127.0.0.1:9650"})
	if err != nil {
		return fmt.Errorf("failed to get extra peers: %w", err)
	}

	signatureAggregator, err := interchain.NewSignatureAggregator(
		network,
		logging.Level(logging.Info),
		subnetID,
		interchain.DefaultQuorumPercentage,
		true,
		peers,
	)
	if err != nil {
		return fmt.Errorf("failed to create signature aggregator: %w", err)
	}

	subnetConversionSignedMessage, err := signatureAggregator.Sign(subnetConversionUnsignedMessage, subnetID[:])
	if err != nil {
		return fmt.Errorf("failed to sign subnet conversion unsigned message: %w", err)
	}

	//as plain text
	privateKey, err := helpers.LoadText(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	type InitialValidatorPayload struct {
		NodeID       []byte
		BlsPublicKey []byte
		Weight       uint64
	}
	type SubnetConversionDataPayload struct {
		SubnetID                     [32]byte
		ValidatorManagerBlockchainID [32]byte
		ValidatorManagerAddress      common.Address
		InitialValidators            []InitialValidatorPayload
	}

	subnetConversionDataPayload := SubnetConversionDataPayload{
		SubnetID:                     subnetID,
		ValidatorManagerBlockchainID: chainID,
		ValidatorManagerAddress:      managerAddress,
		InitialValidators: []InitialValidatorPayload{
			{
				NodeID:       nodeID[:],
				BlsPublicKey: proofOfPossession.PublicKey[:],
				Weight:       constants.BootstrapValidatorWeight,
			},
		},
	}

	tx, _, err := contract.TxToMethodWithWarpMessage(
		fmt.Sprintf("http://127.0.0.1:9650/ext/bc/%s/rpc", chainID),
		strings.TrimSpace(privateKey),
		managerAddress,
		subnetConversionSignedMessage,
		big.NewInt(0),
		"initialize validator set",
		validatormanager.ErrorSignatureToError,
		"initializeValidatorSet((bytes32,bytes32,address,[(bytes,bytes,uint64)]),uint32)",
		subnetConversionDataPayload,
		uint32(0),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize validator set: %w", err)
	}

	fmt.Printf("✅ Successfully initialized validator set. Transaction hash: %s\n", tx.Hash().String())

	helpers.SaveText(helpers.InitializeValidatorSetTxPath, tx.Hash().String())

	return nil
}
