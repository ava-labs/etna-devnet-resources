// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/key"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanche-cli/pkg/txutils"
	"github.com/ava-labs/avalanche-cli/pkg/utils"
	"github.com/ava-labs/avalanche-cli/pkg/validatormanager"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
)

func main() {
	chainIDFilePath := filepath.Join("data", "chain.txt")
	chainIDBytes, err := os.ReadFile(chainIDFilePath)
	if err != nil {
		log.Fatalf("❌ Failed to read chain ID file: %s\n", err)
	}
	chainID := ids.FromStringOrPanic(string(chainIDBytes))

	privKey, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("❌ Failed to load key from file: %s\n", err)
	}
	kc := secp256k1fx.NewKeychain(privKey)

	subnetIDBytes, err := os.ReadFile("data/subnet.txt")
	if err != nil {
		log.Fatalf("❌ Failed to read subnet ID file: %s\n", err)
	}
	subnetID := ids.FromStringOrPanic(string(subnetIDBytes))

	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          lib.ETNA_RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
		SubnetIDs:    []ids.ID{subnetID},
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize wallet: %s\n", err)
	}

	softKey, err := key.NewSoft(lib.NETWORK_ID, key.WithPrivateKey(privKey))
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
	for nodeNumber := 0; nodeNumber < lib.VALIDATORS_COUNT; nodeNumber++ {
		configBytes, err := os.ReadFile(filepath.Join("data", "configs", fmt.Sprintf("config-node%d.json", nodeNumber)))
		if err != nil {
			log.Fatalf("❌ Failed to read config file: %s\n", err)
		}
		nodeConfig := lib.NodeConfig{}
		err = json.Unmarshal(configBytes, &nodeConfig)
		if err != nil {
			log.Fatalf("❌ Failed to unmarshal config: %s\n", err)
		}

		endpoint := fmt.Sprintf("http://%s:%s", nodeConfig.PublicIP, nodeConfig.HTTPPort)

		infoClient := info.NewClient(endpoint)
		ctx, cancel := utils.GetAPILargeContext()
		defer cancel()
		nodeID, proofOfPossession, err := infoClient.GetNodeID(ctx)
		if err != nil {
			log.Fatalf("❌ Failed to get node ID: %s\n", err)
		}
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
	}

	avaGoBootstrapValidators, err := blockchaincmd.ConvertToAvalancheGoSubnetValidator(validators)
	if err != nil {
		log.Fatalf("❌ Failed to convert to AvalancheGo subnet validator: %s\n", err)
	}

	managerAddress := goethereumcommon.HexToAddress(validatormanager.ValidatorContractAddress)
	options := getMultisigTxOptions(subnetAuthKeys, kc)
	unsignedTx, err := wallet.P().Builder().NewConvertSubnetTx(
		subnetID,
		chainID,
		managerAddress.Bytes(),
		avaGoBootstrapValidators,
		options...,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create convert subnet tx: %s\n", err)
	}

	tx := txs.Tx{Unsigned: unsignedTx}
	if err := wallet.P().Signer().Sign(context.Background(), &tx); err != nil {
		log.Fatalf("❌ Failed to sign convert subnet tx: %s\n", err)
	}

	_, remainingSubnetAuthKeys, err := txutils.GetRemainingSigners(&tx, []string{changeOwnerAddress})
	if err != nil {
		log.Fatalf("❌ Failed to get remaining subnet auth keys: %s\n", err)
	}
	isFullySigned := len(remainingSubnetAuthKeys) == 0

	id := tx.TxID
	if isFullySigned {
		fmt.Printf("Tx is fully signed with ID: %s\n", id)
		err = wallet.P().IssueTx(&tx)
		if err != nil {
			log.Fatalf("❌ Failed to commit convert subnet tx: %s\n", err)
		}
	} else {
		log.Fatalf("❌ Convert subnet tx is not fully signed")
	}

	fmt.Printf("✅ Convert subnet tx ID: %s\n", id)
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
