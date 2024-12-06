package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

type linkReference struct {
	Start  int `json:"start"`
	Length int `json:"length"`
}

type compiledJSON struct {
	DeployedBytecode struct {
		Object string `json:"object"`
	} `json:"deployedBytecode"`
	LinkReferences map[string]map[string][]linkReference `json:"linkReferences"`
}

func laodDeployedContractBytecodeFromJSON(path string) ([]byte, map[string]map[string][]linkReference, error) {
	compiled := compiledJSON{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(data, &compiled)
	if err != nil {
		return nil, nil, err
	}

	deployedBytecodeHex := compiled.DeployedBytecode.Object
	if strings.HasPrefix(deployedBytecodeHex, "0x") {
		deployedBytecodeHex = deployedBytecodeHex[2:]
	}

	bytecode, err := hex.DecodeString(deployedBytecodeHex)
	if err != nil {
		return nil, nil, err
	}
	return bytecode, compiled.LinkReferences, nil
}

func insertLinkReferences(bytecode []byte, linkReferences map[string]map[string][]linkReference, addressMap map[string]map[string]string) ([]byte, error) {
	// Make a copy of the bytecode to modify
	result := make([]byte, len(bytecode))
	copy(result, bytecode)

	// Iterate through all link references
	for file, contracts := range linkReferences {
		for contract, refs := range contracts {
			// Get the address from the address map
			addresses, ok := addressMap[file]
			if !ok {
				return nil, fmt.Errorf("missing address map entry for file: %s", file)
			}
			address, ok := addresses[contract]
			if !ok {
				return nil, fmt.Errorf("missing address for contract: %s in file: %s", contract, file)
			}

			// Convert address to bytes
			addressBytes, err := hex.DecodeString(strings.TrimPrefix(address, "0x"))
			if err != nil {
				return nil, fmt.Errorf("invalid address hex: %s", address)
			}
			if len(addressBytes) != 20 {
				return nil, fmt.Errorf("address must be 20 bytes long, got %d bytes", len(addressBytes))
			}

			// Insert address at each reference location
			for _, ref := range refs {
				if ref.Length != 20 {
					return nil, fmt.Errorf("link reference length must be 20, got %d", ref.Length)
				}
				copy(result[ref.Start:ref.Start+ref.Length], addressBytes)
			}
		}
	}

	return result, nil
}
