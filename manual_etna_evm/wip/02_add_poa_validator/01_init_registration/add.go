package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	warpMessage "github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ethereum/go-ethereum/common"
)

func noErrVal[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func main() {
	err := os.RemoveAll(helpers.AddValidatorFolder)
	if err != nil {
		log.Fatalf("❌ Failed to remove add validator node folder: %s\n", err)
	}

	log.Printf("Cleaned up add validator folder %s\n", helpers.AddValidatorFolder)

	if helpers.FileExists(helpers.AddValidatorValidationIdPath) {
		log.Printf("✅ Validation ID already exists, skipping initialization\n")
		return
	}

	credshelper.GenerateCredsIfNotExists(helpers.AddValidatorKeysFolder)

	nodeID, proofOfPossession := credshelper.NodeInfoFromCreds(helpers.AddValidatorKeysFolder)

	chainID := helpers.LoadId(helpers.ChainIdPath)
	evmChainURL := fmt.Sprintf("http://127.0.0.1:9650/ext/bc/%s/rpc", chainID)

	expiry := uint64(time.Now().Add(constants.DefaultValidationIDExpiryDuration).Unix())
	err = helpers.SaveUint64(helpers.AddValidatorExpiryPath, expiry)
	if err != nil {
		return fmt.Errorf("failed to save expiry: %s\n", err)
	}

	managerKey := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)

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
		log.Fatalf("failed to get extra peers: %w", err)
	}

	blsPublicKey := [48]byte(proofOfPossession.PublicKey[:])
	weight := constants.NonBootstrapValidatorWeight

	subnetID := helpers.LoadId(helpers.SubnetIdPath)
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
		log.Fatalf("failed to get subnet validator registration message: %s", err)
	}

	err = helpers.SaveHex(helpers.AddValidatorWarpMessagePath, signedMessage.Bytes())
	if err != nil {
		return fmt.Errorf("failed to save warp message: %s\n", err)
	}

	fmt.Printf("validationID: %s\n", validationID)

	err = helpers.SaveId(helpers.AddValidatorValidationIdPath, validationID)
	if err != nil {
		return fmt.Errorf("failed to save validation ID: %s\n", err)
	}

}
