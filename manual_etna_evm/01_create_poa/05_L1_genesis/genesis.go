package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"math/big"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	_ "embed"

	"github.com/ava-labs/avalanche-cli/pkg/vm"
	pluginEVM "github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

var (
	defaultPoAOwnerBalance = new(big.Int).Mul(vm.OneAvax, big.NewInt(10)) // 10 Native Tokens
)

func main() {
	ownerKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	if err != nil {
		log.Fatalf("failed to load key from file: %s\n", err)
	}

	ethAddr := pluginEVM.PublicKeyToEthAddress(ownerKey.PublicKey())

	now := time.Now().Unix()

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

	// zeroTime := uint64(0)

	genesis := core.Genesis{
		Config: &params.ChainConfig{
			BerlinBlock:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			HomesteadBlock:      big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			FeeConfig:           feeConfig,
			ChainID:             big.NewInt(config.L1_CHAIN_ID),
		},
		Alloc: types.GenesisAlloc{
			ethAddr: {
				Balance: defaultPoAOwnerBalance,
			},
		},
		Difficulty: big.NewInt(0),
		GasLimit:   uint64(12000000),
		Timestamp:  uint64(now),
	}

	proxyAdminBytecode, err := loadHexFile("01_create_poa/04_compile_validator_manager/proxy_compiled/deployed_proxy_admin_bytecode.txt")
	if err != nil {
		log.Fatalf("❌ Failed to get proxy admin deployed bytecode: %s\n", err)
	}

	transparentProxyBytecode, err := loadHexFile("01_create_poa/04_compile_validator_manager/proxy_compiled/deployed_transparent_proxy_bytecode.txt")
	if err != nil {
		log.Fatalf("❌ Failed to get transparent proxy deployed bytecode: %s\n", err)
	}

	validatorMessagesBytecode, err := loadDeployedHexFromJSON("01_create_poa/04_compile_validator_manager/compiled/ValidatorMessages.json", nil)
	if err != nil {
		log.Fatalf("❌ Failed to get validator messages deployed bytecode: %s\n", err)
	}

	poaValidatorManagerLinkRefs := map[string]string{
		"contracts/validator-manager/ValidatorMessages.sol:ValidatorMessages": config.ValidatorMessagesAddress[2:],
	}
	poaValidatorManagerDeployedBytecode, err := loadDeployedHexFromJSON("01_create_poa/04_compile_validator_manager/compiled/PoAValidatorManager.json", poaValidatorManagerLinkRefs)
	if err != nil {
		log.Fatalf("❌ Failed to get PoA deployed bytecode: %s\n", err)
	}

	genesis.Alloc[common.HexToAddress(config.ValidatorMessagesAddress)] = types.Account{
		Code:    validatorMessagesBytecode,
		Balance: big.NewInt(0),
		Nonce:   1,
	}

	genesis.Alloc[common.HexToAddress(config.ValidatorContractAddress)] = types.Account{
		Code:    poaValidatorManagerDeployedBytecode,
		Balance: big.NewInt(0),
		Nonce:   1,
	}

	genesis.Alloc[common.HexToAddress(config.ProxyAdminContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    proxyAdminBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x0"): common.HexToHash(ethAddr.String()),
		},
	}

	genesis.Alloc[common.HexToAddress(config.ProxyContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    transparentProxyBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"): common.HexToHash(config.ValidatorContractAddress),
			common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"): common.HexToHash(config.ProxyAdminContractAddress),
		},
	}

	// Convert genesis to map to add warpConfig
	genesisMap := make(map[string]interface{})
	genesisBytes, err := json.Marshal(genesis)
	if err != nil {
		log.Fatalf("❌ Failed to marshal genesis to map: %s\n", err)
	}
	if err := json.Unmarshal(genesisBytes, &genesisMap); err != nil {
		log.Fatalf("❌ Failed to unmarshal genesis to map: %s\n", err)
	}

	// Add warpConfig to config
	configMap := genesisMap["config"].(map[string]interface{})
	configMap["warpConfig"] = map[string]interface{}{
		"blockTimestamp":               now,
		"quorumNumerator":              67,
		"requirePrimaryNetworkSigners": true,
	}

	prettyJSON, err := json.MarshalIndent(genesisMap, "", "  ")
	if err != nil {
		log.Fatalf("❌ Failed to marshal genesis: %s\n", err)
	}

	if err := helpers.SaveText(helpers.L1GenesisPath, string(prettyJSON)); err != nil {
		log.Fatalf("❌ Failed to write genesis: %s\n", err)
	}

	log.Printf("✅ Successfully wrote genesis to data/L1-genesis.json\n")
}
