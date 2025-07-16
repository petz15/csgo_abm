package util

import (
	"math/rand"
)

// ContestSuccessFunction calculates the probability of winning a round based on team expenditures and other factors.
func ContestSuccessFunction(ctSpend, tSpend int, ctSurvivors, tSurvivors int, bombPlanted bool) float64 {
	baseProbability := float64(ctSpend) / float64(ctSpend+tSpend)
	survivorFactor := float64(ctSurvivors) / float64(ctSurvivors+tSurvivors)
	bombFactor := 1.0
	if bombPlanted {
		bombFactor = 1.2 // Increase probability if the bomb is planted
	}

	return baseProbability * survivorFactor * bombFactor
}

// RandomOutcome simulates a round outcome based on the CSF and returns the winning team.
func RandomOutcome(ctSpend, tSpend int, ctSurvivors, tSurvivors int, bombPlanted bool) string {
	probability := ContestSuccessFunction(ctSpend, tSpend, ctSurvivors, tSurvivors, bombPlanted)
	if rand.Float64() < probability {
		return "CT" // Counter-Terrorists win
	}
	return "T" // Terrorists win
}