package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mypkg/lib"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/ava-labs/avalanche-cli/pkg/constants"
)

func main() {
	err := os.MkdirAll(filepath.Join("data", "configs"), 0755)
	if err != nil {
		log.Fatalf("❌ Failed to create configs directory: %s\n", err)
	}

	openPorts, err := lib.FindMultipleFreePorts(lib.VALIDATORS_COUNT*2, 9650)
	if err != nil {
		log.Fatalf("❌ Failed to find free ports: %s\n", err)
	}

	for i := 0; i < lib.VALIDATORS_COUNT; i++ {

		config := lib.NodeConfig{
			APIAdminEnabled:          "true",
			BootstrapIDs:             strings.Join(constants.EtnaDevnetBootstrapNodeIDs, ","),
			BootstrapIPs:             strings.Join(constants.EtnaDevnetBootstrapIPs, ","),
			DataDir:                  fmt.Sprintf("/data/node%d", i),
			DbDir:                    fmt.Sprintf("/data/node%d/db", i),
			GenesisFile:              "/data/genesis.json",
			HealthCheckFrequency:     "2s",
			HTTPPort:                 fmt.Sprintf("%d", openPorts[i*2]),
			IndexEnabled:             "true",
			LogDir:                   fmt.Sprintf("/data/node%d/logs", i),
			LogDisplayLevel:          "INFO",
			LogLevel:                 "INFO",
			NetworkID:                "76",
			NetworkMaxReconnectDelay: "1s",
			PluginDir:                "/data/plugins/",
			PublicIP:                 "127.0.0.1",
			StakingPort:              fmt.Sprintf("%d", openPorts[i*2+1]),
			TrackSubnets:             "",
			UpgradeFile:              "/data/upgrade.json",
		}

		marshalled, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatalf("❌ Failed to marshal config: %s\n", err)
		}
		err = os.WriteFile(filepath.Join("data", "configs", fmt.Sprintf("config-node%d.json", i)), marshalled, 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write config: %s\n", err)
		}

		err = createMultipleFolders([]string{
			config.DbDir[1:],
			config.LogDir[1:],
		})
		if err != nil {
			log.Fatalf("❌ Failed to create folders: %s\n", err)
		}
	}

	err = os.WriteFile("data/upgrade.json", constants.EtnaDevnetUpgradeData, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write upgrade file: %s\n", err)
	}

	err = os.WriteFile("data/genesis.json", constants.EtnaDevnetGenesisData, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write genesis file: %s\n", err)
	}

	//FIXME: needs plugins
	err = os.MkdirAll(filepath.Join("data", "plugins"), 0755)
	if err != nil {
		log.Fatalf("❌ Failed to create plugins directory: %s\n", err)
	}

	fmt.Println("✅ Successfully created configs")
}

func createMultipleFolders(folders []string) error {
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
