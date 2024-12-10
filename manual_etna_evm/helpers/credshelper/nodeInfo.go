package credshelper

import (
	"encoding/pem"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func NodeInfoFromCreds(folder string) (ids.NodeID, *signer.ProofOfPossession) {
	blsKey := helpers.LoadBLSKey(folder + "signer.key")
	pop := signer.NewProofOfPossession(blsKey)
	certString := helpers.LoadText(folder + "staker.crt")

	block, _ := pem.Decode([]byte(certString))
	if block == nil || block.Type != "CERTIFICATE" {
		panic("failed to decode PEM block containing certificate")
	}

	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to decode PEM block containing certificate")
	}

	nodeID := ids.NodeIDFromCert(cert)

	return nodeID, pop
}
