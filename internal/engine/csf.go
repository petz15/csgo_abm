package engine

import (
	"math"
)

// ContestSuccessFunction calculates the probability of winning a round based on team expenditures and other factors.
func ContestSuccessFunction_simple(x int, y int, r float64) float64 { //theoretically this is currently a tullock contest

	probability := (math.Pow(float64(x), r) / (math.Pow(float64(x), r) + math.Pow(float64(y), r)))
	return probability
}

// SkewedGaussianProbability calculates probability between two values with a skewed normal distribution
// x, y: The two values to compare
// skew: Skewness factor - positive values skew right, negative values skew left, 0 for symmetric
// minOutput, maxOutput: Range for the output values (defaults to [0,1] if not specified)
func SkewedGaussianProbability_std4(x float64, y float64, skew float64, minOutput float64, maxOutput float64) float64 {
	// Calculate ratio of first value to total, normalized to [0,1]
	ratio := 0.0
	if x+y > 0 {
		ratio = x / (x + y)
	}

	// Normalize mean and stdDev to the output range
	rangeWidth := maxOutput - minOutput
	normalizedMean := maxOutput - (rangeWidth / 2)
	normalizedStdDev := rangeWidth / 4 // Adjust stdDev to fit within the output range

	// Apply the skew transformation
	transformedRatio := ratio
	if skew != 0 {
		delta := skew / math.Sqrt(1+skew*skew)
		transformedRatio = ratio - delta*math.Sin(2*math.Pi*ratio)
	}

	// Calculate probability using normal CDF approximation
	z := (transformedRatio - normalizedMean) / normalizedStdDev
	probability := 0.5 * (1 + math.Erf(z/math.Sqrt(2)))

	// Ensure the result is within [0,1]
	probability = math.Max(0, math.Min(1, probability))

	// Scale the result to the desired output range
	return minOutput + probability*rangeWidth
}
