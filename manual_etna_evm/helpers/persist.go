package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking if file exists at %s: %w", path, err)
	}
	return true, nil
}

// SaveId saves an ID to a file for the given type
func SaveId(path string, id ids.ID) error {
	if err := os.WriteFile(path, []byte(id.String()), 0644); err != nil {
		return fmt.Errorf("saving ID to %s: %w", path, err)
	}
	return nil
}

// LoadId loads an ID from a file for the given type
func LoadId(path string) (ids.ID, error) {
	text, err := LoadText(path)
	if err != nil {
		return ids.ID{}, err
	}
	id, err := ids.FromString(text)
	if err != nil {
		return ids.ID{}, fmt.Errorf("parsing ID from %s: %w", path, err)
	}
	return id, nil
}

func SaveText(path string, text string) error {
	return SaveBytes(path, []byte(text))
}

func LoadText(path string) (string, error) {
	textBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("loading text from %s: %w", path, err)
	}
	return strings.TrimSpace(string(textBytes)), nil
}

func SaveBytes(path string, value []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, value, 0644); err != nil {
		return fmt.Errorf("saving bytes to %s: %w", path, err)
	}

	return nil
}

func LoadBytes(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading bytes from %s: %w", path, err)
	}
	return bytes, nil
}

func SaveUint64(path string, value uint64) error {
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0644); err != nil {
		return fmt.Errorf("saving uint64 to %s: %w", path, err)
	}
	return nil
}

func LoadUint64(path string) (uint64, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("loading uint64 from %s: %w", path, err)
	}
	val, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing uint64 from %s: %w", path, err)
	}
	return val, nil
}

func SaveHex(path string, value []byte) error {
	return SaveText(path, hex.EncodeToString(value))
}

func LoadHex(path string) ([]byte, error) {
	text, err := LoadText(path)
	if err != nil {
		return nil, err
	}
	bytes, err := hex.DecodeString(text)
	if err != nil {
		return nil, fmt.Errorf("decoding hex from %s: %w", path, err)
	}
	return bytes, nil
}

func SaveAddress(path string, address common.Address) error {
	return SaveHex(path, address[:])
}

func LoadAddress(path string) (common.Address, error) {
	bytes, err := LoadHex(path)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(bytes), nil
}

func SaveNodeID(path string, nodeID ids.NodeID) error {
	if err := os.WriteFile(path, []byte(nodeID.String()), 0644); err != nil {
		return fmt.Errorf("saving node ID to %s: %w", path, err)
	}
	return nil
}

func LoadNodeID(path string) (ids.NodeID, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return ids.NodeID{}, fmt.Errorf("loading node ID from %s: %w", path, err)
	}
	nodeID, err := ids.NodeIDFromString(string(text))
	if err != nil {
		return ids.NodeID{}, fmt.Errorf("parsing node ID from %s: %w", path, err)
	}
	return nodeID, nil
}

func SaveSecp256k1PrivateKey(path string, key *secp256k1.PrivateKey) error {
	return SaveHex(path, key.Bytes())
}

func LoadSecp256k1PrivateKey(path string) (*secp256k1.PrivateKey, error) {
	keyBytes, err := LoadHex(path)
	if err != nil {
		return nil, err
	}
	key, err := secp256k1.ToPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing secp256k1 private key from %s: %w", path, err)
	}
	return key, nil
}

func LoadSecp256k1PrivateKeyECDSA(path string) (*ecdsa.PrivateKey, error) {
	keyBytes, err := LoadHex(path)
	if err != nil {
		return nil, err
	}
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing ECDSA private key from %s: %w", path, err)
	}
	return key, nil
}

func SaveBLSKey(path string, key *bls.SecretKey) error {
	return SaveBytes(path, bls.SecretKeyToBytes(key))
}

func LoadBLSKey(path string) (*bls.SecretKey, error) {
	keyBytes, err := LoadBytes(path)
	if err != nil {
		return nil, err
	}
	key, err := bls.SecretKeyFromBytes(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing BLS key from %s: %w", path, err)
	}
	return key, nil
}

func CopyFile(src, dst string) error {
	srcBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	return os.WriteFile(dst, srcBytes, 0644)
}
