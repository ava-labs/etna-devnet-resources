package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"mypkg/lib"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
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

var MIN_BALANCE = units.Avax*lib.VALIDATORS_COUNT + 100*units.MilliAvax

func getBalanceString(balance *big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	quotient := new(big.Int).Div(balance, divisor)
	remainder := new(big.Int).Mod(balance, divisor)
	return fmt.Sprintf("%d.%0*d", quotient, decimals, remainder)
}

func main() {
	key, err := lib.LoadKeyFromFile(lib.VALIDATOR_MANAGER_OWNER_KEY_PATH)
	if err != nil {
		log.Fatalf("failed to load key from file: %s\n", err)
	}

	//has to check balance on P-chain, if less than MIN_BALANCE, then check balance on C-chain
	//If not enough balance on C-chain, then exit with message to visit faucet
	//If enough balance on C-chain, then transfer MIN_BALANCE from C-chain to P-chain
	// if initially enough balance on P-chain, then exit 0

	pChainAddr := key.Address()

	pChainBalance, err := checkPChainBalance(context.Background(), pChainAddr)
	if err != nil {
		log.Printf("Failed to check P-chain balance: %s\n", err)
	} else {
		log.Printf("P-chain balance: %s AVAX\n", getBalanceString(pChainBalance, 9))
		if pChainBalance.Cmp(big.NewInt(int64(lib.MIN_BALANCE))) >= 0 {
			log.Printf("✅ P-chain balance sufficient")
			os.Exit(0)
		}
	}

	log.Printf("P-chain balance insufficient on address %s: %s < %d\n", pChainAddr.String(), pChainBalance, lib.MIN_BALANCE)

	ethAddr := evm.PublicKeyToEthAddress(key.PublicKey())

	cChainClient, err := ethclient.Dial(lib.RPC_URL + "/ext/bc/C/rpc")
	if err != nil {
		log.Fatalf("failed to connect to c-chain: %s\n", err)
	}

	cChainBalance, err := cChainClient.BalanceAt(context.Background(), ethAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %s\n", err)
	}
	// The P chain balance is in nDEVAX (10-9), but the C-chain balance is in WEI (10-18)
	// So we need to convert it to the same unit
	cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))

	log.Printf("Balance on c-chain at address %s: %s\n", ethAddr.Hex(), cChainBalance)

	if cChainBalance.Uint64() < lib.MIN_BALANCE {
		log.Printf("❌ Balance %s is less than minimum balance: %d\n", cChainBalance, lib.MIN_BALANCE)
		log.Printf("Please visit https://core.app/tools/testnet-faucet/?subnet=c&token=c \n")
		log.Printf("Use this address to request funds: %s\n", ethAddr.Hex())
		os.Exit(1)
	} else {
		log.Printf("C-chain balance sufficient: current %s, required %d\n", cChainBalance, lib.MIN_BALANCE)
	}

	log.Printf("Transferring %d from C-chain to P-chain\n", lib.MIN_BALANCE)

	// Create keychain and wallet
	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          lib.RPC_URL,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}

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
			Amt:          lib.MIN_BALANCE + units.MilliAvax*100,
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

	// Check P-chain balance again after import
	pChainBalance, err = checkPChainBalance(context.Background(), pChainAddr)
	if err != nil {
		log.Fatalf("failed to get P-chain balance: %s\n", err)
	}
	if pChainBalance.Cmp(new(big.Int).SetUint64(lib.MIN_BALANCE)) < 0 {
		log.Fatalf("❌ Final P-chain balance %s is less than minimum required %d\n", pChainBalance, lib.MIN_BALANCE)
	}
	log.Printf("✅ Final P-chain balance: %s (greater than minimum %d)\n", pChainBalance, lib.MIN_BALANCE)
}

func checkPChainBalance(ctx context.Context, addr ids.ShortID) (*big.Int, error) {
	addresses := set.Of(addr)

	fetchStartTime := time.Now()
	state, err := primary.FetchState(ctx, lib.RPC_URL, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch state: %w", err)
	}
	log.Printf("fetched state in %s\n", time.Since(fetchStartTime))

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
