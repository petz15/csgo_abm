package engine

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// RoundOutcome holds all the results of a round determination
type RoundOutcome struct {
	CTWins                    bool // true if CT wins, false if T wins
	ReasonCode                int
	BombPlanted               bool
	CTSurvivors               int
	TSurvivors                int
	CTEquipmentSharePerPlayer []float64
	TEquipmentSharePerPlayer  []float64
	CTEquipmentPerPlayer      []float64
	TEquipmentPerPlayer       []float64
	CSF                       float64
	CSFKey                    string
	StochasticValues          RNG_Outcomes
}

type RNG_Outcomes struct {
	RNG_CSF          float64
	RNG_RoundOutcome float64
	RNG_Bombplant    float64
	RNG_SurvivorsCT  float64
	RNG_SurvivorsT   float64
	RNG_EquipmentCT  []float64
	RNG_EquipmentT   []float64
}

var (
	distributions       *ProcessedDistributions
	distributionsLoaded bool
)

// LoadDistributions loads the distributions from JSON file and converts to sorted slices.
// Must be called once at application startup before any simulations.
func LoadDistributions(filePath string) error {
	if distributionsLoaded {
		return nil
	}

	processed, err := loadAndProcessDistributions(filePath)
	if err != nil {
		return err
	}

	distributions = processed
	distributionsLoaded = true

	return nil
}

// IsABMModelsLoaded returns whether distributions have been loaded
func IsABMModelsLoaded() bool {
	return distributionsLoaded
}

// GetABMModels returns the loaded distributions (for testing/debugging)
func GetABMModels() *ProcessedDistributions {
	return distributions
}

// ContestSuccessFunction calculates win probability using Tullock CSF.
// Returns the probability that side with expenditure x wins against side with expenditure y.
func CSF(x float64, y float64) float64 {
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

// DetermineRoundOutcome determines all aspects of a round outcome based on CSF probability.
// csfProb should be in [0,1], representing the CT win probability.
func DetermineRoundOutcome(ct_eq_val float64, t_eq_val float64, rng *rand.Rand, gameR GameRules) RoundOutcome {
	assertLoaded("DetermineRoundOutcome")

	outcome := RoundOutcome{}

	outcome.CSF = CSF(ct_eq_val, t_eq_val)

	// 1. Determine winner
	outcome.StochasticValues.RNG_CSF = rng.Float64()
	outcome.CTWins = outcome.StochasticValues.RNG_CSF < outcome.CSF
	side := determineSide(outcome.CTWins)

	outcome.CSFKey = csfKeyForProb(outcome.CSF)

	// 2. Determine round end reason
	sampleRoundEndReason(side, &outcome, rng)

	// 3. Determine bomb planted status
	determineBombPlanted(&outcome, rng)

	if gameR.WithSaves { //to test robustness and certain effects, Survivors and Equipment Saved can be excluded
		// 4. Determine survivors
		winningSide := side
		losingSide := oppositeSide(side)

		winningSurvivors := sampleSurvivors(winningSide, &outcome, rng)
		losingSurvivors := sampleSurvivors(losingSide, &outcome, rng)
		if outcome.CTWins {
			outcome.CTSurvivors = winningSurvivors
			outcome.TSurvivors = losingSurvivors
		} else {
			outcome.TSurvivors = winningSurvivors
			outcome.CTSurvivors = losingSurvivors
		}

		// 5. Determine equipment saved

		total_equipment := ct_eq_val + t_eq_val
		sampleEquipment(winningSide, &outcome, rng)
		sampleEquipment(losingSide, &outcome, rng)

		// 6. Calculate equipment value per surviving player and making sure, players cant save more than total equipment
		determineEquipmentSavedPerPlayer(&outcome, total_equipment)
	} else {
		outcome.CTSurvivors = 0
		outcome.TSurvivors = 0
	}

	return outcome
}

// ============================================================================
// Internal Helpers - Distribution Sampling
// ============================================================================

// sampleRoundEndReason samples the round end reason from distributions
func sampleRoundEndReason(side string, oc *RoundOutcome, rng *rand.Rand) {
	sideMap := distributions.RoundEndReason[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: missing round end reason distributions for side='%s'", side))
	}

	thresholds := sideMap[oc.CSFKey]
	if len(thresholds) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: missing or empty cumulative distribution for side='%s', csfKey='%s'", side, oc.CSFKey))
	}

	randValue := rng.Float64() * 100.0
	oc.StochasticValues.RNG_RoundOutcome = randValue

	// Use pre-sorted slice
	for _, entry := range thresholds {
		if randValue <= entry.Threshold {
			oc.ReasonCode = entry.Reason
			break
		}
	}

	if oc.ReasonCode == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleRoundEndReason: failed to sample reason for side='%s', csfKey='%s', randValue=%.2f", side, oc.CSFKey, randValue))
	}

}

