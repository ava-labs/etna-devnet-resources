package balance

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/chain/p/builder"
	"github.com/ava-labs/avalanchego/wallet/chain/p/wallet"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ava-labs/coreth/plugin/evm"
)

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

	// Connect to C-chain and check balance
	cChainClient, err := ethclient.Dial(rpcUrl + "/ext/bc/C/rpc")
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to connect to c-chain: %w", err)
	}

	cChainBalance, err := cChainClient.BalanceAt(context.Background(), cChainAddr, nil)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert balance from C-Chain representation to P-Chain representation
	cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))
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

	// Add balance checks after import
	newCBalance, err := cChainClient.BalanceAt(context.Background(), cChainAddr, nil)
	if err != nil {
		return importTx.ID(), fmt.Errorf("failed to get new C-chain balance: %w", err)
	}
	newCBalance = newCBalance.Div(newCBalance, big.NewInt(int64(1e9)))
	log.Printf("New C-chain balance: %s AVAX", GetBalanceString(newCBalance, 9))

	return importTx.ID(), nil
}

func GetBalanceString(balance *big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	quotient := new(big.Int).Div(balance, divisor)
	remainder := new(big.Int).Mod(balance, divisor)
	return fmt.Sprintf("%d.%0*d", quotient, decimals, remainder)
}

func CheckPChainBalance(ctx context.Context, addr ids.ShortID, rpcUrl string) (*big.Int, error) {
	addresses := set.Of(addr)

	state, err := primary.FetchState(ctx, rpcUrl, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch state: %w", err)
	}

	pUTXOs := common.NewChainUTXOs(constants.PlatformChainID, state.UTXOs)
	pBackend := wallet.NewBackend(state.PCTX, pUTXOs, nil)
	pBuilder := builder.New(addresses, state.PCTX, pBackend)

	currentBalances, err := pBuilder.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	avaxID := state.PCTX.AVAXAssetID
	avaxBalance := currentBalances[avaxID]
	return big.NewInt(int64(avaxBalance)), nil
}
