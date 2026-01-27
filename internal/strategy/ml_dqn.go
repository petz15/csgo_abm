package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// GameState represents the observable game state
type GameState struct {
	OwnFunds          float64
	OwnScore          int
	OpponentScore     int
	OwnSurvivors      int
	OpponentSurvivors int
	ConsecutiveLosses int
	IsCTSide          bool
	RoundNumber       int
	HalfLength        int
	LastRoundReason   int
	LastBombPlanted   bool
	OwnEquipment      float64
	ScoreDiff         int
	OpponentFunds     float64
	OpponentEquipment float64
}

func InvestDecisionMaking_ml_dqn(ctx StrategyContext_simple) float64 {
	// Example usage
	var model *DQNModel
	var modelErr error
	var modelOnce sync.Once

	modelOnce.Do(func() {
		model, modelErr = LoadModel("ml_models/metadata.json", "ml_models/q_network_weights.json")
	})

	if modelErr != nil {
		panic(modelErr)
	}

	state := GameState{
		OwnFunds:          ctx.Funds,
		OwnScore:          ctx.OwnScore,
		OpponentScore:     ctx.OpponentScore,
		OwnSurvivors:      ctx.OwnSurvivors,
		OpponentSurvivors: ctx.EnemySurvivors,
		ConsecutiveLosses: ctx.ConsecutiveLosses,
		IsCTSide:          ctx.Side,
		RoundNumber:       ctx.CurrentRound,
		HalfLength:        ctx.GameRules_strategy.HalfLength,
		LastRoundReason:   ctx.RoundEndReason,
		LastBombPlanted:   ctx.Is_BombPlanted,
		OwnEquipment:      ctx.Equipment,
		ScoreDiff:         ctx.OwnScore - ctx.OpponentScore,
	}

	action := model.SelectAction(state)
	return action * ctx.Funds
}

// ToArray converts GameState to normalized feature array
func (s *GameState) ToArray() []float64 {
	ctSide := 0.0
	if s.IsCTSide {
		ctSide = 1.0
	}
	bombPlanted := 0.0
	if s.LastBombPlanted {
		bombPlanted = 1.0
	}

	return []float64{
		s.OwnFunds / 50000.0,
		float64(s.OwnScore) / 16.0,
		float64(s.OpponentScore) / 16.0,
		float64(s.OwnSurvivors) / 5.0,
		float64(s.OpponentSurvivors) / 5.0,
		math.Min(float64(s.ConsecutiveLosses), 5) / 5.0,
		ctSide,
		float64(s.RoundNumber) / 30.0,
		float64(s.HalfLength) / 15.0,
		float64(s.LastRoundReason) / 4.0,
		bombPlanted,
		s.OwnEquipment / 999999.0,            // own_starting_equipment
		(float64(s.ScoreDiff) + 15.0) / 30.0, // score_diff
	}
}

// ToArrayForbidden converts GameState to extended feature array including opponent info
func (s *GameState) ToArrayForbidden(opponentFunds float64, opponentEquipment float64) []float64 {
	baseArray := s.ToArray()
	return append(baseArray,
		opponentFunds/999999.0,     // opponent_funds
		opponentEquipment/999999.0, // opponent_starting_equipment
	)
}

// DQNModel represents a trained DQN model
type DQNModel struct {
	StateDim     int       `json:"state_dim"`
	NActions     int       `json:"n_actions"`
	ActionValues []float64 `json:"action_values"`
	Weights      ModelWeights
}

// ModelWeights holds the neural network weights
type ModelWeights struct {
	Layers []LayerWeights
}

// LayerWeights represents weights for one layer
type LayerWeights struct {
	Weight [][]float64
	Bias   []float64
}

// LoadModel loads a DQN model from JSON files
func LoadModel(metadataPath, weightsPath string) (*DQNModel, error) {
	// Load metadata
	metadataFile, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var model DQNModel
	if err := json.Unmarshal(metadataFile, &model); err != nil {
		return nil, err
	}

	// Load weights
	weightsFile, err := os.ReadFile(weightsPath)
	if err != nil {
		return nil, err
	}

	var rawWeights map[string]interface{}
	if err := json.Unmarshal(weightsFile, &rawWeights); err != nil {
		return nil, err
	}

	// Parse weights into layer structure
	// This is simplified - actual implementation would need to parse the layer structure

	return &model, nil
}

// LayerNorm normalization (simplified)
func layerNorm(x []float64) []float64 {
	mean := 0.0
	for _, v := range x {
		mean += v
	}
	mean /= float64(len(x))

	variance := 0.0
	for _, v := range x {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(x))

	result := make([]float64, len(x))
	for i, v := range x {
		result[i] = (v - mean) / math.Sqrt(variance+1e-5)
	}
	return result
}

// Forward pass through linear layer
func linearForward(input []float64, weight [][]float64, bias []float64) []float64 {
	output := make([]float64, len(bias))
	for i := range output {
		sum := bias[i]
		for j := range input {
			sum += input[j] * weight[i][j]
		}
		output[i] = sum
	}
	return output
}

// Predict Q-values for a given state
func (m *DQNModel) Predict(state GameState) []float64 {
	x := state.ToArray()

	// Forward pass through network
	// This is simplified - actual implementation would iterate through layers
	for _, layer := range m.Weights.Layers {
		x = linearForward(x, layer.Weight, layer.Bias)
		// Apply activation (ReLU) except for last layer
		for i := range x {
			x[i] = relu(x[i])
		}
		x = layerNorm(x)
	}

	return x
}

// SelectAction selects the best action based on Q-values
func (m *DQNModel) SelectAction(game GameState) float64 {
	qValues := m.Predict(game)

	// Find action with maximum Q-value
	maxIdx := 0
	maxVal := qValues[0]
	for i := 1; i < len(qValues); i++ {
		if qValues[i] > maxVal {
			maxVal = qValues[i]
			maxIdx = i
		}
	}

	return m.ActionValues[maxIdx]
}

// InvestDecisionMaking_ml_dqn_forbidden uses DQN with extended feature set including opponent info
func InvestDecisionMaking_ml_dqn_forbidden(ctx StrategyContext_simple) float64 {
	// Example usage
	var model *DQNModel
	var modelErr error
	var modelOnce sync.Once

	modelOnce.Do(func() {
		model, modelErr = LoadModel("ml_models/metadata.json", "ml_models/q_network_weights_forbidden.json")
	})

	if modelErr != nil {
		panic(modelErr)
	}

	state := GameState{
		OwnFunds:          ctx.Funds,
		OwnScore:          ctx.OwnScore,
		OpponentScore:     ctx.OpponentScore,
		OwnSurvivors:      ctx.OwnSurvivors,
		OpponentSurvivors: ctx.EnemySurvivors,
		ConsecutiveLosses: ctx.ConsecutiveLosses,
		IsCTSide:          ctx.Side,
		RoundNumber:       ctx.CurrentRound,
		HalfLength:        ctx.GameRules_strategy.HalfLength,
		LastRoundReason:   ctx.RoundEndReason,
		LastBombPlanted:   ctx.Is_BombPlanted,
		OwnEquipment:      ctx.Equipment,
		ScoreDiff:         ctx.OwnScore - ctx.OpponentScore,
		OpponentFunds:     ctx.Funds_opponent_forbidden,
		OpponentEquipment: ctx.Start_Equipment_opponent_forbidden,
	}

	// For forbidden variant, use extended feature set with opponent info
	action := model.SelectAction(state)
	return action * ctx.Funds
}
