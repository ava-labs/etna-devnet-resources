package main

import (
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
)

func main() {
	exists, err := folderExists("data/new_validator")
	if err != nil {
		log.Fatalf("Failed to check if data/new_validator directory exists: %s", err)
	}
	if exists {
		log.Fatalf("data/new_validator directory already exists. Please delete it manually before generating new keys")
	}

	err = os.MkdirAll("data/new_validator/staking", 0755)
	if err != nil {
		log.Fatalf("Failed to create data/new_validator/staking directory: %s", err)
	}

	// Generate signer (BLS) key
	signerKey, err := bls.NewSecretKey()
	if err != nil {
		log.Fatalf("Failed to generate signer key: %s", err)
	}

	// Generate staking certificate and key
	stakerCert, stakerKey, err := staking.NewCertAndKeyBytes()
	if err != nil {
		log.Fatalf("Failed to generate staking keys: %s", err)
	}

	// Write keys and certificate to files
	err = os.WriteFile("data/new_validator/staking/signer.key", signerKey.Serialize(), 0644)
	if err != nil {
		log.Fatalf("Failed to write signer key: %s", err)
	}

	err = os.WriteFile("data/new_validator/staking/staker.crt", stakerCert, 0644)
	if err != nil {
		log.Fatalf("Failed to write staker certificate: %s", err)
	}

	err = os.WriteFile("data/new_validator/staking/staker.key", stakerKey, 0644)
	if err != nil {
		log.Fatalf("Failed to write staker key: %s", err)
	}

	// Generate Node ID from certificate
	block, _ := pem.Decode(stakerCert)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("Failed to decode PEM block containing certificate")
	}

	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %s", err)
	}

	nodeID := ids.NodeIDFromCert(cert)

	// Write Node ID to file
	err = helpers.SaveNodeID("new_validator/nodeId", nodeID)
	if err != nil {
		log.Fatalf("Failed to save node ID: %s", err)
	}

	fmt.Printf("Node ID generated: %s\n", nodeID.String())

	// Generate proof of possession
	pop := signer.NewProofOfPossession(signerKey)

	// Write proof of possession to file
	err = helpers.SaveProofOfPossession("data/new_validator/pop.json", pop)
	if err != nil {
		log.Fatalf("Failed to write proof of possession: %s", err)
	}

	fmt.Printf("Proof of Possession generated: %s\n", hex.EncodeToString(pop.ProofOfPossession[:]))
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
