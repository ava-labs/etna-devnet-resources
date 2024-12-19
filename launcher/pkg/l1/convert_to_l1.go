package l1

import (
	"context"
	"fmt"
	"time"

	"encoding/hex"

	"github.com/ava-labs/avalanche-cli/cmd/blockchaincmd"
	"github.com/ava-labs/avalanche-cli/pkg/constants"
	"github.com/ava-labs/avalanche-cli/pkg/key"
	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	avagoconstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ethereum/go-ethereum/common"
)

type SubnetValidator struct {
	NodeID            ids.NodeID
	Weight            uint64
	Balance           uint64
	ProofOfPossession *signer.ProofOfPossession
}

type ConvertToL1Params struct {
	PrivateKey     *secp256k1.PrivateKey
	SubnetID       ids.ID
	ChainID        ids.ID
	ManagerAddress common.Address
	NodeInfos      []info.GetNodeIDReply
	RpcUrl         string
}

func ConvertToL1(params ConvertToL1Params) (ids.ID, error) {
	// Create keychain from private key
	kc := secp256k1fx.NewKeychain(params.PrivateKey)

	// Create context with 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Initialize wallet with subnet ID
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          params.RpcUrl,
		AVAXKeychain: kc,
		EthKeychain:  kc,
		SubnetIDs:    []ids.ID{params.SubnetID},
	})
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to initialize wallet: %s", err)
	}

	softKey, err := key.NewSoft(avagoconstants.TestnetID, key.WithPrivateKey(params.PrivateKey))
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to create change owner address: %w", err)
	}

	changeOwnerAddress := softKey.P()[0]

	// Convert manager address from hex string to bytes
	managerAddressBytes := params.ManagerAddress.Bytes()

	validators := []models.SubnetValidator{}
	for _, nodeInfo := range params.NodeInfos {
		validators = append(validators, models.SubnetValidator{
			NodeID:               nodeInfo.NodeID.String(),
			Weight:               constants.BootstrapValidatorWeight,
			Balance:              constants.BootstrapValidatorBalance,
			BLSPublicKey:         "0x" + hex.EncodeToString(nodeInfo.NodePOP.PublicKey[:]),
			BLSProofOfPossession: "0x" + hex.EncodeToString(nodeInfo.NodePOP.ProofOfPossession[:]),
			ChangeOwnerAddr:      changeOwnerAddress,
		})
	}

	avaGoBootstrapValidators, err := blockchaincmd.ConvertToAvalancheGoSubnetValidator(validators)
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to convert to AvalancheGo subnet validator: %w", err)
	}

	// Issue convert subnet transaction
	tx, err := wallet.P().IssueConvertSubnetToL1Tx(
		params.SubnetID,
		params.ChainID,
		managerAddressBytes,
		avaGoBootstrapValidators,
	)
	if err != nil {
		return ids.ID{}, fmt.Errorf("failed to issue convert subnet transaction: %s", err)
	}

	return tx.ID(), nil
}
