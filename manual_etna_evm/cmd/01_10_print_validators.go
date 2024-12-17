package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printPChainInfoCmd)
}

var printPChainInfoCmd = &cobra.Command{
	Use:   "validators",
	Short: "Print validators",
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🧱 Printing validators")

		if err := printPChainState(); err != nil {
			return fmt.Errorf("failed to print P-Chain state: %w", err)
		}

		return nil
	},
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
	ID      int         `json:"id"`
}

func makeJSONRPCRequest(client *http.Client, url string, payload map[string]interface{}) (*JSONRPCResponse, error) {
	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	var jsonResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for JSON-RPC error
	if jsonResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %v", jsonResp.Error)
	}

	return &jsonResp, nil
}

type ValidatorInfo struct {
	PublicKey string `json:"publicKey"`
	Weight    string `json:"weight"`
}

type ValidatorsResponse struct {
	Validators map[string]ValidatorInfo
}

func callPChainValidatorsAt(fujiPChainURL string, subnetID string) (*ValidatorsResponse, error) {
	client := &http.Client{}
	validatorsPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "platform.getValidatorsAt",
		"params": map[string]string{
			"height":   "proposed",
			"subnetID": subnetID,
		},
		"id": 1,
	}

	resp, err := makeJSONRPCRequest(client, fujiPChainURL, validatorsPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Parse the result into our struct
	validatorsResp := &ValidatorsResponse{
		Validators: make(map[string]ValidatorInfo),
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	for nodeID, details := range resultMap {
		validatorDetails, ok := details.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid validator details format for node %s", nodeID)
		}

		validatorsResp.Validators[nodeID] = ValidatorInfo{
			PublicKey: validatorDetails["publicKey"].(string),
			Weight:    validatorDetails["weight"].(string),
		}
	}

	return validatorsResp, nil
}

func printPChainState() error {
	subnetID, err := helpers.LoadId(helpers.SubnetIdPath)
	if err != nil {
		return fmt.Errorf("failed to load subnet ID: %w", err)
	}

	// Make HTTP requests
	client := &http.Client{}
	fujiPChainURL := "https://api.avax-test.network/ext/P"

	// Get validators
	validatorsResp, err := callPChainValidatorsAt(fujiPChainURL, subnetID.String())
	if err != nil {
		return fmt.Errorf("failed to get validators: %w", err)
	}

	// Get subnet info
	subnetPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "platform.getSubnet",
		"params": map[string]string{
			"subnetID": subnetID.String(),
		},
		"id": 1,
	}

	subnetResp, err := makeJSONRPCRequest(client, fujiPChainURL, subnetPayload)
	if err != nil {
		return fmt.Errorf("failed to get subnet info: %w", err)
	}

	// Print results
	fmt.Println("P-Chain State:")
	fmt.Println("------------------------")

	fmt.Println("\nValidators:")
	for nodeID, details := range validatorsResp.Validators {
		fmt.Printf("NodeID: %s\n", nodeID)
		fmt.Printf("  Public Key: %s\n", details.PublicKey)
		fmt.Printf("  Weight: %s\n", details.Weight)
	}

	fmt.Println("\nSubnet Info:")
	subnetInfo := subnetResp.Result.(map[string]interface{})
	fmt.Printf("Is Permissioned: %v\n", subnetInfo["isPermissioned"])
	fmt.Printf("Control Keys: %v\n", subnetInfo["controlKeys"])
	fmt.Printf("Threshold: %s\n", subnetInfo["threshold"])
	fmt.Printf("Manager Chain ID: %s\n", subnetInfo["managerChainID"])
	fmt.Printf("Manager Address: %s\n", subnetInfo["managerAddress"])
	fmt.Printf("IsPermissioned: %v\n", subnetInfo["isPermissioned"])

	return nil
}
