package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func loadHexFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Handle 0x prefix if present
	if len(data) > 1 && data[0] == '0' && data[1] == 'x' {
		data = data[2:]
	}
	// Trim whitespace and newlines
	cleanData := []byte(strings.TrimSpace(string(data)))
	return hex.DecodeString(string(cleanData))
}

type compiledJSON struct {
	DeployedBytecode struct {
		Object string `json:"object"`
	} `json:"deployedBytecode"`
}

func loadDeployedHexFromJSON(path string, linkReferences map[string]string) ([]byte, error) {
	compiled := compiledJSON{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &compiled)
	if err != nil {
		return nil, err
	}

	resultHex := compiled.DeployedBytecode.Object

	if linkReferences != nil {
		for refName, address := range linkReferences {
			if len(address) != 40 {
				return nil, fmt.Errorf("invalid placeholder length %d, expected 40: %s", len(address), address)
			}
			if _, err := hex.DecodeString(address); err != nil {
				return nil, fmt.Errorf("invalid hex in placeholder address: %s", address)
			}

			linkRefHash := crypto.Keccak256Hash([]byte(refName))
			linkRefHashStr := linkRefHash.Hex()
			placeholderStr := fmt.Sprintf("__$%s$__", linkRefHashStr[2:36])

			fmt.Printf("Replacing %s with %s\n", placeholderStr, address)

			resultHex = strings.Replace(resultHex, placeholderStr, address, -1)
		}
	}

	if strings.Contains(resultHex, "$__") {
		return nil, fmt.Errorf("unresolved link reference found in bytecode: %s", resultHex)
	}

	// Handle 0x prefix if present
	if len(resultHex) > 1 && resultHex[0] == '0' && resultHex[1] == 'x' {
		resultHex = resultHex[2:]
	}

	return hex.DecodeString(resultHex)
}
