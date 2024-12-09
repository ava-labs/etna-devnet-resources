package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
)

func checkInDataFolderTempREMOVEME(path string) {
	if !strings.HasPrefix(path, "data/") {
		log.Fatalf("path must be in the data folder: %s", path)
	}
}

func FileExists(path string) (bool, error) {
	checkInDataFolderTempREMOVEME(path)
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// SaveId saves an ID to a file for the given type
func SaveId(path string, id ids.ID) error {
	checkInDataFolderTempREMOVEME(path)
	return os.WriteFile(path, []byte(id.String()), 0644)
}

// LoadId loads an ID from a file for the given type
func LoadId(path string) (ids.ID, error) {
	checkInDataFolderTempREMOVEME(path)
	text, err := os.ReadFile(path)
	if err != nil {
		return ids.ID{}, err
	}
	return ids.FromStringOrPanic(string(text)), nil
}

func SaveText(path string, text string) error {
	checkInDataFolderTempREMOVEME(path)
	return os.WriteFile(path, []byte(text), 0644)
}

func LoadText(path string) (string, error) {
	checkInDataFolderTempREMOVEME(path)
	textBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(textBytes), nil
}

func SaveUint64(path string, value uint64) error {
	checkInDataFolderTempREMOVEME(path)
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0644)
}

func LoadUint64(path string) (uint64, error) {
	checkInDataFolderTempREMOVEME(path)
	text, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(text), 10, 64)
}

func SaveHex(path string, value []byte) error {
	checkInDataFolderTempREMOVEME(path)
	return os.WriteFile(path, []byte(hex.EncodeToString(value)), 0644)
}
func LoadHex(path string) ([]byte, error) {
	checkInDataFolderTempREMOVEME(path)
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(text)))
}

func SaveNodeID(path string, nodeID ids.NodeID) error {
	checkInDataFolderTempREMOVEME(path)
	return os.WriteFile(path, []byte(nodeID.String()), 0644)
}

func LoadNodeID(path string) (ids.NodeID, error) {
	checkInDataFolderTempREMOVEME(path)
	text, err := os.ReadFile(path)
	if err != nil {
		return ids.NodeID{}, err
	}
	return ids.NodeIDFromString(string(text))
}

func SaveSecp256k1PrivateKey(path string, key *secp256k1.PrivateKey) error {
	checkInDataFolderTempREMOVEME(path)
	return SaveHex(path, key.Bytes())
}

func LoadSecp256k1PrivateKey(path string) (*secp256k1.PrivateKey, error) {
	checkInDataFolderTempREMOVEME(path)
	keyBytes, err := LoadHex(path)
	if err != nil {
		return nil, err
	}
	return secp256k1.ToPrivateKey(keyBytes)
}

func LoadSecp256k1PrivateKeyECDSA(path string) (*ecdsa.PrivateKey, error) {
	checkInDataFolderTempREMOVEME(path)
	keyBytes, err := LoadHex(path)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(keyBytes)
}
