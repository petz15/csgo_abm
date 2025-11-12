package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

// ============================================================================
// Data Structures
// ============================================================================

// Distributions represents the complete JSON structure for ABM distributions
type Distributions struct {
	Metadata struct {
		CSFRValue float64 `json:"csf_r_value"`
	} `json:"metadata"`

	Distributions struct {
		RoundEndReason map[string]map[string]struct {
			CumulativeDistribution map[string]struct {
				Reason                int     `json:"reason"`
				ReasonName            string  `json:"reason_name"`
				CumulativeProbability float64 `json:"cumulative_probability"`
			} `json:"cumulative_distribution"`
			NRounds int `json:"n_rounds"`
		} `json:"round_end_reason"`

		BombPlanted struct {
			T map[string]float64 `json:"T"`
		} `json:"bomb_planted"`

		Survivors map[string]map[string]struct {
			ReasonName       string `json:"reason_name"`
			CSFDistributions map[string]struct {
				CumulativeLookup map[string]int `json:"cumulative_lookup"`
			} `json:"csf_distributions"`
		} `json:"survivors"`

		EquipmentSaved map[string]map[string]struct {
			ReasonName            string `json:"reason_name"`
			SurvivorDistributions map[string]struct {
				ECDFLookup map[string]float64 `json:"ecdf_lookup"`
			} `json:"survivor_distributions"`
		} `json:"equipment_saved"`
	} `json:"distributions"`
}

// RoundOutcome holds all the results of a round determination
type RoundOutcome struct {
	CTWins               bool
	ReasonCode           string
	ReasonName           string
	BombPlanted          bool
	CTSurvivors          int
	TSurvivors           int
	CTEquipmentSaved     float64
	CTEquipmentPerPlayer float64
	TEquipmentSaved      float64
	TEquipmentPerPlayer  float64
}

// ============================================================================
// Global State
// ============================================================================

var (
	distributions       *Distributions
	distributionsLoaded bool
)

// ============================================================================
// Public API - Initialization
// ============================================================================

// LoadABMModels loads the distributions from JSON file.
// Must be called once at application startup before any simulations.
func LoadABMModels(filePath string) error {
	if distributionsLoaded {
		return nil
	}

	if filePath == "" {
		filePath = "distributions.json"
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if filePath == "distributions.json" {
			parentPath := filepath.Join("..", "distributions.json")
			data, err = os.ReadFile(parentPath)
			if err != nil {
				return fmt.Errorf("failed to read distributions file from '%s' or '%s': %w", filePath, parentPath, err)
			}
		} else {
			return fmt.Errorf("failed to read distributions file '%s': %w", filePath, err)
		}
	}

	var models Distributions
	if err := json.Unmarshal(data, &models); err != nil {
		return fmt.Errorf("failed to unmarshal distributions file: %w", err)
	}

	distributions = &models
	distributionsLoaded = true

	return nil
}

// IsABMModelsLoaded returns whether distributions have been loaded
func IsABMModelsLoaded() bool {
	return distributionsLoaded
}

// GetABMModels returns the loaded distributions (for testing/debugging)
func GetABMModels() *Distributions {
	return distributions
}

// ============================================================================
// Public API - Contest Success Function
// ============================================================================

// ContestSuccessFunction_simples calculates win probability using Tullock CSF.
// Returns the probability that side with expenditure x wins against side with expenditure y.
func ContestSuccessFunction_simples(x, y float64) float64 {
	r := GetCSFRValue()
	return math.Pow(x, r) / (math.Pow(x, r) + math.Pow(y, r))
}

// GetCSFRValue returns the CSF r value from loaded distributions.
func GetCSFRValue() float64 {
	assertLoaded("GetCSFRValue")
	if distributions.Metadata.CSFRValue <= 0 {
		panic("probabilities.go: GetCSFRValue: CSF r value is invalid or not set in distributions metadata")
	}
	return distributions.Metadata.CSFRValue
}

// ============================================================================
// Public API - Round Outcome Determination
// ============================================================================

