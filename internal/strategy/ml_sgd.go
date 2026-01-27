package strategy

import (
	"encoding/json"
	"math"
	"os"
)

// SGDModel represents a neural network model loaded from JSON
type SGDModel struct {
	Architecture struct {
		InputSize    int   `json:"input_size"`
		HiddenLayers []int `json:"hidden_layers"`
		OutputSize   int   `json:"output_size"`
	} `json:"architecture"`
	Weights struct {
		Layer0Weight [][]float64 `json:"0.weight"`
		Layer0Bias   []float64   `json:"0.bias"`
		Layer2Weight [][]float64 `json:"2.weight"`
		Layer2Bias   []float64   `json:"2.bias"`
		Layer4Weight [][]float64 `json:"4.weight"`
		Layer4Bias   []float64   `json:"4.bias"`
	} `json:"weights"`
	StateFeatures []string           `json:"state_features"`
	Normalization map[string]float64 `json:"normalization"`
}

var sgdModelInstance *SGDModel

// LoadSGDModel loads the SGD neural network model from JSON file
func LoadSGDModel(path string) (*SGDModel, error) {
	if sgdModelInstance != nil {
		return sgdModelInstance, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var model SGDModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	sgdModelInstance = &model
	return sgdModelInstance, nil
}

// Forward pass through the neural network
func (m *SGDModel) predict(input []float64) float64 {
	// Layer 0: input -> hidden1 (64 neurons)
	hidden1 := make([]float64, len(m.Weights.Layer0Bias))
	for i := 0; i < len(hidden1); i++ {
		sum := m.Weights.Layer0Bias[i]
		for j := 0; j < len(input); j++ {
			sum += m.Weights.Layer0Weight[i][j] * input[j]
		}
		hidden1[i] = relu(sum)
	}

	// Layer 2: hidden1 -> hidden2 (32 neurons)
	hidden2 := make([]float64, len(m.Weights.Layer2Bias))
	for i := 0; i < len(hidden2); i++ {
		sum := m.Weights.Layer2Bias[i]
		for j := 0; j < len(hidden1); j++ {
			sum += m.Weights.Layer2Weight[i][j] * hidden1[j]
		}
		hidden2[i] = relu(sum)
	}

	// Layer 4: hidden2 -> output (1 neuron)
	output := m.Weights.Layer4Bias[0]
	for j := 0; j < len(hidden2); j++ {
		output += m.Weights.Layer4Weight[0][j] * hidden2[j]
	}

	return output
}

// Normalize and prepare input features
func (m *SGDModel) prepareInput(ctx StrategyContext_simple) []float64 {
	input := make([]float64, 13)

	// Normalize features according to the model's normalization
	input[0] = ctx.Funds / m.Normalization["own_funds"]
	input[1] = float64(ctx.OwnScore) / m.Normalization["own_score"]
	input[2] = float64(ctx.OpponentScore) / m.Normalization["opponent_score"]
	input[3] = float64(ctx.OwnSurvivors) / m.Normalization["own_survivors"]
	input[4] = float64(ctx.EnemySurvivors) / m.Normalization["opponent_survivors"]
	input[5] = float64(ctx.ConsecutiveLosses) / m.Normalization["consecutive_losses"]

	// is_ct_side: 1.0 if CT, 0.0 if T
	if ctx.Side {
		input[6] = 1.0
	} else {
		input[6] = 0.0
	}

	input[7] = float64(ctx.CurrentRound) / m.Normalization["round_number"]
	input[8] = float64(ctx.GameRules_strategy.HalfLength) / m.Normalization["half_length"]
	input[9] = float64(ctx.RoundEndReason) / m.Normalization["last_round_reason"]

	// last_bomb_planted: 1.0 if planted, 0.0 otherwise
	if ctx.Is_BombPlanted {
		input[10] = 1.0
	} else {
		input[10] = 0.0
	}

	// New features
	input[11] = (float64(ctx.OwnScore-ctx.OpponentScore) + 15.0) / 30.0 // score_diff
	input[12] = ctx.Equipment / 999999.0                                // equipment

	return input
}

// Normalize and prepare input features for forbidden variant (includes opponent info)
func (m *SGDModel) prepareInputForbidden(ctx StrategyContext_simple) []float64 {
	baseInput := m.prepareInput(ctx)
	input := make([]float64, len(baseInput)+2)
	copy(input, baseInput)
	// Add opponent features
	input[len(baseInput)] = ctx.Funds_opponent_forbidden / 999999.0
	input[len(baseInput)+1] = ctx.Start_Equipment_opponent_forbidden / 999999.0
	return input
}

// InvestDecisionMaking_ml_sgd uses the SGD neural network model for buy decisions
func InvestDecisionMaking_ml_sgd(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadSGDModel("ml_models/sgd_model.json")
	if err != nil {
		// Fallback to default strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input
	input := model.prepareInput(ctx)

	// Get prediction (0.0 to 1.0 range representing fraction of funds to spend)
	prediction := model.predict(input)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount
	return ctx.Funds * prediction
}

// InvestDecisionMaking_ml_sgd_forbidden uses the SGD neural network model with extended features
func InvestDecisionMaking_ml_sgd_forbidden(ctx StrategyContext_simple) float64 {
	// Load model (cached after first load)
	model, err := LoadSGDModel("ml_models/sgd_model_forbidden.json")
	if err != nil {
		// Fallback to default strategy if model fails to load
		return InvestDecisionMaking_adaptive_v2(ctx)
	}

	// Prepare normalized input with forbidden features
	input := model.prepareInputForbidden(ctx)

	// Get prediction (0.0 to 1.0 range representing fraction of funds to spend)
	prediction := model.predict(input)

	// Clamp prediction to valid range
	prediction = math.Max(0.0, math.Min(1.0, prediction))

	// Return investment amount
	return ctx.Funds * prediction
}
