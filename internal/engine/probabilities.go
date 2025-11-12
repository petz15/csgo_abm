package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

var csf_r = 1.0855

// RoundEndReason represents a possible round end outcome
type RoundEndReason struct {
	WinnerSide     string  `json:"winner_side"`
	RoundEndReason int     `json:"round_end_reason"`
	Reason         string  `json:"reason"`
	Count          int     `json:"count"`
	Percentage     float64 `json:"percentage"`
}

// RoundEndDistribution holds the loaded distribution data
type RoundEndDistribution struct {
	CTReasons []RoundEndReason
	TReasons  []RoundEndReason
}

var roundEndDist *RoundEndDistribution

// LoadRoundEndDistribution loads the round end distribution from JSON file
func LoadRoundEndDistribution(filePath string) error {
	if filePath == "" {
		// Use default path relative to the probabilities package
		filePath = filepath.Join("internal", "engine", "probabilities", "round_end_dist.json")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var reasons []RoundEndReason
	if err := json.Unmarshal(data, &reasons); err != nil {
		return err
	}

	// Separate CT and T reasons
	dist := &RoundEndDistribution{
		CTReasons: make([]RoundEndReason, 0),
		TReasons:  make([]RoundEndReason, 0),
	}

	for _, reason := range reasons {
		if reason.WinnerSide == "CT" {
			dist.CTReasons = append(dist.CTReasons, reason)
		} else if reason.WinnerSide == "T" {
			dist.TReasons = append(dist.TReasons, reason)
		}
	}

	roundEndDist = dist
	return nil
}

// GetRoundEndDist returns the loaded distribution (loads it if not already loaded)
func GetRoundEndDist() *RoundEndDistribution {
	if roundEndDist == nil {
		// Attempt to load from default location
		_ = LoadRoundEndDistribution("")
	}
	return roundEndDist
}

// DetermineRoundEndReason selects a round end reason based on the winner side
// Returns the reason string (e.g., "Elimination", "Defused", "Exploded")
func DetermineRoundEndReason(ctWin bool) string {
	dist := GetRoundEndDist()
	if dist == nil {
		// Fallback to simple logic if distribution not loaded
		if ctWin {
			// CT win: 70% elimination, 30% defused
			if rand.Float64() < 0.70 {
				return "Elimination"
			}
			return "Defused"
		} else {
			// T win: 68% elimination, 32% exploded
			if rand.Float64() < 0.68 {
				return "Elimination"
			}
			return "Exploded"
		}
	}

	// Use loaded distribution
	var reasons []RoundEndReason
	if ctWin {
		reasons = dist.CTReasons
	} else {
		reasons = dist.TReasons
	}

	if len(reasons) == 0 {
		return "Elimination" // Fallback
	}

	// Select based on percentage distribution
	randValue := rand.Float64() * 100.0
	cumulative := 0.0

	for _, reason := range reasons {
		cumulative += reason.Percentage
		if randValue <= cumulative {
			return reason.Reason
		}
	}

	// Fallback to first reason if something goes wrong
	return reasons[0].Reason
}

// ContestSuccessFunction calculates the probability of winning a round based on team expenditures and other factors.
// This is a Tullock contest success function
func ContestSuccessFunction_simples(x float64, y float64, r float64) float64 {
	probability := (math.Pow(x, r) / (math.Pow(x, r) + math.Pow(y, r)))
	return probability
}

func bool_CSF_simple(x float64, y float64, r float64) bool {
	probability := ContestSuccessFunction_simples(x, y, r)
	return rand.Float64() < probability
}

// CSFNormalDistribution_std_4 generates a value from a normal distribution with mean determined by Contest Success Function
// x, y: The two values to compare (e.g., team expenditures)
// r: The parameter for the CSF (higher r = more deterministic outcomes)
// minOutput, maxOutput: Range for the output values
// Returns a value sampled from a normal distribution between minOutput and maxOutput
func CSFNormalDistribution_std_4(x float64, y float64, r float64, minOutput float64, maxOutput float64) float64 {
	// Calculate the range width for scaling
	rangeWidth := maxOutput - minOutput

	// First calculate the base probability using Contest Success Function (Tullock contest)
	csfProb := 0.5
	if x+y > 0 {
		csfProb = ContestSuccessFunction_simples(x, y, r)
	}

	// Calculate the mean position between minOutput and maxOutput based on CSF probability
	mean := minOutput + csfProb*rangeWidth

	// Set stdDev to 1/4 of the range width as specified
	stdDev := rangeWidth / 4.0

	// Generate a random value from normal distribution with these parameters
	// We'll use the inverse CDF (quantile function) method with a uniform random number
	u := rand.Float64() // Random value between 0 and 1

	// Convert to normal distribution using inverse error function
	// z is a standard normal random variable (mean 0, stdDev 1)
	z := math.Sqrt(2) * math.Erfinv(2*u-1)

	// Scale and shift z to our desired mean and stdDev
	result := mean + z*stdDev

	// Clamp the result to be within minOutput and maxOutput
	return math.Max(minOutput, math.Min(maxOutput, result))
}

// CSFNormalDistribution_std_custom_skew generates a value from a normal distribution with mean determined by Contest Success Function,
// with optional skewness applied to the distribution.
// x, y: The two values to compare (e.g., team expenditures)
// r: The parameter for the CSF (higher r = more deterministic outcomes)
// minOutput, maxOutput: Range for the output values
// stdDevFactor: Controls the standard deviation (higher = less variance)
// skew: Skewness factor (-1 for left, 0 for normal, +1 for right, can be any float)
// Returns a value sampled from a skewed normal distribution between minOutput and maxOutput
func CSFNormalDistribution_std_custom_skew(x float64, y float64, r float64, minOutput float64, maxOutput float64, stdDevFactor float64, skew float64) float64 {
	rangeWidth := maxOutput - minOutput

	csfProb := 0.5
	if x+y > 0 {
		csfProb = ContestSuccessFunction_simples(x, y, r)
	}

	mean := minOutput + csfProb*rangeWidth
	stdDev := rangeWidth / stdDevFactor

	u := rand.Float64()

	// Apply skewness using the Azzalini's skew-normal transformation
	// https://en.wikipedia.org/wiki/Skew_normal_distribution
	// alpha controls the skewness: negative = left, positive = right
	alpha := skew * 5 // scale skew to make it more pronounced

	// Generate two independent standard normal variables
	z0 := math.Sqrt(2) * math.Erfinv(2*u-1)
	z1 := math.Sqrt(2) * math.Erfinv(2*rand.Float64()-1)

	// Skew-normal variable
	skewedZ := alpha*math.Abs(z0) + z1

	result := mean + skewedZ*stdDev

	return math.Max(minOutput, math.Min(maxOutput, result))
}

// --- ABM Models loader and samplers ---

// ABMModels is a trimmed representation of the large JSON used for sampling
type ABMModels struct {
	Metadata struct {
		CSFRValue float64 `json:"csf_r_value"`
	} `json:"metadata"`

	RoundEndReasonDistributions map[string]map[string]struct {
		NRounds            int `json:"n_rounds"`
		ReasonDistribution map[string]struct {
			Count       int     `json:"count"`
			Probability float64 `json:"probability"`
			ReasonName  string  `json:"reason_name"`
		} `json:"reason_distribution"`
	} `json:"round_end_reason_distributions"`

	BombPlantedDistributions map[string]map[string]struct {
		N             int `json:"n_rounds"`
		Probabilities map[string]struct {
			Count       int     `json:"count"`
			Probability float64 `json:"probability"`
		} `json:"probability_distribution"`
	} `json:"bomb_planted_distributions"`

	SurvivorDistributions map[string]map[string]map[string]struct {
		NSamples               int     `json:"n_samples"`
		MeanSurvivors          float64 `json:"mean_survivors"`
		MedianSurvivors        float64 `json:"median_survivors"`
		CumulativeDistribution map[string]struct {
			Count                 int     `json:"count"`
			Probability           float64 `json:"probability"`
			CumulativeProbability float64 `json:"cumulative_probability"`
		} `json:"cumulative_distribution"`
	} `json:"survivor_distributions"`

	EquipmentSavedDistributions map[string]map[string]map[string]struct {
		NSamples int     `json:"n_samples"`
		Mean     float64 `json:"mean"`
		Std      float64 `json:"std"`
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
	} `json:"equipment_saved_distributions"`
}

var abmModels *ABMModels

// LoadABMModels loads the large ABM JSON file and stores it for sampling
func LoadABMModels(filePath string) error {
	if filePath == "" {
		filePath = filepath.Join("internal", "engine", "probabilities", "abm_models.json")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var models ABMModels
	if err := json.Unmarshal(data, &models); err != nil {
		return fmt.Errorf("unmarshal abm models: %w", err)
	}

	abmModels = &models

	// Apply csf r value if present
	if abmModels.Metadata.CSFRValue > 0 {
		csf_r = abmModels.Metadata.CSFRValue
	}

	return nil
}

// csfPercentKeyForProb takes a map (any map with string keys) and a probability p in [0,1]
// It returns the string key whose numeric value is closest to p*100
func csfPercentKeyForProb(m interface{}, p float64) string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return "50"
	}
	keys := v.MapKeys()
	if len(keys) == 0 {
		return "50"
	}

	target := p * 100.0
	bestKey := keys[0].String()
	bestDiff := math.Abs(float64(atoiSafe(bestKey)) - target)

	for _, k := range keys[1:] {
		s := k.String()
		n := float64(atoiSafe(s))
		d := math.Abs(n - target)
		if d < bestDiff {
			bestDiff = d
			bestKey = s
		}
	}
	return bestKey
}

