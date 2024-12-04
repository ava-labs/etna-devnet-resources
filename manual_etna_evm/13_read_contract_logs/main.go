package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mypkg/config"
	"mypkg/helpers"
	"net/http"

	"github.com/ava-labs/subnet-evm/interfaces"
	poavalidatormanager "github.com/ava-labs/teleporter/abi-bindings/go/validator-manager/PoAValidatorManager"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

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

func main() {
	if err := printPChainState(); err != nil {
		log.Fatalf("failed to print P-Chain state: %s\n", err)
	}
	fmt.Print("\n\n\n\n")
	if err := printEVMContractLogs(); err != nil {
		log.Fatalf("failed to print EVM contract logs: %w", err)
	}
}

func printPChainState() error {
	// Create JSON-RPC request payloads
	validatorsPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "platform.getValidatorsAt",
		"params": map[string]string{
			"height":   "proposed",
			"subnetID": "mNGzCz4iiZgRDFBWcdGyDHpETVR24inQwcyDHGZUQPtmLdynk",
		},
		"id": 1,
	}

	subnetPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "platform.getSubnet",
		"params": map[string]string{
			"subnetID": "mNGzCz4iiZgRDFBWcdGyDHpETVR24inQwcyDHGZUQPtmLdynk",
		},
		"id": 1,
	}

	// Make HTTP requests
	client := &http.Client{}
	fujiPChainURL := "https://rpc.ankr.com/avalanche_fuji-p"

	// Get validators
	validatorsResp, err := makeJSONRPCRequest(client, fujiPChainURL, validatorsPayload)
	if err != nil {
		return fmt.Errorf("failed to get validators: %w", err)
	}

	// Get subnet info
	subnetResp, err := makeJSONRPCRequest(client, fujiPChainURL, subnetPayload)
	if err != nil {
		return fmt.Errorf("failed to get subnet info: %w", err)
	}

	// Print results
	fmt.Println("P-Chain State:")
	fmt.Println("------------------------")

	fmt.Println("\nValidators:")
	for nodeID, details := range validatorsResp.Result.(map[string]interface{}) {
		validatorDetails := details.(map[string]interface{})
		fmt.Printf("NodeID: %s\n", nodeID)
		fmt.Printf("  Public Key: %s\n", validatorDetails["publicKey"])
		fmt.Printf("  Weight: %s\n", validatorDetails["weight"])
	}

	fmt.Println("\nSubnet Info:")
	subnetInfo := subnetResp.Result.(map[string]interface{})
	fmt.Printf("Is Permissioned: %v\n", subnetInfo["isPermissioned"])
	fmt.Printf("Control Keys: %v\n", subnetInfo["controlKeys"])
	fmt.Printf("Threshold: %s\n", subnetInfo["threshold"])
	fmt.Printf("Manager Chain ID: %s\n", subnetInfo["managerChainID"])
	fmt.Printf("Manager Address: %s\n", subnetInfo["managerAddress"])

	return nil
}

func printEVMContractLogs() error {
	managerAddress := common.HexToAddress(config.ProxyContractAddress)

	ethClient, _, err := helpers.GetLocalEthClient()
	if err != nil {
		return fmt.Errorf("failed to connect to client: %s\n", err)
	}

	contract, err := poavalidatormanager.NewPoAValidatorManager(managerAddress, ethClient)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %s\n", err)
	}

	// Get all logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{managerAddress},
	}

	logs, err := ethClient.FilterLogs(context.Background(), (interfaces.FilterQuery)(query))
	if err != nil {
		log.Fatal(err)
	}

	// Print all logs
	for _, vLog := range logs {
		fmt.Println("------------------------")

		fmt.Printf("Log TxHash: %s\n", vLog.TxHash.Hex())

		if event, err := contract.PoAValidatorManagerFilterer.ParseInitialValidatorCreated(vLog); err == nil {
			fmt.Printf("InitialValidatorCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodCreated(vLog); err == nil {
			fmt.Printf("ValidationPeriodCreated:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  NodeID: %x\n", event.NodeID)
			fmt.Printf("  RegisterValidationMessageID: %x\n", event.RegisterValidationMessageID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  RegistrationExpiry: %d\n", event.RegistrationExpiry)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodEnded(vLog); err == nil {
			fmt.Printf("ValidationPeriodEnded:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Status: %d\n", event.Status)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidationPeriodRegistered(vLog); err == nil {
			fmt.Printf("ValidationPeriodRegistered:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  Timestamp: %s\n", event.Timestamp.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorRemovalInitialized(vLog); err == nil {
			fmt.Printf("ValidatorRemovalInitialized:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			fmt.Printf("  Weight: %s\n", event.Weight.String())
			fmt.Printf("  EndTime: %s\n", event.EndTime.String())
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseValidatorWeightUpdate(vLog); err == nil {
			fmt.Printf("ValidatorWeightUpdate:\n")
			fmt.Printf("  ValidationID: %x\n", event.ValidationID)
			fmt.Printf("  Nonce: %d\n", event.Nonce)
			fmt.Printf("  ValidatorWeight: %d\n", event.ValidatorWeight)
			fmt.Printf("  SetWeightMessageID: %x\n", event.SetWeightMessageID)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseInitialized(vLog); err == nil {
			fmt.Printf("Initialized:\n")
			fmt.Printf("  Version: %d\n", event.Version)
			continue
		}

		if event, err := contract.PoAValidatorManagerFilterer.ParseOwnershipTransferred(vLog); err == nil {
			fmt.Printf("OwnershipTransferred:\n")
			fmt.Printf("  Previous Owner: %s\n", event.PreviousOwner.Hex())
			fmt.Printf("  New Owner: %s\n", event.NewOwner.Hex())
			continue
		}

		log.Printf("Failed to parse log: unknown event type\n")
	}
	return nil
}