// DetermineRoundOutcome determines all aspects of a round outcome based on CSF probability.
// csfProb should be in [0,1], representing the CT win probability.
func DetermineRoundOutcome(csfProb float64) RoundOutcome {
	assertLoaded("DetermineRoundOutcome")

	outcome := RoundOutcome{}

	// 1. Determine winner
	outcome.CTWins = rand.Float64() < csfProb

	// 2. Determine round end reason
	side := determineSide(outcome.CTWins)
	csfKey := csfKeyForProb(csfProb)

	reasonData := sampleRoundEndReason(side, csfKey)
	outcome.ReasonCode = strconv.Itoa(reasonData.reason)
	outcome.ReasonName = reasonData.reasonName

	// 3. Determine bomb planted status
	outcome.BombPlanted = determineBombPlanted(outcome.ReasonCode, csfKey)

	// 4. Determine survivors
	winningSide := side
	losingSide := oppositeSide(side)

	winningSurvivors := sampleSurvivors(winningSide, outcome.ReasonCode, csfKey)
	losingSurvivors := sampleSurvivors(losingSide, outcome.ReasonCode, csfKey)

	if outcome.CTWins {
		outcome.CTSurvivors = winningSurvivors
		outcome.TSurvivors = losingSurvivors
	} else {
		outcome.TSurvivors = winningSurvivors
		outcome.CTSurvivors = losingSurvivors
	}

	// 5. Determine equipment saved
	outcome.CTEquipmentSaved = sampleEquipment("CT", outcome.ReasonCode, outcome.CTSurvivors)
	if outcome.CTSurvivors > 0 {
		outcome.CTEquipmentPerPlayer = outcome.CTEquipmentSaved / float64(outcome.CTSurvivors)
	}

	outcome.TEquipmentSaved = sampleEquipment("T", outcome.ReasonCode, outcome.TSurvivors)
	if outcome.TSurvivors > 0 {
		outcome.TEquipmentPerPlayer = outcome.TEquipmentSaved / float64(outcome.TSurvivors)
	}

	return outcome
}

// ============================================================================
// Internal Helpers - Distribution Sampling
// ============================================================================

type roundEndReasonData struct {
	reason     int
	reasonName string
}

// sampleRoundEndReason samples the round end reason from distributions
func sampleRoundEndReason(side, csfKey string) roundEndReasonData {
	sideMap := distributions.Distributions.RoundEndReason[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: missing round end reason distributions for side='%s'", side))
	}

	bucket := sideMap[csfKey]
	if bucket.CumulativeDistribution == nil || len(bucket.CumulativeDistribution) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: missing or empty cumulative distribution for side='%s', csfKey='%s'", side, csfKey))
	}

	randValue := rand.Float64() * 100.0
	randValue = math.Min(randValue, 99.0) // Ensure it doesn't hit 100.0 exactly, which can cause issues
	//TODO: fix the issues for 100 value
	var selected roundEndReasonData
	minThreshold := math.MaxFloat64

	for thresholdStr, entry := range bucket.CumulativeDistribution {
		threshold, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: invalid threshold string '%s' for side='%s', csfKey='%s': %v", thresholdStr, side, csfKey, err))
		}

		cumulativeThreshold := entry.CumulativeProbability * 100.0
		if randValue <= cumulativeThreshold && threshold < minThreshold {
			selected.reason = entry.Reason
			selected.reasonName = entry.ReasonName
			minThreshold = threshold
		}
	}

	if selected.reasonName == "" {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: failed to sample reason for side='%s', csfKey='%s', randValue=%.2f", side, csfKey, randValue))
	}

	return selected
}

// determineBombPlanted determines if bomb was planted based on reason code
func determineBombPlanted(reasonCode, csfKey string) bool {
	// Reason codes: 1=Target Bombed, 2=T Win Elimination, 3=CT Win Defuse, 4=CT Win Elimination
	switch reasonCode {
	case "1", "3":
		return true
	case "2":
		// Sample from bomb planted distribution
		bombMap := distributions.Distributions.BombPlanted.T
		if bombMap == nil {
			panic(fmt.Sprintf("probabilities.go: determineBombPlanted: missing bomb planted distribution for T side"))
		}
		prob, ok := bombMap[csfKey]
		if !ok {
			panic(fmt.Sprintf("probabilities.go: determineBombPlanted: missing bomb planted probability for csfKey='%s'", csfKey))
		}
		return rand.Float64() < prob
	default:
		return false
	}
}

