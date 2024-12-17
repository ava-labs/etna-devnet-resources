package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func GetValidatorCMD(credsFolder string, nodeIndex int) (string, error) {
	if nodeIndex == 0 {
		return "", fmt.Errorf("node index cannot be 0")
	}

	containerName := fmt.Sprintf("node%d", nodeIndex)
	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return "", fmt.Errorf("failed to load subnet id: %w", err)
	}

	stakerKeyBase64, stakerCertBase64, signerKeyBase64, err := GetCredsBase64(credsFolder)
	if err != nil {
		return "", fmt.Errorf("failed to get creds base64: %w", err)
	}

	httpPort := 9650 + (nodeIndex)*2
	stakingPort := httpPort + 1

	script := fmt.Sprintf(`
docker rm -f %s || true; \
docker run -d \
  --name %s \
  --network host \
  -e AVALANCHEGO_NETWORK_ID=fuji \
  -e AVALANCHEGO_HTTP_PORT=%d \
  -e AVALANCHEGO_STAKING_PORT=%d \
  -e AVALANCHEGO_TRACK_SUBNETS=%s \
  -e AVALANCHEGO_HTTP_ALLOWED_HOSTS=* \
  -e AVALANCHEGO_HTTP_HOST=0.0.0.0 \
  -e AVALANCHEGO_STAKING_TLS_CERT_FILE_CONTENT=%s \
  -e AVALANCHEGO_STAKING_TLS_KEY_FILE_CONTENT=%s \
  -e BLS_KEY_BASE64=%s \
  -e AVALANCHEGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme \
  -e AVALANCHEGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \
  containerman17/avalanchego-subnetevm:v1.12.0_v0.7.0 ;\

	`, containerName, containerName, httpPort, stakingPort, subnetID.String(), stakerCertBase64, stakerKeyBase64, signerKeyBase64, containerName)

	return script, nil
}

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
