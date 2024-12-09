package credshelper

import (
	"encoding/pem"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func NodeInfoFromCreds(folder string) (ids.NodeID, *signer.ProofOfPossession, error) {
	blsKey, err := helpers.LoadBLSKey(folder + "signer.key")
	if err != nil {
		return ids.NodeID{}, nil, err
	}

	pop := signer.NewProofOfPossession(blsKey)

	certString, err := helpers.LoadText(folder + "staker.crt")
	if err != nil {
		return ids.NodeID{}, nil, err
	}

	block, _ := pem.Decode([]byte(certString))
	if block == nil || block.Type != "CERTIFICATE" {
		return ids.NodeID{}, nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		return ids.NodeID{}, nil, err
	}

	nodeID := ids.NodeIDFromCert(cert)

	return nodeID, pop, nil

}
