package lib

import (
	"encoding/hex"
	"os"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

func SaveKeyToFile(key *secp256k1.PrivateKey, path string) error {
	hexStr := hex.EncodeToString(key.Bytes())
	return os.WriteFile(path, []byte(hexStr), 0644)
}

func LoadKeyFromFile(path string) (*secp256k1.PrivateKey, error) {
	hexStr, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keyBytes, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return nil, err
	}

	return secp256k1.ToPrivateKey(keyBytes)
}
