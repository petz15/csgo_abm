package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// LogisticRegressionModel represents a logistic regression model loaded from JSON
type LogisticRegressionModel struct {
	ModelType    string    `json:"model_type"`
	ModelName    string    `json:"model_name"`
	Description  string    `json:"description"`
	Coefficients []float64 `json:"coefficients"`
	Intercept    float64   `json:"intercept"`
	Classes      []int     `json:"classes"`
	Features     []string  `json:"features"`
	FeatureCount int       `json:"feature_count"`
	Performance  struct {
		Accuracy float64 `json:"accuracy"`
		ROCAUC   float64 `json:"roc_auc"`
		LogLoss  float64 `json:"log_loss"`
	} `json:"performance"`
	FeatureImportance map[string]float64 `json:"feature_importance"`
}

var logisticModelInstance *LogisticRegressionModel
var logisticOnce sync.Once

// LoadLogisticModel loads the logistic regression model from JSON file
func LoadLogisticModel(path string) (*LogisticRegressionModel, error) {
	var err error
	logisticOnce.Do(func() {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}

		var model LogisticRegressionModel
		if unmarshalErr := json.Unmarshal(data, &model); unmarshalErr != nil {
			err = unmarshalErr
			return
		}

		logisticModelInstance = &model
	})

	if err != nil {
		return nil, err
	}
	return logisticModelInstance, nil
}

// predict using logistic regression formula: sigmoid(intercept + sum(coeff * feature))
func (m *LogisticRegressionModel) predict(features []float64) float64 {
	if len(features) != len(m.Coefficients) {
		// Return neutral prediction if feature count mismatch
		return 0.5
	}

	// Calculate linear combination
	logit := m.Intercept
	for i := 0; i < len(m.Coefficients); i++ {
		logit += m.Coefficients[i] * features[i]
	}

	// Apply sigmoid function: 1 / (1 + exp(-logit))
	return 1.0 / (1.0 + math.Exp(-logit))
}

// Prepare input features as a normalized array
// Now uses the same feature set as other models
func (m *LogisticRegressionModel) prepareInput(ctx StrategyContext_simple) []float64 {
	features := make([]float64, 13)

	// Features: own_funds, own_score, opponent_score, own_survivors,
	// opponent_survivors, consecutive_losses, is_ct_side, round_number,
	// half_length, last_round_reason, last_bomb_planted, own_starting_equipment, score_diff

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
func (m *LogisticRegressionModel) prepareInputForbidden(ctx StrategyContext_simple) []float64 {
	baseFeatures := m.prepareInput(ctx)
	// Add opponent features
	features := make([]float64, len(baseFeatures)+2)
	copy(features, baseFeatures)
	features[len(baseFeatures)] = ctx.Funds_opponent_forbidden / 999999.0
	features[len(baseFeatures)+1] = ctx.Start_Equipment_opponent_forbidden / 999999.0
	return features
}

// InvestDecisionMaking_ml_logistic uses logistic regression for buy decisions
// Predicts P(CT wins | game state) and uses as investment fraction
func InvestDecisionMaking_ml_logistic(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadLogisticModel("ml_models/logistic_model.json")
	if err != nil {
		// Fallback to adaptive strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input
	features := model.prepareInput(ctx)

	// Get prediction (probability CT wins - 0.0 to 1.0)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	if !ctx.Side {
		// If team is T, invert prediction to get P(T wins)
		prediction = 1.0 - prediction
	}

	// Return investment amount based on win probability
	// Higher win probability -> invest more aggressively
	return ctx.Funds * prediction
}

// InvestDecisionMaking_ml_logistic_forbidden uses logistic regression with extended features
func InvestDecisionMaking_ml_logistic_forbidden(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadLogisticModel("ml_models/logistic_model_forbidden.json")
	if err != nil {
		// Fallback to adaptive strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input with forbidden features
	features := model.prepareInputForbidden(ctx)

	// Get prediction (probability CT wins - 0.0 to 1.0)
	prediction := model.predict(features)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	if !ctx.Side {
		// If team is T, invert prediction to get P(T wins)
		prediction = 1.0 - prediction
	}

	// Return investment amount based on win probability
	// Higher win probability -> invest more aggressively
	return ctx.Funds * prediction
}
