package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
)

func checkInDataFolderTempREMOVEME(path string) {
	if !strings.HasPrefix(path, "data/") {
		log.Fatalf("path must be in the data folder: %s", path)
	}
}

func FileExists(path string) bool {
	checkInDataFolderTempREMOVEME(path)
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
	checkInDataFolderTempREMOVEME(path)
	if err := os.WriteFile(path, []byte(id.String()), 0644); err != nil {
		panic(fmt.Errorf("saving ID to %s: %w", path, err))
	}
}

// LoadId loads an ID from a file for the given type
func LoadId(path string) ids.ID {
	checkInDataFolderTempREMOVEME(path)
	text, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading ID from %s: %w", path, err))
	}
	return ids.FromStringOrPanic(string(text))
}

func SaveText(path string, text string) {
	checkInDataFolderTempREMOVEME(path)
	SaveBytes(path, []byte(text))
}

func LoadText(path string) string {
	checkInDataFolderTempREMOVEME(path)
	textBytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading text from %s: %w", path, err))
	}
	return string(textBytes)
}

func SaveBytes(path string, value []byte) {
	checkInDataFolderTempREMOVEME(path)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Errorf("creating directory %s: %w", dir, err))
	}

	if err := os.WriteFile(path, value, 0644); err != nil {
		panic(fmt.Errorf("saving bytes to %s: %w", path, err))
	}
}

func LoadBytes(path string) []byte {
	checkInDataFolderTempREMOVEME(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("loading bytes from %s: %w", path, err))
	}
	return bytes
}

func SaveUint64(path string, value uint64) {
	checkInDataFolderTempREMOVEME(path)
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0644); err != nil {
		panic(fmt.Errorf("saving uint64 to %s: %w", path, err))
	}
}

func LoadUint64(path string) uint64 {
	checkInDataFolderTempREMOVEME(path)
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
	checkInDataFolderTempREMOVEME(path)
	SaveText(path, hex.EncodeToString(value))
}

func LoadHex(path string) []byte {
	checkInDataFolderTempREMOVEME(path)
	text := LoadText(path)
	bytes, err := hex.DecodeString(strings.TrimSpace(string(text)))
	if err != nil {
		panic(fmt.Errorf("decoding hex from %s: %w", path, err))
	}
	return bytes
}

func SaveNodeID(path string, nodeID ids.NodeID) {
	checkInDataFolderTempREMOVEME(path)
	if err := os.WriteFile(path, []byte(nodeID.String()), 0644); err != nil {
		panic(fmt.Errorf("saving node ID to %s: %w", path, err))
	}
}

func LoadNodeID(path string) ids.NodeID {
	checkInDataFolderTempREMOVEME(path)
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
	checkInDataFolderTempREMOVEME(path)
	SaveHex(path, key.Bytes())
}

func LoadSecp256k1PrivateKey(path string) *secp256k1.PrivateKey {
	checkInDataFolderTempREMOVEME(path)
	keyBytes := LoadHex(path)
	key, err := secp256k1.ToPrivateKey(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing secp256k1 private key from %s: %w", path, err))
	}
	return key
}

func LoadSecp256k1PrivateKeyECDSA(path string) *ecdsa.PrivateKey {
	checkInDataFolderTempREMOVEME(path)
	keyBytes := LoadHex(path)
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing ECDSA private key from %s: %w", path, err))
	}
	return key
}

func SaveBLSKey(path string, key *bls.SecretKey) {
	checkInDataFolderTempREMOVEME(path)
	SaveBytes(path, bls.SecretKeyToBytes(key))
}

func LoadBLSKey(path string) *bls.SecretKey {
	checkInDataFolderTempREMOVEME(path)
	keyBytes := LoadBytes(path)
	key, err := bls.SecretKeyFromBytes(keyBytes)
	if err != nil {
		panic(fmt.Errorf("parsing BLS key from %s: %w", path, err))
	}
	return key
}