// determineBombPlanted determines if bomb was planted based on reason code
func determineBombPlanted(oc *RoundOutcome, rng *rand.Rand) {
	// Reason codes: 1=Target Bombed, 2=T Win Elimination, 3=CT Win Defuse, 4=CT Win Elimination
	switch oc.ReasonCode {
	case 1, 3:
		oc.BombPlanted = true
	case 2:
		// Sample from bomb planted distribution
		bombMap := distributions.BombPlantedT
		if bombMap == nil {
			panic("probabilities.go: determineBombPlanted: missing bomb planted distribution for T side")
		}
		prob, ok := bombMap[oc.CSFKey]
		if !ok {
			panic(fmt.Sprintf("probabilities.go: determineBombPlanted: missing bomb planted probability for csfKey='%s'", oc.CSFKey))
		}
		oc.StochasticValues.RNG_Bombplant = rng.Float64()
		oc.BombPlanted = oc.StochasticValues.RNG_Bombplant <= prob
	default:
		oc.BombPlanted = false
	}
}

// sampleSurvivors samples number of survivors from distributions
func sampleSurvivors(side string, oc *RoundOutcome, rng *rand.Rand) int {
	sideMap := distributions.Survivors[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing survivor distributions for side='%s'", side))
	}

	reasonData := sideMap[strconv.Itoa(oc.ReasonCode)]
	if reasonData == nil {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing survivor distributions for side='%s', reasonCode='%d'", side, oc.ReasonCode))
	}

	thresholds := reasonData[oc.CSFKey]
	if len(thresholds) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleSurvivors: missing or empty cumulative lookup for side='%s', reasonCode='%d', csfKey='%s'", side, oc.ReasonCode, oc.CSFKey))
	}

	randValue := rng.Float64() * 100.0
	if side == "CT" {
		oc.StochasticValues.RNG_SurvivorsCT = randValue
	} else {
		oc.StochasticValues.RNG_SurvivorsT = randValue
	}

	var survivors int = 0

	// Use pre-sorted slice
	for _, entry := range thresholds {
		if randValue <= entry.Threshold {
			survivors = entry.Value
			break
		}
	}

	return survivors

}

