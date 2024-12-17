package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/contract"
	"github.com/ava-labs/avalanche-cli/pkg/evm"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	validatorManagerSDK "github.com/ava-labs/avalanche-cli/sdk/validatormanager"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
)

//TODO:
//FinishValidatorRemoval
//- GetPChainSubnetValidatorRegistrationWarpMessage
//- SetupProposerVM
//- CompleteValidatorRemoval

func FinishValidatorRemoval(validationID ids.ID) error {
	chainID, err := helpers.LoadId(helpers.ChainIdPath)
	if err != nil {
		return fmt.Errorf("failed to load chain ID: %w", err)
	}
	rpcURL := fmt.Sprintf("http://%s:%s/ext/bc/%s/rpc", "127.0.0.1", "9650", chainID)

	network := models.NewFujiNetwork()
	aggregatorLogLevel := logging.Level(logging.Info)
	aggregatorAllowPrivateIPs := true
	aggregatorQuorumPercentage := uint64(0)
	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return fmt.Errorf("failed to load subnet id: %w", err)
	}
	aggregatorExtraPeerEndpoints, err := blockchaincmd.ConvertURIToPeers([]string{"http://127.0.0.1:9650"})
	if err != nil {
		return fmt.Errorf("failed to get extra peers: %w", err)
	}
	registered := false

	signedMessage, err := ValidatorManagerGetPChainSubnetValidatorRegistrationWarpMessage(
		network,
		rpcURL,
		aggregatorLogLevel,
		aggregatorQuorumPercentage,
		aggregatorAllowPrivateIPs,
		aggregatorExtraPeerEndpoints,
		subnetID,
		validationID,
		registered,
	)
	if err != nil {
		log.Fatalf("failed to get P-chain subnet validator registration warp message: %s", err)
	}

	privateKey, err := helpers.LoadText(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load manager key: %w", err)
	}

	if err := evm.SetupProposerVM(
		rpcURL,
		privateKey,
	); err != nil {
		return err
	}

	managerAddress := goethereumcommon.HexToAddress(config.ProxyContractAddress)

	tx, _, err := contract.TxToMethodWithWarpMessage(
		rpcURL,
		privateKey,
		managerAddress,
		signedMessage,
		big.NewInt(0),
		"complete poa validator removal",
		validatorManagerSDK.ErrorSignatureToError,
		"completeEndValidation(uint32)",
		uint32(0),
	)
	if err != nil {
		return evm.TransactionError(tx, err, "failure completing validator removal")
	}
	return nil
}
