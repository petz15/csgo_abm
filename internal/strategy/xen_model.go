package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// ============================================================================
// Xen Model Strategy: ML XGBoost Buy Type Optimization
// Loads XGBoost model from xen_model folder and tests different buy types
// to determine the optimal buy that maximizes win probability
// ============================================================================

// BuyTypeDefinition represents a specific buy strategy with its characteristics
type BuyTypeDefinition struct {
	Name                 string
	MinStartingEquipment float64
	MaxStartingEquipment float64
	MinSpend             float64
	MaxSpend             float64
	Description          string
	EncodedValue         int // Encoded value for the model
}

// BuyTypeDefinitions defines all available buy types based on Paper's definitions
var BuyTypeDefinitions = map[string]BuyTypeDefinition{
	"eco": {
		Name:                 "Eco",
		MinStartingEquipment: 0,
		MaxStartingEquipment: 3000,
		MinSpend:             0,
		MaxSpend:             2000,
		Description:          "Minimal buy to save money - starting equipment 0-3k, spent 0-2k",
		EncodedValue:         0,
	},
	"low_buy": {
		Name:                 "Low Buy",
		MinStartingEquipment: 0,
		MaxStartingEquipment: 3000,
		MinSpend:             2000,
		MaxSpend:             7500,
		Description:          "Light buy with weak utility - starting equipment 0-3k, spent 2k-7.5k",
		EncodedValue:         5,
	},
	"half_buy": {
		Name:                 "Half Buy",
		MinStartingEquipment: 0,
		MaxStartingEquipment: 3000,
		MinSpend:             7500,
		MaxSpend:             20000,
		Description:          "Moderate buy with decent utility - starting equipment 0-3k, spent 7.5-20k",
		EncodedValue:         2,
	},
	"hero_low": {
		Name:                 "Hero Low Buy",
		MinStartingEquipment: 3000,
		MaxStartingEquipment: 20000,
		MinSpend:             0,
		MaxSpend:             7500,
		Description:          "Strong economy with weak utility - starting equipment 3-20k, spent 0-7.5k",
		EncodedValue:         4,
	},
	"hero_half": {
		Name:                 "Hero Half Buy",
		MinStartingEquipment: 3000,
		MaxStartingEquipment: 20000,
		MinSpend:             7500,
		MaxSpend:             17000,
		Description:          "Strong economy with solid utility - starting equipment 3-20k, spent 7.5-17k",
		EncodedValue:         3,
	},
	"full_buy": {
		Name:                 "Full Buy",
		MinStartingEquipment: 3000,
		MaxStartingEquipment: 999999,
		MinSpend:             17500,
		MaxSpend:             999999,
		Description:          "Full buy for max firepower - starting equipment + spent > 20k",
		EncodedValue:         1,
	},
}

// XGBoostTree represents a single tree in the XGBoost model
type XGBoostTree_xen struct {
	ID              int64     `json:"id"`
	BaseWeights     []float64 `json:"base_weights"`
	LeftChildren    []int64   `json:"left_children"`
	RightChildren   []int64   `json:"right_children"`
	SplitConditions []float64 `json:"split_conditions"`
	SplitIndices    []int64   `json:"split_indices"`
	DefaultLeft     []int     `json:"default_left"`
}

// XGBoostModel represents the complete XGBoost model loaded from JSON
type XGBoostModel_xen struct {
	Booster struct {
		ModelParam struct {
			NumFeature int `json:"num_feature"`
		} `json:"model_param"`
		GBTree struct {
			Forest []XGBoostTree_xen `json:"forest"`
		} `json:"gbtree"`
		LearnerModelParam struct {
			BaseScore string `json:"base_score"`
		} `json:"learner_model_param"`
	} `json:"booster"`
	Version [3]int `json:"version"`
}

// Singleton pattern for XGBoost model
var (
	xenModelInstance *XGBoostModel_xen
	xenModelOnce     sync.Once
	xenModelErr      error
)