// sampleEquipment samples equipment saved value from distributions
func sampleEquipment(side string, oc *RoundOutcome, rng *rand.Rand) {
	var survivors int = 0

	if side == "CT" {
		survivors = oc.CTSurvivors
	} else {
		survivors = oc.TSurvivors
	}

	if survivors == 0 {
		return
	}

	sideMap := distributions.EquipmentSaved[side]
	if sideMap == nil {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing equipment saved distributions for side='%s'", side))
	}

	reasonCodeStr := strconv.Itoa(oc.ReasonCode)
	reasonData := sideMap[reasonCodeStr]
	if reasonData == nil {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing equipment saved distributions for side='%s', reasonCode='%d'", side, oc.ReasonCode))
	}

	survivorKey := strconv.Itoa(survivors)
	percentiles := reasonData[survivorKey]
	if len(percentiles) == 0 {
		panic(fmt.Sprintf("probabilities.go: sampleEquipment: missing or empty ECDF lookup for side='%s', reasonCode='%d', survivors=%d", side, oc.ReasonCode, survivors))
	}

	saved_eq_pct := 0.0
	for i := 0; i < survivors; i++ {
		if side == "CT" {
			oc.StochasticValues.RNG_EquipmentCT = append(oc.StochasticValues.RNG_EquipmentCT, rng.Float64())
			saved_eq_pct = sampleFromECDFLookup(percentiles, oc.StochasticValues.RNG_EquipmentCT[i])
			oc.CTEquipmentSharePerPlayer = append(oc.CTEquipmentSharePerPlayer, saved_eq_pct)
		} else {
			oc.StochasticValues.RNG_EquipmentT = append(oc.StochasticValues.RNG_EquipmentT, rng.Float64())
			saved_eq_pct = sampleFromECDFLookup(percentiles, oc.StochasticValues.RNG_EquipmentT[i])
			oc.TEquipmentSharePerPlayer = append(oc.TEquipmentSharePerPlayer, saved_eq_pct)
		}
	}

}

func determineEquipmentSavedPerPlayer(oc *RoundOutcome, total_equipment float64) {

	if sumArray(oc.CTEquipmentSharePerPlayer)+sumArray(oc.TEquipmentSharePerPlayer) > 100.0 {
		// Normalize shares
		total_share := sumArray(oc.CTEquipmentSharePerPlayer) + sumArray(oc.TEquipmentSharePerPlayer)
		for i := range oc.CTEquipmentSharePerPlayer {
			eq_value := oc.CTEquipmentSharePerPlayer[i] / total_share
			oc.CTEquipmentPerPlayer = append(oc.CTEquipmentPerPlayer, eq_value)
		}
		for i := range oc.TEquipmentSharePerPlayer {
			eq_value := oc.TEquipmentSharePerPlayer[i] / total_share
			oc.TEquipmentPerPlayer = append(oc.TEquipmentPerPlayer, eq_value)
		}
	} else {

		for _, share := range oc.CTEquipmentSharePerPlayer {
			eq_value := (share / 100) * total_equipment
			oc.CTEquipmentPerPlayer = append(oc.CTEquipmentPerPlayer, eq_value)
		}

		for _, share := range oc.TEquipmentSharePerPlayer {
			eq_value := (share / 100) * total_equipment
			oc.TEquipmentPerPlayer = append(oc.TEquipmentPerPlayer, eq_value)
		}
	}

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

// the csfKey in the distribution is limited as some values do not exist (e.g., 0, 1, 98, 99, 100)
// values which do not exist between min and max have been interpolated when calculating the distributions
// csfKeyForProb converts probability [0,1] to CSF key string [min-max]
func csfKeyForProb(prob float64) string {
	csfPercent := int(math.Round(prob * 100))
	if csfPercent < distributions.Metadata.CSFRanges.Min {
		csfPercent = distributions.Metadata.CSFRanges.Min
	} else if csfPercent > distributions.Metadata.CSFRanges.Max {
		csfPercent = distributions.Metadata.CSFRanges.Max
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

// sampleFromECDFLookup samples from a pre-sorted ECDF lookup table
func sampleFromECDFLookup(percentiles []PercentileValue, RNG_Eq float64) float64 {
	randPercentile := RNG_Eq * 100.0
	result := 0.0

	// Use pre-sorted slice
	for _, entry := range percentiles {
		if entry.Percentile <= randPercentile {
			result = entry.Value
		} else {
			break
		}
	}

	return result
}

func sumArray(arr []float64) float64 {
	sum := 0.0
	for _, val := range arr {
		sum += val
	}
	return sum
}
