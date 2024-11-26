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

	openPorts := []int{9650, 9651} //FindMultipleFreePorts(VALIDATORS_COUNT*2, 9650)
	if len(openPorts) != VALIDATORS_COUNT*2 {
		return fmt.Errorf("failed to find free ports: %w", err)
	}

	for i := 0; i < VALIDATORS_COUNT; i++ {
		config := NodeConfig{
			APIAdminEnabled:          "true",
			BootstrapIDs:             strings.Join(constants.EtnaDevnetBootstrapNodeIDs, ","),
			BootstrapIPs:             strings.Join(constants.EtnaDevnetBootstrapIPs, ","),
			DataDir:                  fmt.Sprintf("/data/node%d", i),
			DbDir:                    fmt.Sprintf("/data/node%d/db", i),
			GenesisFile:              "/data/genesis_fuji.json",
			HealthCheckFrequency:     "2s",
			HTTPPort:                 fmt.Sprintf("%d", openPorts[i*2]),
			IndexEnabled:             "true",
			LogDir:                   fmt.Sprintf("/data/node%d/logs", i),
			LogDisplayLevel:          "INFO",
			LogLevel:                 "INFO",
			NetworkID:                fmt.Sprintf("%d", NETWORK_ID),
			NetworkMaxReconnectDelay: "1s",
			PluginDir:                "/plugins/",
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
