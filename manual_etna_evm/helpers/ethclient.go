package helpers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func DeriveContractAddress(from common.Address, nonce uint64) common.Address {
	encoded, err := rlp.EncodeToBytes([]interface{}{from, nonce})
	if err != nil {
		panic(err)
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToAddress(hash[12:])
}
