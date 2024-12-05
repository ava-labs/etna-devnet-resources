package main

import (
	"encoding/hex"
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
