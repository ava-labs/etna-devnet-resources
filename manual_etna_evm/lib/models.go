package lib

type NodeConfig struct {
	APIAdminEnabled          string `json:"api-admin-enabled"`
	BootstrapIDs             string `json:"bootstrap-ids"`
	BootstrapIPs             string `json:"bootstrap-ips"`
	DataDir                  string `json:"data-dir"`
	DbDir                    string `json:"db-dir"`
	GenesisFile              string `json:"genesis-file"`
	HealthCheckFrequency     string `json:"health-check-frequency"`
	HTTPPort                 string `json:"http-port"`
	IndexEnabled             string `json:"index-enabled"`
	LogDir                   string `json:"log-dir"`
	LogDisplayLevel          string `json:"log-display-level"`
	LogLevel                 string `json:"log-level"`
	NetworkID                string `json:"network-id"`
	NetworkMaxReconnectDelay string `json:"network-max-reconnect-delay"`
	PluginDir                string `json:"plugin-dir"`
	PublicIP                 string `json:"public-ip"`
	StakingPort              string `json:"staking-port"`
	TrackSubnets             string `json:"track-subnets,omitempty"`
	UpgradeFile              string `json:"upgrade-file"`
}
