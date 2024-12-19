package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/spf13/cobra"

	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/chain/p/builder"
	"github.com/ava-labs/avalanchego/wallet/chain/p/wallet"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ava-labs/coreth/plugin/evm"
)

var TransferCoinsCmd = &cobra.Command{
	Use:   "transfer-coins",
	Short: "Transfer coins between C and P chains",
	Long:  `Transfer coins between C and P chains`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("💰 Transferring AVAX between C and P chains")

		key, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			log.Fatalf("failed to load validator manager owner key: %s\n", err)
		}

		pChainAddr := key.Address()
		cChainAddr := evm.PublicKeyToEthAddress(key.PublicKey())

		pChainBalance, err := CheckPChainBalance(context.Background(), pChainAddr)
		if err != nil {
			log.Printf("Failed to check P-chain balance: %s\n", err)
		} else {
			log.Printf("P-chain balance: %s AVAX\n", GetBalanceString(pChainBalance, 9))
			if pChainBalance.Cmp(big.NewInt(int64(MIN_BALANCE))) >= 0 {
				log.Printf("P-chain balance sufficient")
				return nil
			}
		}

		log.Printf("P-chain balance insufficient on address %s: %s < %s\n", pChainAddr.String(), GetBalanceString(pChainBalance, 9), MIN_BALANCE_STRING)

		cChainClient, err := ethclient.Dial(config.RPC_URL + "/ext/bc/C/rpc")
		if err != nil {
			log.Fatalf("failed to connect to c-chain: %s\n", err)
		}

		cChainBalance, err := cChainClient.BalanceAt(context.Background(), cChainAddr, nil)
		if err != nil {
			log.Fatalf("failed to get balance: %s\n", err)
		}
		// The P chain balance is in nDEVAX (10-9), but the C-chain balance is in WEI (10-18)
		// So we need to convert it to the same unit
		cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))

		log.Printf("Balance on c-chain at address %s: %s\n", cChainAddr.Hex(), GetBalanceString(cChainBalance, 9))

		if cChainBalance.Uint64() < MIN_BALANCE {
			log.Printf("Balance %s is less than minimum balance: %s\n", GetBalanceString(cChainBalance, 9), MIN_BALANCE_STRING)
			log.Printf("Please visit https://test.core.app/tools/testnet-faucet/?subnet=c&token=c \n")
			log.Printf("Use this address to request funds: %s\n", cChainAddr.Hex())
			return fmt.Errorf("transfer to your Fuji C-chain address %s balance to at least %s AVAX", cChainAddr.Hex(), MIN_BALANCE_STRING)
		} else {
			log.Printf("C-chain balance sufficient: current %s, required %s\n", GetBalanceString(cChainBalance, 9), MIN_BALANCE_STRING)
		}

		log.Printf("Transferring balance from C-chain to P-chain\n")

		// Create keychain and wallet
		kc := secp256k1fx.NewKeychain(key)
		wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
			URI:          config.RPC_URL,
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

		log.Println("constants.PlatformChainID", constants.PlatformChainID)

		// Export from C-chain
		exportTx, err := cWallet.IssueExportTx(
			constants.PlatformChainID,
			[]*secp256k1fx.TransferOutput{{
				Amt:          cChainBalance.Uint64() - 100*units.MilliAvax,
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
		pChainBalance, err = CheckPChainBalance(context.Background(), pChainAddr)
		if err != nil {
			log.Fatalf("failed to get P-chain balance: %s\n", err)
		}
		if pChainBalance.Cmp(big.NewInt(int64(MIN_BALANCE))) < 0 {
			log.Fatalf("❌ Final P-chain balance %s is less than minimum required %s\n", GetBalanceString(pChainBalance, 9), MIN_BALANCE_STRING)
		}
		log.Printf("✅ Final P-chain balance: %s (greater than minimum %s)\n", GetBalanceString(pChainBalance, 9), MIN_BALANCE_STRING)

		return nil
	},
}

var MIN_BALANCE = units.Avax + 100*units.MilliAvax
var MIN_BALANCE_STRING = GetBalanceString(big.NewInt(int64(MIN_BALANCE)), 9)

func GetBalanceString(balance *big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	quotient := new(big.Int).Div(balance, divisor)
	remainder := new(big.Int).Mod(balance, divisor)
	return fmt.Sprintf("%d.%0*d", quotient, decimals, remainder)
}

func CheckPChainBalance(ctx context.Context, addr ids.ShortID) (*big.Int, error) {
	addresses := set.Of(addr)

	fetchStartTime := time.Now()
	state, err := primary.FetchState(ctx, config.RPC_URL, addresses)
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

func init() {
	rootCmd.AddCommand(TransferCoinsCmd)
}
