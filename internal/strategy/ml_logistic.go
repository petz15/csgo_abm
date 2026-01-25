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
// Features expected: ct_score_start, t_score_start, score_diff, ct_starting_equipment,
// t_starting_equipment, ct_money_start, t_money_start
func (m *LogisticRegressionModel) prepareInput(ctx StrategyContext_simple) []float64 {
	features := make([]float64, 7)

	// Feature indices match the model's features array:
	// 0: ct_score_start
	// 1: t_score_start
	// 2: score_diff
	// 3: ct_starting_equipment
	// 4: t_starting_equipment
	// 5: ct_money_start
	// 6: t_money_start

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

// InvestDecisionMaking_ml_logistic uses logistic regression for buy decisions
// Predicts P(CT wins | game state) and uses as investment fraction
func InvestDecisionMaking_ml_logistic(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadLogisticModel("logistic.json")
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
