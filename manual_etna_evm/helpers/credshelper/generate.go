package credshelper

import (
	"encoding/base64"
	"fmt"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func GetCredsBase64(folder string) (string, string, string, error) {
	stakerKey, err := helpers.LoadBytes(folder + "staker.key")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load staker key: %w", err)
	}
	stakerCert, err := helpers.LoadBytes(folder + "staker.crt")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load staker cert: %w", err)
	}
	signerKey, err := helpers.LoadBytes(folder + "signer.key")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load signer key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(stakerKey),
		base64.StdEncoding.EncodeToString(stakerCert),
		base64.StdEncoding.EncodeToString(signerKey), nil
}
