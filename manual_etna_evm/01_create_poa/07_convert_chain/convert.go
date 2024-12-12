// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers/credshelper"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/key"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanchego/ids"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
)

func main() {
	exists := helpers.FileExists(helpers.ConversionIdPath)

	if exists {
		log.Println("✅ Subnet was already converted to L1")
		os.Exit(0)
	}

	chainID := helpers.LoadId(helpers.ChainIdPath)

	privKey := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
	kc := secp256k1fx.NewKeychain(privKey)

	subnetID := helpers.LoadId(helpers.SubnetIdPath)

	wallet, err := primary.MakeWallet(context.Background(), config.RPC_URL, kc, kc, primary.WalletConfig{})
	if err != nil {
		log.Fatalf("❌ Failed to initialize wallet: %s\n", err)
	}

	//TODO: replace with address.Format
	softKey, err := key.NewSoft(avagoconstants.TestnetID, key.WithPrivateKey(privKey))
	if err != nil {
		log.Fatalf("❌ Failed to create change owner address: %s\n", err)
	}

	changeOwnerAddress := softKey.P()[0]

	fmt.Printf("Using changeOwnerAddress: %s\n", changeOwnerAddress)

	subnetAuthKeys, err := address.ParseToIDs([]string{changeOwnerAddress})
	if err != nil {
		log.Fatalf("❌ Failed to parse subnet auth keys: %s\n", err)
	}

	validators := []models.SubnetValidator{}

	nodeID, proofOfPossession := credshelper.NodeInfoFromCreds(helpers.Node0KeysFolder)

	publicKey := "0x" + hex.EncodeToString(proofOfPossession.PublicKey[:])
	pop := "0x" + hex.EncodeToString(proofOfPossession.ProofOfPossession[:])

	validator := models.SubnetValidator{
		NodeID:               nodeID.String(),
		Weight:               constants.BootstrapValidatorWeight,
		Balance:              constants.BootstrapValidatorBalance,
		BLSPublicKey:         publicKey,
		BLSProofOfPossession: pop,
		ChangeOwnerAddr:      changeOwnerAddress,
	}
	validators = append(validators, validator)

	avaGoBootstrapValidators, err := blockchaincmd.ConvertToAvalancheGoSubnetValidator(validators)
	if err != nil {
		log.Fatalf("❌ Failed to convert to AvalancheGo subnet validator: %s\n", err)
	}

	managerAddress := goethereumcommon.HexToAddress(config.ProxyContractAddress)
	options := getMultisigTxOptions(subnetAuthKeys, kc)

	convertLog := fmt.Sprintf("Issuing convert subnet tx\n"+
		"subnetID: %s\n"+
		"chainID: %s\n"+
		"managerAddress: %x\n"+
		"avaGoBootstrapValidators[0]:\n"+
		"\tNodeID: %x\n"+
		"\tBLS Public Key: %x\n"+
		"\tWeight: %d\n"+
		"\tBalance: %d\n",
		subnetID.String(),
		chainID.String(),
		managerAddress[:],
		avaGoBootstrapValidators[0].NodeID[:],
		avaGoBootstrapValidators[0].Signer.PublicKey[:],
		int(avaGoBootstrapValidators[0].Weight),
		int(avaGoBootstrapValidators[0].Balance),
	)

	log.Println(convertLog)
	err = os.WriteFile("./data/convert_log.txt", []byte(convertLog), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write convert log: %s\n", err)
	}

	if len(avaGoBootstrapValidators) > 1 {
		fmt.Printf("⚠️ WARNING! Only the first validator's info is printed\n")
	}

	tx, err := wallet.P().IssueConvertSubnetToL1Tx(
		subnetID,
		chainID,
		managerAddress.Bytes(),
		avaGoBootstrapValidators,
		options...,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create convert subnet tx: %s\n", err)
	}

	helpers.SaveId(helpers.ConversionIdPath, tx.ID())

	fmt.Printf("✅ Convert subnet tx ID: %s\n", tx.ID().String())
}

func getMultisigTxOptions(subnetAuthKeys []ids.ShortID, kc *secp256k1fx.Keychain) []common.Option {
	options := []common.Option{}
	walletAddrs := kc.Addresses().List()
	changeAddr := walletAddrs[0]
	// addrs to use for signing
	customAddrsSet := set.Set[ids.ShortID]{}
	customAddrsSet.Add(walletAddrs...)
	customAddrsSet.Add(subnetAuthKeys...)
	options = append(options, common.WithCustomAddresses(customAddrsSet))
	// set change to go to wallet addr (instead of any other subnet auth key)
	changeOwner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{changeAddr},
	}
	options = append(options, common.WithChangeOwner(changeOwner))
	return options
}
