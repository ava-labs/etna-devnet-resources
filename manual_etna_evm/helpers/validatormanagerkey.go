package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"

	goethereumcrypto "github.com/ethereum/go-ethereum/crypto"
)

func ValidatorManagerKeyExists() (bool, error) {
	return TextFileExists("validator_manager_owner_key")
}

func GenerateValidatorManagerKeyAndSave() error {
	key, err := secp256k1.NewPrivateKey()
	if err != nil {
		return err
	}
	hexStr := hex.EncodeToString(key.Bytes())
	return SaveText("validator_manager_owner_key", hexStr)
}

func LoadValidatorManagerKey() (*secp256k1.PrivateKey, error) {
	hexStr, err := LoadText("validator_manager_owner_key")
	if err != nil {
		return nil, err
	}

	keyBytes, err := hex.DecodeString(strings.TrimSpace(hexStr))
	if err != nil {
		return nil, err
	}

	return secp256k1.ToPrivateKey(keyBytes)
}

func LoadValidatorManagerKeyECDSA() (*ecdsa.PrivateKey, error) {
	hexStr, err := LoadText("validator_manager_owner_key")
	if err != nil {
		return nil, err
	}

	keyBytes, err := hex.DecodeString(strings.TrimSpace(hexStr))
	if err != nil {
		return nil, err
	}

	return goethereumcrypto.ToECDSA(keyBytes)
}