// LoadXenModel loads the XGBoost model from the xen_model folder (singleton pattern)
func LoadXenModel() (*XGBoostModel_xen, error) {
	xenModelOnce.Do(func() {
		// Try multiple possible paths
		paths := []string{
			"xen_model/xen_xgboost_model.json",
			"./xen_model/xen_xgboost_model.json",
			"../xen_model/xen_xgboost_model.json",
		}

		var lastErr error

		for _, path := range paths {
			data, err := os.ReadFile(path)
			if err == nil {
				model := &XGBoostModel_xen{}
				if err := json.Unmarshal(data, model); err == nil {
					xenModelInstance = model
					return
				}
				lastErr = err
			}
			lastErr = err
		}

		xenModelErr = lastErr
	})

	return xenModelInstance, xenModelErr
}

// predictTree recursively traverses a single tree
func (m *XGBoostModel_xen) predictTree_xen(features []float64, tree *XGBoostTree_xen) float64 {
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
func (m *XGBoostModel_xen) predict_xen(features []float64) float64 {
	// Parse base score
	baseScore := 0.5062822746
	prediction := baseScore

	// Add prediction from each tree
	for i := range m.Booster.GBTree.Forest {
		prediction += m.predictTree_xen(features, &m.Booster.GBTree.Forest[i])
	}

	// Apply sigmoid for binary classification (P(CT wins))
	return 1.0 / (1.0 + math.Exp(-prediction))
}

// prepareXenFeatures prepares input features for the Xen model
// Features: ct_score_start, t_score_start, score_diff, ct_starting_equipment,
// t_starting_equipment, ct_money_start, t_money_start, ct_buy_encoded, t_buy_encoded
func prepareXenFeatures(ctx StrategyContext_simple, ctBuyType int, tBuyType int, isCT bool) []float64 {
	features := make([]float64, 9)

	// Determine CT and T scores based on perspective
	ctScore := ctx.OwnScore
	tScore := ctx.OpponentScore
	ctMoney := ctx.Funds
	tMoney := ctx.Funds_opponent_forbidden
	ctEquipment := ctx.Equipment
	tEquipment := ctx.Start_Equipment_opponent_forbidden

	// If this is T side perspective, swap the values
	if !isCT {
		ctScore, tScore = tScore, ctScore
		ctMoney, tMoney = tMoney, ctMoney
		ctEquipment, tEquipment = tEquipment, ctEquipment
	}

	features[0] = float64(ctScore)          // ct_score_start
	features[1] = float64(tScore)           // t_score_start
	features[2] = float64(ctScore - tScore) // score_diff
	features[3] = ctEquipment / 999999.0    // ct_starting_equipment (normalized)
	features[4] = tEquipment / 999999.0     // t_starting_equipment (normalized)
	features[5] = ctMoney / 999999.0        // ct_money_start (normalized)
	features[6] = tMoney / 999999.0         // t_money_start (normalized)
	features[7] = float64(ctBuyType)        // ct_buy_encoded
	features[8] = float64(tBuyType)         // t_buy_encoded

	return features
}

// PredictCTWinProbability uses the XGBoost model to predict P(CT wins) for a given buy configuration
func PredictCTWinProbability(model *XGBoostModel_xen, ctx StrategyContext_simple, ctBuyType int, tBuyType int, isCT bool) float64 {

	features := prepareXenFeatures(ctx, ctBuyType, tBuyType, isCT)
	prediction := model.predict_xen(features)

	// Clamp to valid probability range
	if prediction < 0.0 {
		prediction = 0.0
	} else if prediction > 1.0 {
		prediction = 1.0
	}

	return prediction
}

// CalculateInvestmentForBuyType determines how much to spend to achieve a specific buy type
func CalculateInvestmentForBuyType(buyType string, currentFunds float64, ctx StrategyContext_simple) float64 {
	profile, exists := BuyTypeDefinitions[buyType]
	if !exists {
		return 0.0
	}

	//make sure starting equipment is in range
	if profile.MinStartingEquipment > ctx.Equipment || profile.MaxStartingEquipment < ctx.Equipment {
		return -1 // Not enough starting equipment for this buy type
	}

	//make sure min spend is achievable and set investment = to currentFunds or maxspend
	investment := profile.MaxSpend
	if currentFunds < profile.MinSpend {
		return -1 // Not enough funds for this buy type
	} else if currentFunds < profile.MaxSpend {
		return currentFunds
	}

	return investment
}

// CanAffordBuyType checks if a buy type is achievable with current funds
func CanAffordBuyType(profile BuyTypeDefinition, currentFunds float64, currentEquipment float64) bool {

	if profile.Name == "" {
		return false // Invalid buy type
	}

	// Check starting equipment range
	if currentEquipment < profile.MinStartingEquipment || currentEquipment > profile.MaxStartingEquipment {
		return false
	}

	//return true if minSpend is affordable
	return currentFunds >= profile.MinSpend
}

// InvestDecisionMaking_xen_model evaluates all buy types using XGBoost model
// and selects the optimal one based on predicted win probability
// Uses minimax strategy: picks the buy that performs best in worst case
func InvestDecisionMaking_xen_model(ctx StrategyContext_simple) float64 {
	currentFunds := ctx.Funds

	// Load the XGBoost model
	model_xen, err := LoadXenModel()
	if err != nil {
		panic("failed to load XGBoost model: " + err.Error())
	}

	bestWorstCaseWinProb := 0.0
	bestInvestment := 0.0

	// Determine if we're CT or T side
	isCT := ctx.Side

	// All possible opponent buy types (0-5)
	opponentBuyTypes := []int{0, 1, 2, 3, 4, 5}

	// Evaluate each affordable buy type
	for buyTypeName, profile := range BuyTypeDefinitions {

		// Calculate investment needed and if not enough money, continue
		investment := CalculateInvestmentForBuyType(buyTypeName, currentFunds, ctx)
		if investment < 0 {
			continue
		}

		// Test this buy type against all possible opponent buy types
		// Find the worst-case (minimum) win probability
		worstCaseWinProb := 1.0 // Start at maximum

		for _, opponentBuyType := range opponentBuyTypes {
			// Predict win probability with this buy type vs opponent buy type
			if !CanAffordBuyType(GetBuyTypeByEncodedValue(opponentBuyType), ctx.Funds_opponent_forbidden, ctx.Start_Equipment_opponent_forbidden) {
				continue // Skip opponent buy types that are not affordable
			}

			winProb := PredictCTWinProbability(model_xen, ctx, profile.EncodedValue, opponentBuyType, isCT)

			// For T side, convert to T win probability
			effectiveWinProb := winProb
			if !isCT {
				effectiveWinProb = 1.0 - winProb
			}

			// Track the worst case (minimum) for this buy type
			if effectiveWinProb < worstCaseWinProb {
				worstCaseWinProb = effectiveWinProb
			}
		}

		// Select buy type with the highest worst-case win probability (minimax strategy)
		if worstCaseWinProb > bestWorstCaseWinProb {
			bestWorstCaseWinProb = worstCaseWinProb
			bestInvestment = investment
		}
	}

	// Safety fallback: if nothing selected, do eco
	if bestInvestment <= 0 {
		bestInvestment = math.Min(currentFunds*0.75, 2000) // Eco buy
	}

	return bestInvestment
}

// GetBuyTypeByEncodedValue searches the buy type definitions by encoded value
func GetBuyTypeByEncodedValue(encodedValue int) BuyTypeDefinition {
	for _, profile := range BuyTypeDefinitions {
		if profile.EncodedValue == encodedValue {
			return profile
		}
	}
	return BuyTypeDefinition{}
}
