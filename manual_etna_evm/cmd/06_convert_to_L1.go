package cmd

import (
	"context"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/key"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	goethereumcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ConvertToL1Cmd)
}

var ConvertToL1Cmd = &cobra.Command{
	Use:   "convert-to-L1",
	Short: "Convert the subnet to L1",
	Long:  `Convert the subnet to L1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🔌 Converting subnet to L1")

		exists, err := helpers.FileExists(helpers.ConversionIdPath)
		if err != nil {
			return fmt.Errorf("failed to check if conversion ID file exists: %w", err)
		}

		if exists {
			log.Println("✅ Subnet was already converted to L1")
			return nil
		}

		chainID, err := helpers.LoadId(helpers.ChainIdPath)
		if err != nil {
			return fmt.Errorf("failed to load chain ID: %w", err)
		}

		privKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load validator manager owner key: %w", err)
		}
		kc := secp256k1fx.NewKeychain(privKey)

		subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
		if err != nil {
			return fmt.Errorf("failed to load subnet ID: %w", err)
		}

		wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
			URI:          config.RPC_URL,
			AVAXKeychain: kc,
			EthKeychain:  kc,
			SubnetIDs:    []ids.ID{subnetID},
		})
		if err != nil {
			return fmt.Errorf("❌ Failed to initialize wallet: %w", err)
		}

		//TODO: replace with address.Format
		softKey, err := key.NewSoft(avagoconstants.TestnetID, key.WithPrivateKey(privKey))
		if err != nil {
			return fmt.Errorf("❌ Failed to create change owner address: %w", err)
		}

		changeOwnerAddress := softKey.P()[0]

		fmt.Printf("Using changeOwnerAddress: %s\n", changeOwnerAddress)

		subnetAuthKeys, err := address.ParseToIDs([]string{changeOwnerAddress})
		if err != nil {
			return fmt.Errorf("❌ Failed to parse subnet auth keys: %w", err)
		}

		validators := []models.SubnetValidator{}

		nodeID, proofOfPossession, err := NodeInfoFromCreds(helpers.Node0KeysFolder)
		if err != nil {
			return fmt.Errorf("failed to get node info from creds: %w", err)
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

		avaGoBootstrapValidators, err := blockchaincmd.ConvertToAvalancheGoSubnetValidator(validators)
		if err != nil {
			return fmt.Errorf("❌ Failed to convert to AvalancheGo subnet validator: %w", err)
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
			return fmt.Errorf("❌ Failed to write convert log: %w", err)
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
			return fmt.Errorf("❌ Failed to create convert subnet tx: %w", err)
		}

		err = helpers.SaveId(helpers.ConversionIdPath, tx.ID())
		if err != nil {
			return fmt.Errorf("failed to save conversion ID: %w", err)
		}

		log.Printf("✅ Convert subnet tx ID: %s\n", tx.ID().String())
		return nil
	},
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

func NodeInfoFromCreds(folder string) (ids.NodeID, *signer.ProofOfPossession, error) {
	blsKey, err := helpers.LoadBLSKey(folder + "signer.key")
	if err != nil {
		return ids.NodeID{}, nil, fmt.Errorf("failed to load BLS key: %w", err)
	}

	pop := signer.NewProofOfPossession(blsKey)
	certString, err := helpers.LoadText(folder + "staker.crt")
	if err != nil {
		return ids.NodeID{}, nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	block, _ := pem.Decode([]byte(certString))
	if block == nil || block.Type != "CERTIFICATE" {
		panic("failed to decode PEM block containing certificate")
	}

	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to decode PEM block containing certificate")
	}

	nodeID := ids.NodeIDFromCert(cert)

	return nodeID, pop, nil
}
