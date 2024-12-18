package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ava-labs/etna-devnet-resources/launcher/pkg/genesis"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files from dist directory
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/generateGenesis", func(w http.ResponseWriter, r *http.Request) {
		ownerEthAddressString := r.URL.Query().Get("ownerEthAddressString")
		evmChainIdStr := r.URL.Query().Get("evmChainId")

		evmChainId := 0
		if evmChainIdStr != "" {
			var err error
			evmChainId, err = strconv.Atoi(evmChainIdStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
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
		w.Write([]byte(genesis))
	})

	port := "3000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
