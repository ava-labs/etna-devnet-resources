package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"mypkg/lib"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/coreth/ethclient"
)

func main() {
	keyHex := os.Args[1]
	wallet, kc, err := walletFromHex(keyHex)
	if err != nil {
		log.Fatalf("failed to create key chain: %s\n", err)
	}

	ethAddr := kc.EthAddrs.List()[0]
	pChainAddr := kc.Addrs.List()[0]

	cChainClient, err := ethclient.Dial(lib.ETNA_RPC_URL + "/ext/bc/C/rpc")
	if err != nil {
		log.Fatalf("failed to connect to c-chain: %s\n", err)
	}

	cChainBalance, err := cChainClient.BalanceAt(context.Background(), ethAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %s\n", err)
	}
	cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))

	toImport := new(big.Int).Sub(cChainBalance, big.NewInt(int64(100*units.MilliAvax)))
	toImportUint64 := toImport.Uint64()

	// Get P-chain and C-chain wallets
	pWallet := wallet.P()
	cWallet := wallet.C()

	// Setup owner configuration
	owner := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			pChainAddr,
		},
	}

	fmt.Println("constants.PlatformChainID", constants.PlatformChainID)

	// Export from C-chain
	exportTx, err := cWallet.IssueExportTx(
		constants.PlatformChainID,
		[]*secp256k1fx.TransferOutput{{
			Amt:          toImportUint64,
			OutputOwners: owner,
		}},
	)
	if err != nil {
		log.Fatalf("failed to issue export transaction: %s\n", err)
	}
	log.Printf("✅ Issued export %s\n", exportTx.ID())

	// Import to P-chain
	importTx, err := pWallet.IssueImportTx(cWallet.Builder().Context().BlockchainID, &owner)
	if err != nil {
		log.Fatalf("failed to issue import transaction: %s\n", err)
	}
	log.Printf("✅ Issued import %s\n", importTx.ID())
}

func walletFromHex(keyHex string) (primary.Wallet, *secp256k1fx.Keychain, error) {
	if strings.HasPrefix(keyHex, "0x") {
		keyHex = keyHex[2:]
	}

	keyBytes, err := hex.DecodeString(strings.TrimSpace(keyHex))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode key hex: %s\n", err)
	}

	key, err := secp256k1.ToPrivateKey(keyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode key: %s\n", err)
	}

	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          lib.ETNA_RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize wallet: %s\n", err)
	}

	return wallet, kc, nil
}
