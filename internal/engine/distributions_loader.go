package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// Sorted slice types for efficient ordered access
type ThresholdReason struct {
	Threshold float64
	Reason    int
}

type ThresholdInt struct {
	Threshold float64
	Value     int
}

type PercentileValue struct {
	Percentile float64
	Value      float64
}

// ProcessedDistributions contains distributions with sorted slices for efficient sampling
type ProcessedDistributions struct {
	Metadata struct {
		CSFRValue float64
		CSFRanges struct {
			Min int
			Max int
		}
	}

	// RoundEndReason: side -> csfKey -> sorted thresholds
	RoundEndReason map[string]map[string][]ThresholdReason

	// BombPlanted: simple probability map (no sorting needed)
	BombPlantedT map[string]float64

	// Survivors: side -> reasonCode -> csfKey -> sorted thresholds
	Survivors map[string]map[string]map[string][]ThresholdInt

	// EquipmentSaved: side -> reasonCode -> survivorCount -> sorted percentiles
	EquipmentSaved map[string]map[string]map[string][]PercentileValue
}

// Raw JSON structures for unmarshaling
type rawDistributions struct {
	Metadata struct {
		CSFRValue float64 `json:"csf_r_value"`
		CSFRanges struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"csf_ranges"`
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

// loadAndProcessDistributions loads and converts distributions to sorted slices
func loadAndProcessDistributions(filePath string) (*ProcessedDistributions, error) {
	if filePath == "" {
		filePath = "distributions.json"
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if filePath == "distributions.json" {
			parentPath := filepath.Join("..", "distributions.json")
			data, err = os.ReadFile(parentPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read distributions file from '%s' or '%s': %w", filePath, parentPath, err)
			}
		} else {
			return nil, fmt.Errorf("failed to read distributions file '%s': %w", filePath, err)
		}
	}

	// Check if file is empty
	if len(data) == 0 {
		return nil, fmt.Errorf("distributions file is empty")
	}

	var raw rawDistributions
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal distributions file: %w", err)
	}

	// Check if distributions are present
	if len(raw.Distributions.RoundEndReason) == 0 {
		return nil, fmt.Errorf("distributions file does not contain round end reason data")
	}
	if len(raw.Distributions.Survivors) == 0 {
		return nil, fmt.Errorf("distributions file does not contain survivors data")
	}
	if len(raw.Distributions.EquipmentSaved) == 0 {
		return nil, fmt.Errorf("distributions file does not contain equipment saved data")
	}
	if len(raw.Distributions.BombPlanted.T) == 0 {
		return nil, fmt.Errorf("distributions file does not contain bomb planted data")
	}
	if raw.Metadata.CSFRValue < 0 {
		return nil, fmt.Errorf("distributions file does not contain valid CSF R value")
	}

	processed := &ProcessedDistributions{
		RoundEndReason: make(map[string]map[string][]ThresholdReason),
		BombPlantedT:   raw.Distributions.BombPlanted.T,
		Survivors:      make(map[string]map[string]map[string][]ThresholdInt),
		EquipmentSaved: make(map[string]map[string]map[string][]PercentileValue),
	}

	// Copy metadata fields manually
	processed.Metadata.CSFRValue = raw.Metadata.CSFRValue
	processed.Metadata.CSFRanges.Min = raw.Metadata.CSFRanges.Min
	processed.Metadata.CSFRanges.Max = raw.Metadata.CSFRanges.Max

	// Convert RoundEndReason
	for side, sideData := range raw.Distributions.RoundEndReason {
		processed.RoundEndReason[side] = make(map[string][]ThresholdReason)
		for csfKey, bucket := range sideData {
			// Skip empty distributions
			if len(bucket.CumulativeDistribution) == 0 {
				continue
			}

			var thresholds []ThresholdReason
			for thresholdStr, entry := range bucket.CumulativeDistribution {
				threshold, err := strconv.ParseFloat(thresholdStr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid threshold '%s' in RoundEndReason[%s][%s]: %w", thresholdStr, side, csfKey, err)
				}
				thresholds = append(thresholds, ThresholdReason{
					Threshold: threshold,
					Reason:    entry.Reason,
				})
			}
			// Sort by threshold (ascending)
			sort.Slice(thresholds, func(i, j int) bool {
				return thresholds[i].Threshold < thresholds[j].Threshold
			})
			processed.RoundEndReason[side][csfKey] = thresholds
		}
	}

	// Convert Survivors
	for side, sideData := range raw.Distributions.Survivors {
		processed.Survivors[side] = make(map[string]map[string][]ThresholdInt)
		for reasonCode, reasonData := range sideData {
			processed.Survivors[side][reasonCode] = make(map[string][]ThresholdInt)
			for csfKey, bucket := range reasonData.CSFDistributions {
				// Skip empty distributions
				if len(bucket.CumulativeLookup) == 0 {
					continue
				}

				var thresholds []ThresholdInt
				for thresholdStr, value := range bucket.CumulativeLookup {
					threshold, err := strconv.ParseFloat(thresholdStr, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid threshold '%s' in Survivors[%s][%s][%s]: %w", thresholdStr, side, reasonCode, csfKey, err)
					}
					thresholds = append(thresholds, ThresholdInt{
						Threshold: threshold,
						Value:     value,
					})
				}
				// Sort by threshold (ascending)
				sort.Slice(thresholds, func(i, j int) bool {
					return thresholds[i].Threshold < thresholds[j].Threshold
				})
				processed.Survivors[side][reasonCode][csfKey] = thresholds
			}
		}
	}

	// Convert EquipmentSaved
	for side, sideData := range raw.Distributions.EquipmentSaved {
		processed.EquipmentSaved[side] = make(map[string]map[string][]PercentileValue)
		for reasonCode, reasonData := range sideData {
			processed.EquipmentSaved[side][reasonCode] = make(map[string][]PercentileValue)
			for survivorKey, survivorDist := range reasonData.SurvivorDistributions {
				// Skip empty distributions
				if len(survivorDist.ECDFLookup) == 0 {
					continue
				}

				var percentiles []PercentileValue
				for percentileStr, value := range survivorDist.ECDFLookup {
					percentile, err := strconv.ParseFloat(percentileStr, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid percentile '%s' in EquipmentSaved[%s][%s][%s]: %w", percentileStr, side, reasonCode, survivorKey, err)
					}
					percentiles = append(percentiles, PercentileValue{
						Percentile: percentile,
						Value:      value,
					})
				}
				// Sort by percentile (ascending)
				sort.Slice(percentiles, func(i, j int) bool {
					return percentiles[i].Percentile < percentiles[j].Percentile
				})
				processed.EquipmentSaved[side][reasonCode][survivorKey] = percentiles
			}
		}
	}

	return processed, nil
}
