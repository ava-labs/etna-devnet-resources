package lib

import "github.com/ava-labs/avalanchego/utils/units"

const (
	RPC_URL     = "https://api.avax-test.network"
	MIN_BALANCE = units.Avax*VALIDATORS_COUNT + 100*units.MilliAvax

	VALIDATOR_MANAGER_OWNER_KEY_PATH = "data/poa_validator_manager_owner_key.txt"

	VALIDATORS_COUNT = 1
	NETWORK_ID       = 76
)
