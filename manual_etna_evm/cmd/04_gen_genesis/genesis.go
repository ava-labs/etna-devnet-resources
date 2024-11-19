package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"math/big"
	"mypkg/lib"
	"os"
	"time"

	icmgenesis "github.com/ava-labs/avalanche-cli/pkg/teleporter/genesis"
	"github.com/ava-labs/avalanche-cli/pkg/validatormanager"
	blockchainSDK "github.com/ava-labs/avalanche-cli/sdk/blockchain"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/coreth/utils"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var (
	OneAvax                = new(big.Int).SetUint64(1000000000000000000)
	defaultPoAOwnerBalance = new(big.Int).Mul(OneAvax, big.NewInt(10))
)

var (
	// 600 AVAX: to deploy teleporter contract, registry contract, and fund
	// starting relayer operations
	teleporterBalance = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(600))
	// 1000 AVAX: to deploy teleporter contract, registry contract, fund
	// starting relayer operations, and deploy bridge contracts
	externalGasTokenBalance = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
)

func main() {
	ownerKey, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("failed to load key from file: %s\n", err)
	}
	ethAddr := evm.PublicKeyToEthAddress(ownerKey.PublicKey())

	const CHAIN_ID = 12345

	feeConfig := commontype.FeeConfig{
		GasLimit:                 big.NewInt(12000000),
		TargetBlockRate:          2,
		MinBaseFee:               big.NewInt(25000000000),
		TargetGas:                big.NewInt(60000000),
		BaseFeeChangeDenominator: big.NewInt(36),
		MinBlockGasCost:          big.NewInt(0),
		MaxBlockGasCost:          big.NewInt(1000000),
		BlockGasCostStep:         big.NewInt(200000),
	}

	allocation := core.GenesisAlloc{
		common.Address{}: core.GenesisAccount{Balance: teleporterBalance},
	}
	allocation[ethAddr] = core.GenesisAccount{Balance: defaultPoAOwnerBalance}
	icmgenesis.AddICMMessengerContractToAllocations(allocation)

	teleporterAddress := common.HexToAddress("0x2656b6eb2ba42b6e0b5be696021f94de4d8a16d8")
	// TODO: figure out how to get the address of the teleporter contract
	allocation[teleporterAddress] = core.GenesisAccount{Balance: teleporterBalance}

	validatormanager.AddPoAValidatorManagerContractToAllocations(allocation)

	genesisTimestamp := utils.TimeToNewUint64(time.Now())

	precompiles := params.Precompiles{}
	precompiles[warp.ConfigKey] = &warp.Config{
		QuorumNumerator:              warp.WarpDefaultQuorumNumerator,
		RequirePrimaryNetworkSigners: true,
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: genesisTimestamp,
		},
	}

	subnetConfig, err := blockchainSDK.New(
		&blockchainSDK.SubnetParams{
			SubnetEVM: &blockchainSDK.SubnetEVMParams{
				ChainID:     new(big.Int).SetUint64(CHAIN_ID),
				FeeConfig:   feeConfig,
				Allocation:  allocation,
				Precompiles: precompiles,
				Timestamp:   genesisTimestamp,
			},
			Name: "TestSubnet",
		})
	if err != nil {
		log.Fatalf("❌ Failed to create subnet: %s\n", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, subnetConfig.Genesis, "", "    ")
	if err != nil {
		log.Fatalf("❌ Failed to indent genesis: %s\n", err)
	}

	if err := os.WriteFile("data/L1-genesis.json", prettyJSON.Bytes(), 0644); err != nil {
		log.Fatalf("❌ Failed to write genesis: %s\n", err)
	}

	log.Printf("✅ Successfully wrote genesis to data/L1-genesis.json\n")
}
