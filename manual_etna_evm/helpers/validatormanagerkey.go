package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"

	goethereumcrypto "github.com/ethereum/go-ethereum/crypto"
)

const VALIDATOR_MANAGER_KEY_PATH = "data/validator_manager_owner_key.txt"

func ValidatorManagerKeyExists() (bool, error) {
	_, err := os.Stat(VALIDATOR_MANAGER_KEY_PATH)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateValidatorManagerKeyAndSave() error {
	key, err := secp256k1.NewPrivateKey()
	if err != nil {
		return err
	}
	hexStr := hex.EncodeToString(key.Bytes())
	return os.WriteFile(VALIDATOR_MANAGER_KEY_PATH, []byte(hexStr), 0644)
}

func LoadValidatorManagerKey() (*secp256k1.PrivateKey, error) {
	hexStr, err := os.ReadFile(VALIDATOR_MANAGER_KEY_PATH)
	if err != nil {
		return nil, err
	}

	hexStr = []byte(strings.TrimSpace(string(hexStr)))

	keyBytes, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return nil, err
	}

	return secp256k1.ToPrivateKey(keyBytes)
}

func LoadValidatorManagerKeyECDSA() (*ecdsa.PrivateKey, error) {
	keyHex, err := os.ReadFile(VALIDATOR_MANAGER_KEY_PATH)
	if err != nil {
		return nil, err
	}
	keyBytes, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
	if err != nil {
		return nil, err
	}

	key, err := goethereumcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
