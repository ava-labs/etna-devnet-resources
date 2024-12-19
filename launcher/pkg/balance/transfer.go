package balance

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ethereum/go-ethereum/common"
)

func GetCChainBalance(addr common.Address, rpcUrl string) (*big.Int, error) {
	// Connect to C-chain and check balance
	cChainClient, err := ethclient.Dial(rpcUrl + "/ext/bc/C/rpc")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to c-chain: %w", err)
	}

	balance, err := cChainClient.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert balance from C-Chain representation to P-Chain representation
	balance = balance.Div(balance, big.NewInt(int64(1e9)))
	return balance, nil
}

func ImportCToP(PrivateKey *secp256k1.PrivateKey, rpcUrl string) (ids.ID, error) {
	// Get addresses for both chains
	pChainAddr := PrivateKey.Address()
	cChainAddr := evm.PublicKeyToEthAddress(PrivateKey.PublicKey())

	log.Printf("C-chain address: %s", cChainAddr.Hex())

	pBalance, err := CheckPChainBalance(context.Background(), pChainAddr, rpcUrl)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to get new P-chain balance: %w", err)
	}
	log.Printf("P-chain balance: %s AVAX", GetBalanceString(pBalance, 9))

	// Replace the C-chain balance check with the new function
	cChainBalance, err := GetCChainBalance(cChainAddr, rpcUrl)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to get C-chain balance: %w", err)
	}
	log.Printf("C-chain balance: %s AVAX", GetBalanceString(cChainBalance, 9))

	// Check if balance is greater than 1 AVAX
	if cChainBalance.Uint64() < units.Avax {
		fmt.Println("Balance on C-chain is less than 1 AVAX, skipping transfer")
		return ids.Empty, nil
	}

	// Create wallet
	kc := secp256k1fx.NewKeychain(PrivateKey)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          rpcUrl,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Setup owner configuration
	owner := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			pChainAddr,
		},
	}

	// Export from C-chain
	cWallet := wallet.C()
	exportTx, err := cWallet.IssueExportTx(
		constants.PlatformChainID,
		[]*secp256k1fx.TransferOutput{{
			Amt:          cChainBalance.Uint64() - 100*units.MilliAvax, // Leave some for fees
			OutputOwners: owner,
		}},
	)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to export: %w", err)
	}
	log.Printf("Exported from C-chain: %s", exportTx.ID())

	// Import to P-chain
	pWallet := wallet.P()
	importTx, err := pWallet.IssueImportTx(cWallet.Builder().Context().BlockchainID, &owner)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to import: %w", err)
	}
	log.Printf("Imported to P-chain: %s", importTx.ID())

	// Update the final balance check
	newCBalance, err := GetCChainBalance(cChainAddr, rpcUrl)
	if err != nil {
		return importTx.ID(), fmt.Errorf("failed to get new C-chain balance: %w", err)
	}
	log.Printf("New C-chain balance: %s AVAX", GetBalanceString(newCBalance, 9))

	return importTx.ID(), nil
}

func GetBalanceString(balance *big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	quotient := new(big.Int).Div(balance, divisor)
	remainder := new(big.Int).Mod(balance, divisor)
	return fmt.Sprintf("%d.%0*d", quotient, decimals, remainder)
}

func CheckCChainBalance(ctx context.Context, addr common.Address, rpcUrl string) (*big.Int, error) {
	client, err := ethclient.Dial(rpcUrl + "/ext/bc/C/rpc")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to c-chain: %w", err)
	}

	balance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func CheckPChainBalance(ctx context.Context, addr ids.ShortID, rpcUrl string) (*big.Int, error) {
	client := platformvm.NewClient(rpcUrl)

	balance, err := client.GetBalance(ctx, []ids.ShortID{addr})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert json.Uint64 to *big.Int
	return new(big.Int).SetUint64(uint64(balance.Balance)), nil
}
