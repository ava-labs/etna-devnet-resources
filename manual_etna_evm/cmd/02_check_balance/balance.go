package main

import (
	"context"
	"encoding/hex"
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
	goethereumcommon "github.com/ethereum/go-ethereum/common"
	goethereumtypes "github.com/ethereum/go-ethereum/core/types"
	goethereumcrypto "github.com/ethereum/go-ethereum/crypto"
	goethereumethclient "github.com/ethereum/go-ethereum/ethclient"
)

const PrefundedEwoqPrivate = "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

func checkPChainBalance(ctx context.Context, addr ids.ShortID) (*big.Int, error) {
	uri := lib.ETNA_RPC_URL
	addresses := set.Of(addr)

	fetchStartTime := time.Now()
	state, err := primary.FetchState(ctx, uri, addresses)
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
		log.Printf("P-chain balance: %s\n", pChainBalance)
		if pChainBalance.Cmp(big.NewInt(int64(lib.MIN_BALANCE))) >= 0 {
			log.Printf("✅ P-chain balance sufficient")
			os.Exit(0)
		}
	}

	log.Printf("P-chain balance insufficient on address %s: %s < %d\n", pChainAddr.String(), pChainBalance, lib.MIN_BALANCE)

	ethAddr := evm.PublicKeyToEthAddress(key.PublicKey())

	cChainClient, err := ethclient.Dial(lib.ETNA_RPC_URL + "/ext/bc/C/rpc")
	if err != nil {
		log.Fatalf("failed to connect to c-chain: %s\n", err)
	}

	cChainBalance, err := cChainClient.BalanceAt(context.Background(), ethAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %s\n", err)
	}
	cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))
	// The P chain balance is in nDEVAX (10-9), but the C-chain balance is in WEI (10-18)
	// So we need to convert it to the same unit

	log.Printf("Balance on c-chain at address %s: %s\n", ethAddr.Hex(), cChainBalance)

	if cChainBalance.Uint64() < lib.MIN_BALANCE {
		log.Printf("Balance %s is less than minimum balance: %d\n", cChainBalance, lib.MIN_BALANCE)
		err := transferFromEwoq(ethAddr.Hex(), lib.MIN_BALANCE*2)
		if err != nil {
			log.Fatalf("failed to transfer from ewoq: %s\n", err)
		}
		cChainBalance, err = cChainClient.BalanceAt(context.Background(), ethAddr, nil)
		if err != nil {
			log.Fatalf("failed to get balance: %s\n", err)
		}
		cChainBalance = cChainBalance.Div(cChainBalance, big.NewInt(int64(1e9)))
		if cChainBalance.Uint64() < lib.MIN_BALANCE {
			log.Printf("❌ Balance %s is less than minimum balance: %d\n", cChainBalance, lib.MIN_BALANCE)
			log.Printf("Please visit " + lib.FAUCET_LINK)
			log.Printf("Use this address to request funds: %s\n", ethAddr.Hex())
			os.Exit(1)
		}
	} else {
		log.Printf("C-chain balance sufficient: current %s, required %d\n", cChainBalance, lib.MIN_BALANCE)
	}

	log.Printf("Transferring %d from C-chain to P-chain\n", lib.MIN_BALANCE)

	// Create keychain and wallet
	kc := secp256k1fx.NewKeychain(key)
	wallet, err := primary.MakeWallet(context.Background(), &primary.WalletConfig{
		URI:          lib.ETNA_RPC_URL,
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

func transferFromEwoq(receiverAddr string, amountNDevax uint64) error {
	// Convert ewoq key to ECDSA private key
	ewoqkeyBytes, err := hex.DecodeString(PrefundedEwoqPrivate)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}
	sourceKey, err := goethereumcrypto.ToECDSA(ewoqkeyBytes)
	if err != nil {
		return fmt.Errorf("failed to convert to private key: %w", err)
	}

	// Convert receiver address string to address type
	destAddr := goethereumcommon.HexToAddress(receiverAddr)

	// Connect to the network
	client, err := goethereumethclient.Dial(lib.ETNA_RPC_URL + "/ext/bc/C/rpc")
	if err != nil {
		return fmt.Errorf("failed to connect to network: %w", err)
	}

	// Get the sender's address and nonce
	fromAddress := goethereumcrypto.PubkeyToAddress(sourceKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Convert amount from nDEVAX to wei (nDEVAX = 10^-9 AVAX, wei = 10^-18 AVAX)
	// So we multiply by 10^9 to get to wei
	value := new(big.Int).Mul(big.NewInt(int64(amountNDevax)), big.NewInt(1000000000))
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create and sign transaction
	tx := goethereumtypes.NewTransaction(nonce, destAddr, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	signedTx, err := goethereumtypes.SignTx(tx, goethereumtypes.NewEIP155Signer(chainID), sourceKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
	return nil
}
