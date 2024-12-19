package l1

import (
	"fmt"
	"log"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
)

type CreateL1Params struct {
	PrivateKey     *secp256k1.PrivateKey
	RpcURL         string
	Genesis        string
	ManagerAddress common.Address
	NodeInfos      []info.GetNodeIDReply
	ChainName      string
}

func CreateL1(params CreateL1Params) (ids.ID, ids.ID, ids.ID, error) {
	subnetID, err := CreateSubnet(CreateSubnetParams{
		PrivateKey: params.PrivateKey,
		RpcURL:     params.RpcURL,
	})
	if err != nil {
		return ids.ID{}, ids.ID{}, ids.ID{}, fmt.Errorf("failed to create subnet: %s", err)
	}

	log.Printf("Created subnet: %s", subnetID.String())

	chainID, err := CreateChain(CreateChainParams{
		PrivateKey:  params.PrivateKey,
		SubnetID:    subnetID,
		GenesisData: params.Genesis,
		RpcURL:      params.RpcURL,
		ChainName:   params.ChainName,
	})
	if err != nil {
		return ids.ID{}, ids.ID{}, ids.ID{}, fmt.Errorf("failed to create chain: %s", err)
	}

	log.Printf("Created chain: %s", chainID.String())

	conversionID, err := ConvertToL1(ConvertToL1Params{
		PrivateKey:     params.PrivateKey,
		SubnetID:       subnetID,
		ChainID:        chainID,
		ManagerAddress: params.ManagerAddress,
		NodeInfos:      params.NodeInfos,
		RpcUrl:         params.RpcURL,
	})
	if err != nil {
		return ids.ID{}, ids.ID{}, ids.ID{}, fmt.Errorf("failed to convert to L1: %s", err)
	}

	log.Printf("Converted to L1: %s", conversionID.String())

	return chainID, subnetID, conversionID, nil
}
