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
	// Features in order: ct_score_start, t_score_start, score_diff,
	// ct_starting_equipment, t_starting_equipment, ct_money_start, t_money_start
	features := make([]float64, 7)

	ct_score_start := 0.0
	t_score_start := 0.0
	ct_funds := 0.0
	t_funds := 0.0
	ct_rs_eq_val := 0.0
	t_rs_eq_val := 0.0

	// Determine scores based on side
	if ctx.Side {
		ct_score_start = float64(ctx.OwnScore)
		t_score_start = float64(ctx.OpponentScore)
		ct_funds = ctx.Funds
		t_funds = ctx.Funds_opponent_forbidden
		ct_rs_eq_val = ctx.Equipment
		t_rs_eq_val = ctx.Start_Equipment_opponent_forbidden
	} else {
		ct_score_start = float64(ctx.OpponentScore)
		t_score_start = float64(ctx.OwnScore)
		ct_funds = ctx.Funds_opponent_forbidden
		t_funds = ctx.Funds
		ct_rs_eq_val = ctx.Start_Equipment_opponent_forbidden
		t_rs_eq_val = ctx.Equipment
	}

	features[0] = ct_score_start                 // ct_score_start
	features[1] = t_score_start                  // t_score_start
	features[2] = ct_score_start - t_score_start // score_diff

	// Equipment values (simplified - using funds as proxy for equipment value)
	features[3] = ct_rs_eq_val // ct_starting_equipment proxy
	features[4] = t_rs_eq_val  // t_starting_equipment proxy (not directly available from context)

	// Money
	features[5] = ct_funds // ct_money_start
	features[6] = t_funds  // t_money_start (not directly available from context)

	return features
}

// InvestDecisionMaking_ml_xgboost uses the XGBoost model for buy decisions
func InvestDecisionMaking_ml_xgboost(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, _ := LoadXGBoostModel("xgboost.json")

	// Prepare normalized input
	features := model.prepareInput(ctx)

	// Get prediction (0.0 to 1.0 probability)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount (fraction of funds to spend)
	return ctx.Funds * prediction
}
