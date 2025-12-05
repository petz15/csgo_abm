package strategy

import (
	"encoding/json"
	"math"
	"os"
)

// TreeNode represents a node in the decision tree
type TreeNode struct {
	IsLeaf    bool      `json:"is_leaf"`
	Feature   string    `json:"feature,omitempty"`
	Threshold float64   `json:"threshold,omitempty"`
	Value     float64   `json:"value,omitempty"`
	Left      *TreeNode `json:"left,omitempty"`
	Right     *TreeNode `json:"right,omitempty"`
}

// TreeModel represents a decision tree model loaded from JSON
type TreeModel struct {
	Tree          TreeNode           `json:"tree"`
	StateFeatures []string           `json:"state_features"`
	Normalization map[string]float64 `json:"normalization"`
}

var treeModelInstance *TreeModel

// LoadTreeModel loads the decision tree model from JSON file
func LoadTreeModel(path string) (*TreeModel, error) {
	if treeModelInstance != nil {
		return treeModelInstance, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var model TreeModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	treeModelInstance = &model
	return treeModelInstance, nil
}

// Traverse the decision tree to get prediction
func (m *TreeModel) predict(features map[string]float64, node *TreeNode) float64 {
	if node.IsLeaf {
		return node.Value
	}

	featureValue := features[node.Feature]
	if featureValue < node.Threshold {
		return m.predict(features, node.Left)
	}
	return m.predict(features, node.Right)
}

// Prepare input features as a map
func (m *TreeModel) prepareInput(ctx StrategyContext_simple) map[string]float64 {
	features := make(map[string]float64)

	// Normalize features according to the model's normalization
	features["own_funds"] = ctx.Funds / m.Normalization["own_funds"]
	features["own_score"] = float64(ctx.OwnScore) / m.Normalization["own_score"]
	features["opponent_score"] = float64(ctx.OpponentScore) / m.Normalization["opponent_score"]
	features["own_survivors"] = float64(ctx.OwnSurvivors) / m.Normalization["own_survivors"]
	features["opponent_survivors"] = float64(ctx.EnemySurvivors) / m.Normalization["opponent_survivors"]
	features["consecutive_losses"] = float64(ctx.ConsecutiveLosses) / m.Normalization["consecutive_losses"]

	// is_ct_side: 1.0 if CT, 0.0 if T
	if ctx.Side {
		features["is_ct_side"] = 1.0
	} else {
		features["is_ct_side"] = 0.0
	}

	features["round_number"] = float64(ctx.CurrentRound) / m.Normalization["round_number"]
	features["half_length"] = float64(ctx.GameRules_strategy.HalfLength) / m.Normalization["half_length"]
	features["last_round_reason"] = float64(ctx.RoundEndReason) / m.Normalization["last_round_reason"]

	// last_bomb_planted: 1.0 if planted, 0.0 otherwise
	if ctx.Is_BombPlanted {
		features["last_bomb_planted"] = 1.0
	} else {
		features["last_bomb_planted"] = 0.0
	}

	return features
}

// InvestDecisionMaking_ml_tree uses the decision tree model for buy decisions
func InvestDecisionMaking_ml_tree(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadTreeModel("tree_model.json")
	if err != nil {
		// Fallback to default strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input
	features := model.prepareInput(ctx)

	// Get prediction (0.0 to 1.0 range representing fraction of funds to spend)
	prediction := model.predict(features, &model.Tree)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount
	return ctx.Funds * prediction
}
