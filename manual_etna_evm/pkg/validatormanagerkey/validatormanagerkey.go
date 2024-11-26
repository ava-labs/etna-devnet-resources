package validatormanagerkey

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

const PATH = "data/validator_manager_owner_key.txt"

func Exists() (bool, error) {
	_, err := os.Stat(PATH)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateAndSave() error {
	key, err := secp256k1.NewPrivateKey()
	if err != nil {
		return err
	}
	hexStr := hex.EncodeToString(key.Bytes())
	return os.WriteFile(PATH, []byte(hexStr), 0644)
}

func Load() (*secp256k1.PrivateKey, error) {
	hexStr, err := os.ReadFile(PATH)
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
