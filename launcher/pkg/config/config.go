package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

func GetRPCUrl() string {
	return "https://api.avax-test.network"
}

var privateKeyFolder = "."

func init() {
	if os.Getenv("PRIV_KEY_FOLDER") != "" {
		privateKeyFolder = os.Getenv("PRIV_KEY_FOLDER")
	}
}

var privKeyCache *secp256k1.PrivateKey = nil

func LoadOrGeneratePrivateKey() *secp256k1.PrivateKey {
	if privKeyCache == nil {
		// Try to load private key from file
		keyBytes, err := os.ReadFile(privateKeyFolder + "/priv.hex")
		if err == nil {
			// File exists, try to load key
			text := strings.TrimPrefix(strings.TrimSpace(string(keyBytes)), "0x")
			bytes, err := hex.DecodeString(text)
			if err != nil {
				panic(fmt.Errorf("failed to decode private key hex: %w", err))
			}
			key, err := secp256k1.ToPrivateKey(bytes)
			if err != nil {
				panic(fmt.Errorf("failed to parse private key: %w", err))
			}
			privKeyCache = key
			return key
		}

		// File doesn't exist, generate new key
		key, err := secp256k1.NewPrivateKey()
		if err != nil {
			panic(fmt.Errorf("failed to generate private key: %w", err))
		}

		// Save key to file
		keyHex := hex.EncodeToString(key.Bytes())
		err = os.WriteFile("./priv.hex", []byte(keyHex), 0600)
		if err != nil {
			panic(fmt.Errorf("failed to save private key: %w", err))
		}

		privKeyCache = key
	}

	return privKeyCache
}
