package main

import (
	"fmt"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"
)

func main() {
	nodeID, _ := credshelper.NodeInfoFromCreds(helpers.AddValidatorKeysFolder)
	containerName := nodeID.String()
	subnetID := helpers.LoadId(helpers.SubnetIdPath)

	stakerKeyBase64, stakerCertBase64, signerKeyBase64 := credshelper.GetCredsBase64(helpers.AddValidatorKeysFolder)

	script := fmt.Sprintf(`
docker rm -f %s || true; \
docker run -d \
  --name %s \
  --network host \
  -e AVALANCHEGO_NETWORK_ID=fuji \
  -e AVALANCHEGO_HTTP_PORT=9652 \
  -e AVALANCHEGO_STAKING_PORT=9653 \
  -e AVALANCHEGO_TRACK_SUBNETS=%s \
  -e AVALANCHEGO_HTTP_ALLOWED_HOSTS=* \
  -e AVALANCHEGO_HTTP_HOST=0.0.0.0 \
  -e AVALANCHEGO_STAKING_TLS_CERT_FILE_CONTENT=%s \
  -e AVALANCHEGO_STAKING_TLS_KEY_FILE_CONTENT=%s \
  -e BLS_KEY_BASE64=%s \
  -e AVALANCHEGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme \
  containerman17/avalanchego-subnetevm:v1.12.0_v0.6.12 ;\
  docker logs -f %s
	`, containerName, containerName, subnetID.String(), stakerCertBase64, stakerKeyBase64, signerKeyBase64, containerName)

	fmt.Println(script)
}