func atoiSafe(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		// try to strip non-digits
		digits := ""
		for _, r := range s {
			if r >= '0' && r <= '9' {
				digits += string(r)
			}
		}
		if digits == "" {
			return 50
		}
		i, err = strconv.Atoi(digits)
		if err != nil {
			return 50
		}
		return i
	}
	return i
}

// SampleRoundEndFromABM samples the round end reason using ABM distributions
// csfProb should be in [0,1]
func SampleRoundEndFromABM(ctWin bool, csfProb float64) string {
	if abmModels == nil || abmModels.RoundEndReasonDistributions == nil {
		// fallback
		if ctWin {
			if rand.Float64() < 0.7 {
				return "Elimination"
			}
			return "Defused"
		} else {
			if rand.Float64() < 0.68 {
				return "Elimination"
			}
			return "Exploded"
		}
	}

	side := "T"
	if ctWin {
		side = "CT"
	}
	sideMap, ok := abmModels.RoundEndReasonDistributions[side]
	if !ok || len(sideMap) == 0 {
		// fallback simple
		return SampleRoundEndFromABM(ctWin, csfProb)
	}

	key := csfPercentKeyForProb(sideMap, csfProb)
	bucket := sideMap[key]

	// Build slice of entries sorted by probability for deterministic accumulation
	type entry struct {
		Prob float64
		Name string
	}
	entries := make([]entry, 0, len(bucket.ReasonDistribution))
	for _, e := range bucket.ReasonDistribution {
		entries = append(entries, entry{Prob: e.Probability, Name: e.ReasonName})
	}

	// accumulate
	r := rand.Float64()
	cum := 0.0
	for _, e := range entries {
		cum += e.Prob
		if r <= cum {
			return e.Name
		}
	}
	if len(entries) > 0 {
		return entries[0].Name
	}
	return "Elimination"
}

