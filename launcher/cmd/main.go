package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/genesis"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files from dist directory
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/generateGenesis", generateGenesis)
	mux.HandleFunc("/api/createL1", createL1)
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

	if evmChainId < 1 || evmChainId > 1000000 {
		http.Error(w, "Invalid evmChainId", http.StatusBadRequest)
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
