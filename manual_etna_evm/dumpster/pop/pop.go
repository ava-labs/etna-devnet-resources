package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	_ "embed"

	"encoding/pem"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
)

//I know that 3 files below are somehow become nodeid, publickey and pop. Your goal is to do some magic to make them become nodeid, publickey and pop and print them out
//On the app-level the task is to generate a node pop, pubkey and nodeID before launching it.
//RN you have to launch the node which generates them

//go:embed signer.key
var testKeyBin []byte

//go:embed staker.crt
var stakerCertX509 string

//go:embed staker.key
var stakerKeyX509 string

func main() {

	infoClient := info.NewClient("http://localhost:9650")

	nodeIDoriginal, poporiginal, err := infoClient.GetNodeID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Generate a new BLS secret key
	sk, err := bls.SecretKeyFromBytes(testKeyBin)
	if err != nil {
		log.Fatal(err)
	}

	// Create proof of possession from secret key
	pop := signer.NewProofOfPossession(sk)

	// Get public key
	pk := bls.PublicFromSecretKey(sk)
	pkBytes := bls.PublicKeyToCompressedBytes(pk)

	fmt.Println("--------------------------------")
	fmt.Println("PublicKey generated:", hex.EncodeToString(pkBytes))
	fmt.Println("PublicKey original:", hex.EncodeToString(poporiginal.PublicKey[:]))

	fmt.Print("\n")
	fmt.Println("ProofOfPossession original:", hex.EncodeToString(poporiginal.ProofOfPossession[:]))
	fmt.Println("ProofOfPossession generated:", hex.EncodeToString(pop.ProofOfPossession[:]))
	fmt.Print("\n")

	// Parse PEM block
	block, _ := pem.Decode([]byte(stakerCertX509))
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("failed to decode PEM block containing certificate")
	}

	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	nodeID := ids.NodeIDFromCert(cert)

	fmt.Printf("NodeID generated: %s\n", nodeID.String())
	fmt.Println("NodeID original:", nodeIDoriginal.String())

	// fmt.Println("\nEmbedded Values:")
	// fmt.Println("TestKeyBin:", hex.EncodeToString(testKeyBin))
	// fmt.Println("StakerCertX509:", stakerCertX509)
	// fmt.Println("StakerKeyX509:", stakerKeyX509)
}
