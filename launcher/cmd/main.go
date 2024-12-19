package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/balance"
	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/config"
	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/genesis"
	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/l1"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}
}

var lastImportTime time.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
var compileTs int64 = 0

func main() {
	mux := http.NewServeMux()

	privKey := config.LoadOrGeneratePrivateKey()

	// Serve static files from dist directory
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/genesis", logRequest(generateGenesis))
	mux.HandleFunc("/api/create", logRequest(createL1))
	mux.HandleFunc("/api/compiled", logRequest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if compileTs == 0 {
			fi, _ := os.Stat(os.Args[0])
			compileTs = fi.ModTime().Unix()
		}

		w.Write([]byte(strconv.FormatInt(compileTs, 10)))
	}))
	mux.HandleFunc("/api/addr/c", logRequest(func(w http.ResponseWriter, r *http.Request) {
		cChainAddr := evm.PublicKeyToEthAddress(privKey.PublicKey())
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(cChainAddr.Hex()))
	}))
	mux.HandleFunc("/api/balance/c", logRequest(func(w http.ResponseWriter, r *http.Request) {
		cChainAddr := evm.PublicKeyToEthAddress(privKey.PublicKey())
		myBalance, err := balance.CheckCChainBalance(context.Background(), cChainAddr, config.GetRPCUrl())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if time.Since(lastImportTime) > 1*time.Minute && myBalance.Cmp(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)) > 0 {
			lastImportTime = time.Now()
			_, err = balance.ImportCToP(privKey, config.GetRPCUrl())
			if err != nil {
				log.Printf("Failed to import C-chain balance to P-chain: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(balance.GetBalanceString(myBalance, 18)))
	}))
	mux.HandleFunc("/api/balance/p", logRequest(func(w http.ResponseWriter, r *http.Request) {
		myBalance, err := balance.CheckPChainBalance(context.Background(), privKey.Address(), config.GetRPCUrl())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(balance.GetBalanceString(myBalance, 9)))
	}))

	port := "3000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type CreateL1Request struct {
	GenesisString string                `json:"genesisString"`
	Nodes         []info.GetNodeIDReply `json:"nodes"`
	L1Name        string                `json:"l1Name"`
}

func createL1(w http.ResponseWriter, r *http.Request) {
	var req CreateL1Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Nodes) == 0 {
		http.Error(w, "No nodes provided", http.StatusBadRequest)
		return
	}

	if req.GenesisString == "" {
		http.Error(w, "No genesis string provided", http.StatusBadRequest)
		return
	}

	for _, node := range req.Nodes {
		if node.NodePOP == nil {
			http.Error(w, "No node POP provided", http.StatusBadRequest)
			return
		}
		if err := node.NodePOP.Verify(); err != nil {
			http.Error(w, "Invalid node POP", http.StatusBadRequest)
			return
		}
	}

	if len(req.L1Name) < 1 || len(req.L1Name) > 32 {
		http.Error(w, "L1Name name must be between 1 and 32 characters", http.StatusBadRequest)
		return
	}

	if _, hasMock := r.URL.Query()["mock"]; hasMock {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"chainID":      "24dnYhG4YgY2xpmccLJvgrq91DZuTtciU8GiLcHHBzNdSdN7TM",
			"subnetID":     "uBqfpRFjEDnaA4tJhULzrgszt1yaH65Cb4NigfY59LL56wcDH",
			"conversionID": "272Da6ib1eGfH72UYEnhNJHhJA77bNgLMZwT3vGwZSJeao9WEH",
		})
		return
	}

	privKey := config.LoadOrGeneratePrivateKey()

	params := l1.CreateL1Params{
		PrivateKey:     privKey,
		RpcURL:         config.GetRPCUrl(),
		Genesis:        req.GenesisString,
		ManagerAddress: evm.PublicKeyToEthAddress(privKey.PublicKey()),
		NodeInfos:      req.Nodes,
		ChainName:      req.L1Name,
	}

	chainID, subnetID, conversionID, err := l1.CreateL1(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Chain ID: %s", chainID.String())
	log.Printf("Subnet ID: %s", subnetID.String())
	log.Printf("Conversion ID: %s", conversionID.String())

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"chainID":      chainID.String(),
		"subnetID":     subnetID.String(),
		"conversionID": conversionID.String(),
	})
}

func generateGenesis(w http.ResponseWriter, r *http.Request) {
	ownerEthAddressString := r.URL.Query().Get("ownerEthAddressString")
	evmChainIdStr := r.URL.Query().Get("evmChainId")

	evmChainId := 0
	if evmChainIdStr != "" {
		var err error
		evmChainId, err = strconv.Atoi(evmChainIdStr)
		if err != nil {
			http.Error(w, "Invalid evmChainId: must be a valid number", http.StatusBadRequest)
			return
		}
	}

	maxEvmChainId := 1000000
	if evmChainId < 1 || evmChainId > maxEvmChainId {
		http.Error(w, fmt.Sprintf("Invalid evmChainId, should be between 1 and %d", maxEvmChainId), http.StatusBadRequest)
		return
	}

	payload := genesis.GeneratePayload{
		OwnerEthAddressString: ownerEthAddressString,
		EvmChainId:            evmChainId,
	}

	genesis, err := genesis.Generate(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add content type header
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(genesis))
}