// SampleSurvivorsFromABM samples number of survivors for the winning side based on ABM survivor distributions
// side should be "CT" or "T"; reasonCode is the code string used in the abm_models (e.g., "8" for elimination)
func SampleSurvivorsFromABM(side string, reasonCode string, csfProb float64) int {
	if abmModels == nil || abmModels.SurvivorDistributions == nil {
		// fallback assume typical survivors: 3 for elimination, 0 for exploded/defused
		if reasonCode == "8" {
			return 3
		}
		if reasonCode == "7" {
			return 0
		}
		if reasonCode == "1" {
			return 0
		}
		return 2
	}

	sideMap, ok := abmModels.SurvivorDistributions[side]
	if !ok {
		return 0
	}
	reasonMap, ok := sideMap[reasonCode]
	if !ok {
		return 0
	}
	key := csfPercentKeyForProb(reasonMap, csfProb)
	bucket := reasonMap[key]

	// Build ordered list from cumulative distribution
	type sdEntry struct {
		Survivors int
		Prob      float64
	}
	entries := make([]sdEntry, 0, len(bucket.CumulativeDistribution))
	for sStr, ent := range bucket.CumulativeDistribution {
		sInt := atoiSafe(sStr)
		entries = append(entries, sdEntry{Survivors: sInt, Prob: ent.Probability})
	}

	r := rand.Float64()
	cum := 0.0
	for _, e := range entries {
		cum += e.Prob
		if r <= cum {
			return e.Survivors
		}
	}
	if len(entries) > 0 {
		return entries[0].Survivors
	}
	return 0
}

// SampleEquipmentSavedFromABM provides a rough sample of equipment saved (mean) given side, reason and survivors
func SampleEquipmentSavedFromABM(side string, reasonCode string, survivors int) float64 {
	if abmModels == nil || abmModels.EquipmentSavedDistributions == nil {
		// fallback simple heuristic
		if survivors == 0 {
			return 0
		}
		return float64(survivors) * 100.0
	}

	sideMap, ok := abmModels.EquipmentSavedDistributions[side]
	if !ok {
		return 0
	}
	reasonMap, ok := sideMap[reasonCode]
	if !ok {
		return 0
	}
	survKey := fmt.Sprintf("%d", survivors)
	bucket, ok := reasonMap[survKey]
	if !ok {
		// fallback to any available
		for _, b := range reasonMap {
			return b.Mean
		}
		return 0
	}
	return bucket.Mean
}
