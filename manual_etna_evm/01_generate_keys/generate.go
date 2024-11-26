package main

import (
	"log"
	"mypkg/pkg/datafiles"
	"os"
)

func main() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data directory: %s\n", err)
	}

	exists, err := datafiles.ValidatorManagerKeyExists()
	if err != nil {
		log.Fatalf("failed to check if POA validator manager key exists: %s\n", err)
	}

	if exists {
		log.Println("POA validator manager key already exists in ./data/ folder")
	} else {
		if err := datafiles.GenerateValidatorManagerKeyAndSave(); err != nil {
			log.Fatalf("failed to generate POA validator manager key: %s\n", err)
		}
		log.Println("✅ POA validator manager keys generated and saved in ./data/ folder")
	}

}
