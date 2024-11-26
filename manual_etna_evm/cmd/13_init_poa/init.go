package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mypkg/lib"
	"os"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"

	"github.com/ava-labs/avalanchego/api/info"

	pluginEVM "github.com/ava-labs/coreth/plugin/evm"

	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/key"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	blockchainSDK "github.com/ava-labs/avalanche-cli/sdk/blockchain"
	"github.com/ava-labs/avalanchego/ids"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/logging"
)

func main() {
	poaInitializedFilepath := "data/poa_initialized.txt"
	if _, err := os.Stat(poaInitializedFilepath); err == nil {
		fmt.Println("✅ PoA already initialized")
		return
	}

	// Read subnet ID
	subnetIDBytes, err := os.ReadFile("data/subnet.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}
	subnetID := ids.FromStringOrPanic(string(subnetIDBytes))

	// Read blockchain ID
	chainIDBytes, err := os.ReadFile("data/chain.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	blockchainID := ids.FromStringOrPanic(string(chainIDBytes))

	// Read owner key and get address
	ownerKey, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
	}

	// Convert to Ethereum address format
	ownerEthAddress := pluginEVM.PublicKeyToEthAddress(ownerKey.PublicKey())

	softKey, err := key.NewSoft(avagoconstants.TestnetID, key.WithPrivateKey(ownerKey))
	if err != nil {
		log.Fatalf("❌ Failed to create change owner address: %s\n", err)
	}

	changeOwnerAddress := softKey.P()[0]

	extraPeers := []string{}

	// Get bootstrap validators from convert_chain.go
	validators := []models.SubnetValidator{}
	for nodeNumber := 0; nodeNumber < lib.VALIDATORS_COUNT; nodeNumber++ {
		configBytes, err := os.ReadFile(fmt.Sprintf("data/configs/config-node%d.json", nodeNumber))
		if err != nil {
			log.Fatalf("❌ Failed to read config file: %s\n", err)
		}
		nodeConfig := lib.NodeConfig{}
		if err := json.Unmarshal(configBytes, &nodeConfig); err != nil {
			log.Fatalf("❌ Failed to unmarshal config: %s\n", err)
		}

		extraPeers = append(extraPeers, fmt.Sprintf("http://%s:%s", nodeConfig.PublicIP, nodeConfig.HTTPPort))

		endpoint := fmt.Sprintf("http://%s:%s", nodeConfig.PublicIP, nodeConfig.HTTPPort)
		infoClient := info.NewClient(endpoint)

		nodeID, proofOfPossession, err := infoClient.GetNodeID(context.Background())
		if err != nil {
			log.Fatalf("❌ Failed to get node info: %s\n", err)
		}

		publicKey := "0x" + hex.EncodeToString(proofOfPossession.PublicKey[:])
		pop := "0x" + hex.EncodeToString(proofOfPossession.ProofOfPossession[:])

		validator := models.SubnetValidator{
			NodeID:               nodeID.String(),
			Weight:               constants.BootstrapValidatorWeight,
			Balance:              constants.BootstrapValidatorBalance,
			BLSPublicKey:         publicKey,
			BLSProofOfPossession: pop,
			ChangeOwnerAddr:      changeOwnerAddress,
		}
		validators = append(validators, validator)
	}

	avaGoBootstrapValidators, err := blockchaincmd.ConvertToAvalancheGoSubnetValidator(validators)
	if err != nil {
		log.Fatalf("❌ Failed to convert to AvalancheGo subnet validator: %s\n", err)
	}

	subnetSDK := blockchainSDK.Subnet{
		SubnetID:            subnetID,
		BlockchainID:        blockchainID,
		OwnerAddress:        &ownerEthAddress,
		RPC:                 fmt.Sprintf("%s/ext/bc/%s/rpc", extraPeers[0], blockchainID),
		BootstrapValidators: avaGoBootstrapValidators,
	}

	subnetSDK.BootstrapValidators = avaGoBootstrapValidators

	// Get private key hex string for genesis
	genesisPrivateKey := hex.EncodeToString(ownerKey.Bytes())

	network := models.Network{
		Kind:        models.Fuji,
		ID:          avagoconstants.TestnetID,
		Endpoint:    extraPeers[0],
		ClusterName: "",
	}

	peers, err := blockchaincmd.ConvertURIToPeers(extraPeers)
	if err != nil {
		log.Fatalf("❌ Failed to get extra peers: %s\n", err)
	}

	// Initialize PoA
	if err := subnetSDK.InitializeProofOfAuthority(
		network,
		genesisPrivateKey,
		peers,
		logging.Debug,
	); err != nil {
		log.Fatalf("❌ Failed to initialize Proof of Authority: %s\n", err)
	}

	fmt.Println("✅ Successfully initialized Proof of Authority")

	// Create file to indicate that PoA is initialized
	if err := os.WriteFile(poaInitializedFilepath, []byte("true"), 0644); err != nil {
		log.Fatalf("❌ Failed to create PoA initialized file: %s\n", err)
	}
}
