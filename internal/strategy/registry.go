package strategy

import (
	"fmt"
	"sort"
	"strings"
)

// StrategyFunc is the function signature for all strategies
type StrategyFunc func(StrategyContext_simple) float64

// StrategyRegistry maps strategy names to their functions
var StrategyRegistry = map[string]StrategyFunc{
	// Aggressive
	"all_in":    InvestDecisionMaking_allin,
	"all_in_v2": InvestDecisionMaking_allin_v2,

	// Counter
	"anti_allin":    InvestDecisionMaking_anti_allin,
	"anti_allin_v2": InvestDecisionMaking_anti_allin_v2,
	"anti_allin_v3": InvestDecisionMaking_anti_allin_v3,

	// Optimized
	"min_max":        InvestDecisionMaking_min_max,
	"min_max_v2":     InvestDecisionMaking_min_max_v2,
	"min_max_v3":     InvestDecisionMaking_min_max_v3,
	"expected_value": InvestDecisionMaking_expected_value,

	// Adaptive
	"adaptive_eco_v1": InvestDecisionMaking_adaptive_v1,
	"adaptive_eco_v2": InvestDecisionMaking_adaptive_v2,
	"smart_v1":        InvestDecisionMaking_smart_v1,
	"half":            InvestDecisionMaking_half,

	// ML-based
	"ml_dqn":    InvestDecisionMaking_ml_dqn,
	"ml_sgd":    InvestDecisionMaking_ml_sgd,
	"ml_tree":   InvestDecisionMaking_ml_tree,
	"ml_forest": InvestDecisionMaking_ml_forest,

	// Basic
	"casual":  InvestDecisionMaking_casual,
	"random":  InvestDecisionMaking_random,
	"scrooge": InvestDecisionMaking_scrooge,
}

// ValidateStrategy checks if a strategy exists
func ValidateStrategy(name string) error {
	if _, exists := StrategyRegistry[name]; !exists {
		available := GetAvailableStrategies()
		return fmt.Errorf("unknown strategy '%s'. Available strategies:\n  %s",
			name, strings.Join(available, ", "))
	}
	return nil
}

// GetStrategy returns the strategy function, or error if not found
func GetStrategy(name string) (StrategyFunc, error) {
	fn, exists := StrategyRegistry[name]
	if !exists {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return fn, nil
}

// GetAvailableStrategies returns sorted list of all strategies
func GetAvailableStrategies() []string {
	strategies := make([]string, 0, len(StrategyRegistry))
	for name := range StrategyRegistry {
		strategies = append(strategies, name)
	}
	sort.Strings(strategies)
	return strategies
}
