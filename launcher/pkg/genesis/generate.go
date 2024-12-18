package genesis

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	_ "embed"

	_ "embed"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:embed proxy_compiled/deployed_proxy_admin_bytecode.txt
var proxyAdminBytecodeHexString string

//go:embed proxy_compiled/deployed_transparent_proxy_bytecode.txt
var transparentProxyBytecodeHexString string

var (
	OneAvax                = new(big.Int).SetUint64(1000000000000000000)
	defaultPoAOwnerBalance = new(big.Int).Mul(OneAvax, big.NewInt(10)) // 10 Native Tokens

)

const (
	ProxyContractAddress      = "0xFEEDC0DE0000000000000000000000000000000"
	ProxyAdminContractAddress = "0xC0FFEE1234567890aBcDEF1234567890AbCdEf34"
	ZeroAddress               = "0x0000000000000000000000000000000000000000"
)

type GeneratePayload struct {
	OwnerEthAddressString string
	EvmChainId            int
}

func Generate(payload GeneratePayload) (string, error) {
	if !common.IsHexAddress(payload.OwnerEthAddressString) {
		return "", fmt.Errorf("invalid owner address: %s\n", payload.OwnerEthAddressString)
	}
	ownerAddress := common.HexToAddress(payload.OwnerEthAddressString)

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
			ChainID:             big.NewInt(int64(payload.EvmChainId)),
		},
		Alloc: types.GenesisAlloc{
			ownerAddress: types.Account{
				Balance: defaultPoAOwnerBalance,
			},
		},
		Difficulty: big.NewInt(0),
		GasLimit:   uint64(12000000),
		Timestamp:  uint64(now),
	}

	proxyAdminBytecodeHexString = strings.TrimSpace(strings.TrimPrefix(proxyAdminBytecodeHexString, "0x"))
	transparentProxyBytecodeHexString = strings.TrimSpace(strings.TrimPrefix(transparentProxyBytecodeHexString, "0x"))

	proxyAdminBytecode, err := hex.DecodeString(proxyAdminBytecodeHexString)
	if err != nil {
		return "", fmt.Errorf("failed to decode proxy admin bytecode: %s\n", err)
	}

	transparentProxyBytecode, err := hex.DecodeString(transparentProxyBytecodeHexString)
	if err != nil {
		return "", fmt.Errorf("failed to decode transparent proxy bytecode: %s\n", err)
	}

	genesis.Alloc[common.HexToAddress(ProxyAdminContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    proxyAdminBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x0"): common.HexToHash(ownerAddress.String()),
		},
	}

	genesis.Alloc[common.HexToAddress(ProxyContractAddress)] = types.Account{
		Balance: big.NewInt(0),
		Code:    transparentProxyBytecode,
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"): common.HexToHash(ZeroAddress),
			common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"): common.HexToHash(ProxyAdminContractAddress),
		},
	}

	// Convert genesis to map to add warpConfig
	genesisMap := make(map[string]interface{})
	genesisBytes, err := json.Marshal(genesis)
	if err != nil {
		return "", fmt.Errorf("failed to marshal genesis to map: %s\n", err)
	}
	if err := json.Unmarshal(genesisBytes, &genesisMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal genesis to map: %s\n", err)
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
		return "", fmt.Errorf("failed to marshal genesis: %s\n", err)
	}

	return string(prettyJSON), nil
}

func MustDeriveContractAddress(from common.Address, nonce uint64) common.Address {
	encoded, err := rlp.EncodeToBytes([]interface{}{from, nonce})
	if err != nil {
		panic(err)
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToAddress(hash[12:])
}
