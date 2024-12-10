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
	"github.com/ethereum/go-ethereum/crypto"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		panic(fmt.Errorf("checking if file exists at %s: %w", path, err))
	}
	return true
}

// SaveId saves an ID to a file for the given type
func SaveId(path string, id ids.ID) {
	if err := os.WriteFile(path, []byte(id.String()), 0644); err != nil {
		panic(fmt.Errorf("saving ID to %s: %w", path, err))
	}
}

// LoadId loads an ID from a file for the given type
func LoadId(path string) ids.ID {
	return ids.FromStringOrPanic(
		LoadText(path),
	)
}

func SaveText(path string, text string) {
	SaveBytes(path, []byte(text))
}

func LoadText(path string) string {
	textBytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading text from %s: %w", path, err))
	}
	return strings.TrimSpace(string(textBytes))
}

func SaveBytes(path string, value []byte) {

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Errorf("creating directory %s: %w", dir, err))
	}

	if err := os.WriteFile(path, value, 0644); err != nil {
		panic(fmt.Errorf("saving bytes to %s: %w", path, err))
	}
}

func LoadBytes(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading bytes from %s: %w", path, err))
	}
	return bytes
}

func SaveUint64(path string, value uint64) {
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0644); err != nil {
		panic(fmt.Errorf("saving uint64 to %s: %w", path, err))
	}
}

func LoadUint64(path string) uint64 {
	text, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading uint64 from %s: %w", path, err))
	}
	val, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		panic(fmt.Errorf("parsing uint64 from %s: %w", path, err))
	}
	return val
}

func SaveHex(path string, value []byte) {
	SaveText(path, hex.EncodeToString(value))
}

func LoadHex(path string) []byte {
	text := LoadText(path)
	bytes, err := hex.DecodeString(text)
	if err != nil {
		panic(fmt.Errorf("decoding hex from %s: %w", path, err))
	}
	return bytes
}

func SaveNodeID(path string, nodeID ids.NodeID) {
	if err := os.WriteFile(path, []byte(nodeID.String()), 0644); err != nil {
		panic(fmt.Errorf("saving node ID to %s: %w", path, err))
	}
}

func LoadNodeID(path string) ids.NodeID {
	text, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading node ID from %s: %w", path, err))
	}
	nodeID, err := ids.NodeIDFromString(string(text))
	if err != nil {
		panic(fmt.Errorf("parsing node ID from %s: %w", path, err))
	}
	return nodeID
}

func SaveSecp256k1PrivateKey(path string, key *secp256k1.PrivateKey) {
	SaveHex(path, key.Bytes())
}

func LoadSecp256k1PrivateKey(path string) *secp256k1.PrivateKey {
	keyBytes := LoadHex(path)
	key, err := secp256k1.ToPrivateKey(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing secp256k1 private key from %s: %w", path, err))
	}
	return key
}

func LoadSecp256k1PrivateKeyECDSA(path string) *ecdsa.PrivateKey {
	keyBytes := LoadHex(path)
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing ECDSA private key from %s: %w", path, err))
	}
	return key
}

func SaveBLSKey(path string, key *bls.SecretKey) {
	SaveBytes(path, bls.SecretKeyToBytes(key))
}

func LoadBLSKey(path string) *bls.SecretKey {
	keyBytes := LoadBytes(path)
	key, err := bls.SecretKeyFromBytes(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing BLS key from %s: %w", path, err))
	}
	return key
}
