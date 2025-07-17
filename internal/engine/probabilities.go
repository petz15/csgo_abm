package engine

import (
	"math"
	"math/rand"
)

// ContestSuccessFunction calculates the probability of winning a round based on team expenditures and other factors.
// This is a Tullock contest success function
func ContestSuccessFunction_simple(x float64, y float64, r float64) float64 {
	probability := (math.Pow(x, r) / (math.Pow(x, r) + math.Pow(y, r)))
	return probability
}

func bool_CSF_simple(x float64, y float64, r float64) bool {
	probability := ContestSuccessFunction_simple(x, y, r)
	return rand.Float64() < probability
}

// CSFNormalDistribution generates a value from a normal distribution with mean determined by Contest Success Function
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
		csfProb = ContestSuccessFunction_simple(x, y, r)
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

// CSFNormalDistribution generates a value from a normal distribution with mean determined by Contest Success Function
// x, y: The two values to compare (e.g., team expenditures)
// r: The parameter for the CSF (higher r = more deterministic outcomes)
// minOutput, maxOutput: Range for the output values
// Returns a value sampled from a normal distribution between minOutput and maxOutput
func CSFNormalDistribution_std_custom(x float64, y float64, r float64, minOutput float64, maxOutput float64, stdDevFactor float64) float64 {
	// Calculate the range width for scaling
	rangeWidth := maxOutput - minOutput

	// First calculate the base probability using Contest Success Function (Tullock contest)
	csfProb := 0.5
	if x+y > 0 {
		csfProb = math.Pow(x, r) / (math.Pow(x, r) + math.Pow(y, r))
	}

	// Calculate the mean position between minOutput and maxOutput based on CSF probability
	mean := minOutput + csfProb*rangeWidth

	// Set stdDev to 1/4 of the range width as specified
	stdDev := rangeWidth / stdDevFactor

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
