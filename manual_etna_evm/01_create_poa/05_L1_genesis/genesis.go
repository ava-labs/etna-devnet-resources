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

	_ "embed"

	"github.com/ava-labs/avalanche-cli/pkg/vm"
	pluginEVM "github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	defaultPoAOwnerBalance = new(big.Int).Mul(vm.OneAvax, big.NewInt(10)) // 10 Native Tokens
)

//go:embed proxy_compiled/deployed_proxy_admin_bytecode.txt
var proxyAdminBytecode []byte

//go:embed proxy_compiled/deployed_transparent_proxy_bytecode.txt
var transparentProxyBytecode []byte

func main() {

	ownerKey := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)

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

	genesis.Alloc[common.HexToAddress(config.ProxyAdminContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    proxyAdminBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x0"): common.HexToHash(ethAddr.String()),
		},
	}

	expectedValidatorManagerAddress := DeriveContractAddress(ethAddr, 1)
	log.Printf("🔍 Derived contract address: %s\n", expectedValidatorManagerAddress.Hex())

	genesis.Alloc[common.HexToAddress(config.ProxyContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    transparentProxyBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			//address 0 is would be ValidatorMessages, 1 - ValidatorManager
			common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"): common.BytesToHash(expectedValidatorManagerAddress.Bytes()),
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

	helpers.SaveText(helpers.L1GenesisPath, string(prettyJSON))

	log.Printf("✅ Successfully wrote genesis to data/L1-genesis.json\n")
}

func DeriveContractAddress(from common.Address, nonce uint64) common.Address {
	encoded, err := rlp.EncodeToBytes([]interface{}{from, nonce})
	if err != nil {
		panic(err)
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToAddress(hash[12:])
}
