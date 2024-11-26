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
	path := getPath(idType)
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
func SaveId(idType string, id ids.ID) error {
	return os.WriteFile(getPath(idType), []byte(id.String()), 0644)
}

// LoadId loads an ID from a file for the given type
func LoadId(idType string) (ids.ID, error) {
	idBytes, err := os.ReadFile(getPath(idType))
	if err != nil {
		return ids.ID{}, err
	}
	return ids.FromStringOrPanic(string(idBytes)), nil
}