// sampleSurvivors samples number of survivors from distributions
func sampleSurvivors(side, reasonCode, csfKey string) int {
	sideMap := distributions.Distributions.Survivors[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing survivor distributions for side='%s'", side))
	}

	reasonData := sideMap[reasonCode]
	if reasonData.CSFDistributions == nil {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing survivor distributions for side='%s', reasonCode='%s'", side, reasonCode))
	}

	bucket := reasonData.CSFDistributions[csfKey]
	if bucket.CumulativeLookup == nil || len(bucket.CumulativeLookup) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing or empty cumulative lookup for side='%s', reasonCode='%s', csfKey='%s'", side, reasonCode, csfKey))
	}

	return sampleFromCumulativeLookup(bucket.CumulativeLookup)
}

// sampleEquipment samples equipment saved value from distributions
func sampleEquipment(side, reasonCode string, survivors int) float64 {
	if survivors == 0 {
		return 0.0
	}

	sideMap := distributions.Distributions.EquipmentSaved[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing equipment saved distributions for side='%s'", side))
	}

	reasonData := sideMap[reasonCode]
	if reasonData.SurvivorDistributions == nil {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing equipment saved distributions for side='%s', reasonCode='%s'", side, reasonCode))
	}

	survivorKey := strconv.Itoa(survivors)
	survivorDist := reasonData.SurvivorDistributions[survivorKey]
	if survivorDist.ECDFLookup == nil || len(survivorDist.ECDFLookup) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing or empty ECDF lookup for side='%s', reasonCode='%s', survivors=%d", side, reasonCode, survivors))
	}

	return sampleFromECDFLookup(survivorDist.ECDFLookup)
}

// ============================================================================
// Internal Helpers - Utilities
// ============================================================================

// assertLoaded panics if distributions are not loaded
func assertLoaded(functionName string) {
	if !distributionsLoaded || distributions == nil {
		panic(fmt.Sprintf("probabilities.go: %s: distributions not loaded - call LoadABMModels() before running simulations", functionName))
	}
}

// csfKeyForProb converts probability [0,1] to CSF key string [2-97]
func csfKeyForProb(prob float64) string {
	csfPercent := int(math.Round(prob * 100))
	if csfPercent < 2 {
		csfPercent = 2
	}
	if csfPercent > 97 {
		csfPercent = 97
	}
	return strconv.Itoa(csfPercent)
}

// determineSide returns "CT" or "T" based on CT win status
func determineSide(ctWins bool) string {
	if ctWins {
		return "CT"
	}
	return "T"
}

// oppositeSide returns the opposite side
func oppositeSide(side string) string {
	if side == "CT" {
		return "T"
	}
	return "CT"
}

// sampleFromCumulativeLookup samples from a cumulative distribution lookup table
func sampleFromCumulativeLookup(lookup map[string]int) int {
	randValue := int(math.Floor(rand.Float64() * 100))
	bestThreshold := -1
	result := 0

	for thresholdStr, value := range lookup {
		threshold, err := strconv.Atoi(thresholdStr)
		if err != nil {
			continue
		}
		if threshold <= randValue && threshold > bestThreshold {
			bestThreshold = threshold
			result = value
		}
	}

	return result
}

// sampleFromECDFLookup samples from an ECDF lookup table
func sampleFromECDFLookup(lookup map[string]float64) float64 {
	randPercentile := rand.Float64() * 100.0
	bestPercentile := -1.0
	result := 0.0

	for percentileStr, value := range lookup {
		percentile, err := strconv.ParseFloat(percentileStr, 64)
		if err != nil {
			continue
		}
		if percentile <= randPercentile && percentile > bestPercentile {
			bestPercentile = percentile
			result = value
		}
	}

	return result
}
