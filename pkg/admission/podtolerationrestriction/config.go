package podtolerationrestriction

import (
	"encoding/json"
	"os"

	coreV1 "k8s.io/api/core/v1"
)

// Get configuration map
func getConfiguredTolerationsMap() (map[string][]coreV1.Toleration, error) {
	// Don't process if no configuration is set
	s := os.Getenv(ENV_POD_TOLERATION_RESTRICTION_CONFIG)
	tolerations := map[string][]coreV1.Toleration{}
	if s != "" {
		if err := json.Unmarshal([]byte(s), &tolerations); err != nil {
			return nil, err
		}
	}
	return tolerations, nil

}
