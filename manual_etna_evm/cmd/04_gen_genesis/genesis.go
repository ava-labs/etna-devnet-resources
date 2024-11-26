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

	_ "embed"

	"github.com/ava-labs/avalanche-cli/pkg/validatormanager"
	"github.com/ava-labs/avalanche-cli/pkg/vm"
	blockchainSDK "github.com/ava-labs/avalanche-cli/sdk/blockchain"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/coreth/utils"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
)

var (
	defaultPoAOwnerBalance  = new(big.Int).Mul(vm.OneAvax, big.NewInt(10))          // 10 Native Tokens
	defaultEVMAirdropAmount = new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil) // 10^24
)

//go:embed genesis_fuji.json
var genesisFuji []byte

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
		// FIXME: This is definitely a bug in the CLI, CLI allocates funds to a zero address here
		// It is filled in here: https://github.com/ava-labs/avalanche-cli/blob/6debe4169dce2c64352d8c9d0d0acac49e573661/pkg/vm/evm_prompts.go#L178
		ethAddr: core.GenesisAccount{Balance: defaultPoAOwnerBalance},
	}

	validatormanager.AddPoAValidatorManagerContractToAllocations(allocation)
	validatormanager.AddTransparentProxyContractToAllocations(allocation, "0x0000000000000000000000000000000000000000")

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

	err = os.WriteFile("data/genesis_fuji.json", genesisFuji, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write genesis: %s\n", err)
	}

	log.Printf("✅ Successfully wrote genesis to data/L1-genesis.json\n")
}
