package main

import (
	"fmt"
	"log"
	"mypkg/helpers"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/sdk/interchain"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
)

func main() {
	if err := retrievePChainSubnetConversionWarpMessage(); err != nil {
		log.Fatalf("❌ Failed to retrieve P-Chain subnet conversion warp message: %s\n", err)
	}
}

func retrievePChainSubnetConversionWarpMessage() error {
	managerAddressHex, err := helpers.LoadText("validator_manager_address")
	if err != nil {
		log.Fatalf("❌ Failed to load validator manager address: %s\n", err)
	}

	managerAddress := goethereumcommon.HexToAddress(managerAddressHex)

	subnetID, err := helpers.LoadId("subnet")
	if err != nil {
		return fmt.Errorf("failed to load subnet ID: %w", err)
	}

	subnetConversionIDFromFile, err := helpers.LoadId("conversion_id")
	if err != nil {
		return fmt.Errorf("failed to load subnet conversion ID: %w", err)
	}

	_ = subnetConversionIDFromFile

	nodeID, proofOfPossession, err := helpers.GetNodeInfoRetry(fmt.Sprintf("http://%s:%s", "127.0.0.1", "9650"))
	if err != nil {
		return fmt.Errorf("failed to get node info: %w", err)
	}

	chainID, err := helpers.LoadId("chain")
	if err != nil {
		log.Fatalf("❌ Failed to load chain ID: %s\n", err)
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
		logging.Level(logging.Debug),
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

	fmt.Println("subnetConversionSignedMessage", subnetConversionSignedMessage.String())

	return nil
}
