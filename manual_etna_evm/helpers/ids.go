package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/ava-labs/avalanchego/ids"
)

// getPath constructs the file path for a given ID type
func getPath(idType string) string {
	return fmt.Sprintf("data/%s.txt", idType)
}

// IdFileExists checks if an ID file exists for the given type
func IdFileExists(idType string) (bool, error) {
	return TextFileExists(idType)
}

// SaveId saves an ID to a file for the given type
func SaveId(idType string, id ids.ID) error {
	return SaveText(idType, id.String())
}

// LoadId loads an ID from a file for the given type
func LoadId(idType string) (ids.ID, error) {
	text, err := LoadText(idType)
	if err != nil {
		return ids.ID{}, err
	}
	return ids.FromStringOrPanic(text), nil
}

// TextFileExists checks if a text file exists for the given type
func TextFileExists(textType string) (bool, error) {
	path := getPath(textType)
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// SaveText saves a text string to a file for the given type
func SaveText(textType string, text string) error {
	return os.WriteFile(getPath(textType), []byte(text), 0644)
}

// LoadText loads a text string from a file for the given type
func LoadText(textType string) (string, error) {
	textBytes, err := os.ReadFile(getPath(textType))
	if err != nil {
		return "", err
	}
	return string(textBytes), nil
}
