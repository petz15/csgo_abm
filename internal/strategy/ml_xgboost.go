package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// XGBoostTree represents a single tree in the XGBoost model
type XGBoostTree struct {
	ID              int64     `json:"id"`
	BaseWeights     []float64 `json:"base_weights"`
	LeftChildren    []int64   `json:"left_children"`
	RightChildren   []int64   `json:"right_children"`
	SplitConditions []float64 `json:"split_conditions"`
	SplitIndices    []int64   `json:"split_indices"`
	DefaultLeft     []int     `json:"default_left"`
}

// XGBoostModel represents the complete XGBoost model loaded from JSON
type XGBoostModel struct {
	Booster struct {
		ModelParam struct {
			NumFeature int `json:"num_feature"`
		} `json:"model_param"`
		GBTree struct {
			Forest []XGBoostTree `json:"forest"`
		} `json:"gbtree"`
		LearnerModelParam struct {
			BaseScore string `json:"base_score"`
		} `json:"learner_model_param"`
	} `json:"booster"`
	Version [3]int `json:"version"`
}

var xgboostModelInstance *XGBoostModel
var xgboostOnce sync.Once

// LoadXGBoostModel loads the XGBoost model from JSON file
func LoadXGBoostModel(path string) (*XGBoostModel, error) {
	var err error
	xgboostOnce.Do(func() {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}

		var model XGBoostModel
		if unmarshalErr := json.Unmarshal(data, &model); unmarshalErr != nil {
			err = unmarshalErr
			return
		}

		xgboostModelInstance = &model
	})

	if err != nil {
		return nil, err
	}
	return xgboostModelInstance, nil
}

// predict a single tree recursively
func (m *XGBoostModel) predictTree(features []float64, tree *XGBoostTree) float64 {
	nodeID := int64(0)

	for {
		// Check if leaf node
		if tree.LeftChildren[nodeID] == -1 && tree.RightChildren[nodeID] == -1 {
			return tree.BaseWeights[nodeID]
		}

		// Get split information
		splitIndex := tree.SplitIndices[nodeID]
		if splitIndex < 0 || int(splitIndex) >= len(features) {
			return tree.BaseWeights[nodeID]
		}

		featureValue := features[int(splitIndex)]
		splitCondition := tree.SplitConditions[nodeID]

		// Navigate tree
		if featureValue < splitCondition {
			nextID := tree.LeftChildren[nodeID]
			if nextID == -1 {
				return tree.BaseWeights[nodeID]
			}
			nodeID = nextID
		} else {
			nextID := tree.RightChildren[nodeID]
			if nextID == -1 {
				return tree.BaseWeights[nodeID]
			}
			nodeID = nextID
		}
	}
}

// predict using all trees and sum the results
func (m *XGBoostModel) predict(features []float64) float64 {
	// Parse base score
	baseScore := 0.5062822746
	prediction := baseScore

	// Add prediction from each tree
	for i := range m.Booster.GBTree.Forest {
		prediction += m.predictTree(features, &m.Booster.GBTree.Forest[i])
	}

	// Apply sigmoid for binary classification
	return 1.0 / (1.0 + math.Exp(-prediction))
}

// Prepare input features as a normalized array
func (m *XGBoostModel) prepareInput(ctx StrategyContext_simple) []float64 {
	// Features: own_funds, own_score, opponent_score, own_survivors,
	// opponent_survivors, consecutive_losses, is_ct_side, round_number,
	// half_length, last_round_reason, last_bomb_planted, own_starting_equipment, score_diff
	features := make([]float64, 13)

	features[0] = ctx.Funds / 999999.0
	features[1] = float64(ctx.OwnScore) / 16.0
	features[2] = float64(ctx.OpponentScore) / 16.0
	features[3] = float64(ctx.OwnSurvivors) / 5.0
	features[4] = float64(ctx.EnemySurvivors) / 5.0
	features[5] = float64(ctx.ConsecutiveLosses) / 5.0

	if ctx.Side {
		features[6] = 1.0
	} else {
		features[6] = 0.0
	}

	features[7] = float64(ctx.CurrentRound) / 30.0
	features[8] = float64(ctx.GameRules_strategy.HalfLength) / 15.0
	features[9] = float64(ctx.RoundEndReason) / 4.0

	if ctx.Is_BombPlanted {
		features[10] = 1.0
	} else {
		features[10] = 0.0
	}

	features[11] = ctx.Equipment / 999999.0
	features[12] = (float64(ctx.OwnScore-ctx.OpponentScore) + 15.0) / 30.0

	return features
}

// Prepare input features for forbidden variant (includes opponent info)
func (m *XGBoostModel) prepareInputForbidden(ctx StrategyContext_simple) []float64 {
	baseFeatures := m.prepareInput(ctx)
	// Add opponent features
	features := make([]float64, len(baseFeatures)+2)
	copy(features, baseFeatures)
	features[len(baseFeatures)] = ctx.Funds_opponent_forbidden / 999999.0
	features[len(baseFeatures)+1] = ctx.Start_Equipment_opponent_forbidden / 999999.0
	return features
}

// InvestDecisionMaking_ml_xgboost uses the XGBoost model for buy decisions
func InvestDecisionMaking_ml_xgboost(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, _ := LoadXGBoostModel("ml_models/xgboost_model.json")

	// Prepare normalized input
	features := model.prepareInput(ctx)

	// Get prediction (0.0 to 1.0 probability)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount (fraction of funds to spend)
	return ctx.Funds * prediction
}

// InvestDecisionMaking_ml_xgboost_forbidden uses the XGBoost model with extended features
func InvestDecisionMaking_ml_xgboost_forbidden(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, _ := LoadXGBoostModel("ml_models/xgboost_model_forbidden.json")

	// Prepare normalized input with forbidden features
	features := model.prepareInputForbidden(ctx)

	// Get prediction (0.0 to 1.0 probability)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount (fraction of funds to spend)
	return ctx.Funds * prediction
}
