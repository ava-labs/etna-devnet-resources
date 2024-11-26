package datafiles

import (
	"errors"
	"os"

	"github.com/ava-labs/avalanchego/ids"
)

const SUBNET_ID_PATH = "data/subnet.txt"

func SubnetIDExists() (bool, error) {
	_, err := os.Stat(SUBNET_ID_PATH)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func SaveSubnetID(subnetID ids.ID) error {
	return os.WriteFile(SUBNET_ID_PATH, []byte(subnetID.String()), 0644)
}
