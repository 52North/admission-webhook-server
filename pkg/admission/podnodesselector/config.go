package podnodesselector

import (
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
)

// Get configuration map
func getConfiguredSelectorMap() (map[string]labels.Set, error) {
	// Don't process if no configuration is set
	env := os.Getenv(ENV_POD_NODES_SELECTOR_CONFIG)
	if len(env) == 0 {
		return nil, nil
	}

	selectors := make(map[string]labels.Set)
	for _, ns := range strings.Split(env, namespaceSeparator) {
		conf := strings.Split(ns, namespaceLabelSeparator)

		// If no namespace name or label not set, move on
		if len(conf) != 2 || len(conf[0]) == 0 || len(conf[1]) == 0 {
			continue
		}

		set, err := labels.ConvertSelectorToLabelsMap(conf[1])
		if err != nil {
			return nil, err
		}

		selectors[conf[0]] = set
	}

	return selectors, nil
}
