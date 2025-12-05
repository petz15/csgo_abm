package strategy

import (
	"encoding/json"
	"math"
	"os"
)

// ForestModel represents a random forest model loaded from JSON
type ForestModel struct {
	NTrees        int                `json:"n_trees"`
	Trees         []TreeNode         `json:"trees"`
	StateFeatures []string           `json:"state_features"`
	Normalization map[string]float64 `json:"normalization"`
}

var forestModelInstance *ForestModel

// LoadForestModel loads the random forest model from JSON file
func LoadForestModel(path string) (*ForestModel, error) {
	if forestModelInstance != nil {
		return forestModelInstance, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var model ForestModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	forestModelInstance = &model
	return forestModelInstance, nil
}

// Traverse a single tree to get prediction
func (m *ForestModel) predictTree(features map[string]float64, node *TreeNode) float64 {
	if node.IsLeaf {
		return node.Value
	}

	featureValue := features[node.Feature]
	if featureValue < node.Threshold {
		return m.predictTree(features, node.Left)
	}
	return m.predictTree(features, node.Right)
}

// Predict using all trees and average the results
func (m *ForestModel) predict(features map[string]float64) float64 {
	sum := 0.0
	for i := 0; i < len(m.Trees); i++ {
		sum += m.predictTree(features, &m.Trees[i])
	}
	return sum / float64(len(m.Trees))
}

// Prepare input features as a map
func (m *ForestModel) prepareInput(ctx StrategyContext_simple) map[string]float64 {
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

// InvestDecisionMaking_ml_forest uses the random forest model for buy decisions
func InvestDecisionMaking_ml_forest(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadForestModel("forest_model.json")
	if err != nil {
		// Fallback to default strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input
	features := model.prepareInput(ctx)

	// Get prediction (0.0 to 1.0 range representing fraction of funds to spend)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount
	return ctx.Funds * prediction
}
