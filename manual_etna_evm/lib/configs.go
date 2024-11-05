package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ava-labs/avalanche-cli/pkg/constants"
)

func FillNodeConfigs(trackSubnets string) error {
	err := os.MkdirAll(filepath.Join("data", "configs"), 0755)
	if err != nil {
		return fmt.Errorf("failed to create configs directory: %w", err)
	}

	openPorts, err := FindMultipleFreePorts(VALIDATORS_COUNT*2, 9650)
	if err != nil {
		return fmt.Errorf("failed to find free ports: %w", err)
	}

	for i := 0; i < VALIDATORS_COUNT; i++ {
		config := NodeConfig{
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
			TrackSubnets:             trackSubnets,
			UpgradeFile:              "/data/upgrade.json",
		}

		marshalled, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		err = os.WriteFile(filepath.Join("data", "configs", fmt.Sprintf("config-node%d.json", i)), marshalled, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		err = createMultipleFolders([]string{
			config.DbDir[1:],
			config.LogDir[1:],
		})
		if err != nil {
			return fmt.Errorf("failed to create folders: %w", err)
		}
	}

	err = os.WriteFile("data/upgrade.json", constants.EtnaDevnetUpgradeData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write upgrade file: %w", err)
	}

	err = os.WriteFile("data/genesis.json", constants.EtnaDevnetGenesisData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write genesis file: %w", err)
	}

	//FIXME: needs plugins
	err = os.MkdirAll(filepath.Join("data", "plugins"), 0755)
	if err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	return nil
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
