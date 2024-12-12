package helpers

import "os"

func GetDesiredContractName() string {
	contractName := os.Getenv("CONTRACT_NAME")
	if contractName == "" {
		return "PoAValidatorManager"
	}
	return contractName
}
