//change the file to a .go extension and run it. shows how distributions of different probabilities change.
// helpful to visualize probabilities.go 

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// Copy of the CSF function for analysis
func contestSuccessFunction(x float64, y float64, r float64) float64 {
	return (math.Pow(x, r) / (math.Pow(x, r) + math.Pow(y, r)))
}

// Copy of the distribution function for analysis
func csfNormalDistribution(x float64, y float64, r float64, minOutput float64, maxOutput float64, stdDevFactor float64, skew float64) float64 {
	rangeWidth := maxOutput - minOutput

	csfProb := 0.5
	if x+y > 0 {
		csfProb = contestSuccessFunction(x, y, r)
	}

	mean := minOutput + csfProb*rangeWidth
	stdDev := rangeWidth / stdDevFactor

	u := rand.Float64()

	// Apply skewness using the Azzalini's skew-normal transformation
	alpha := skew * 5 // scale skew to make it more pronounced

	// Generate two independent standard normal variables
	z0 := math.Sqrt(2) * math.Erfinv(2*u-1)
	z1 := math.Sqrt(2) * math.Erfinv(2*rand.Float64()-1)

	// Skew-normal variable
	skewedZ := alpha*math.Abs(z0) + z1

	result := mean + skewedZ*stdDev

	return math.Max(minOutput, math.Min(maxOutput, result))
}

func main() {
	const (
		minOutput = 0.0
		maxOutput = 5.0
		samples   = 10000
	)

	fmt.Printf("Distribution Analysis: minOutput=%.1f, maxOutput=%.1f\n", minOutput, maxOutput)
	fmt.Printf("Samples per test: %d\n\n", samples)

	// Test scenarios
	scenarios := []struct {
		name         string
		x, y, r      float64
		stdDevFactor float64
		skew         float64
	}{
		// Basic scenarios with equal teams
		{"Equal teams, normal dist, std/4", 100, 100, 1.0, 4.0, 0.0},
		{"Equal teams, normal dist, std/2", 100, 100, 1.0, 2.0, 0.0},
		{"Equal teams, normal dist, std/8", 100, 100, 1.0, 8.0, 0.0},

		// Different CSF parameters (r values)
		{"Equal teams, r=0.5 (less deterministic)", 100, 100, 0.5, 4.0, 0.0},
		{"Equal teams, r=2.0 (more deterministic)", 100, 100, 2.0, 4.0, 0.0},

		// Unequal teams
		{"Team1 stronger (2:1)", 200, 100, 1.0, 4.0, 0.0},
		{"Team1 much stronger (3:1)", 300, 100, 1.0, 4.0, 0.0},
		{"Team1 weaker (1:2)", 100, 200, 1.0, 4.0, 0.0},

		// Skewness effects
		{"Equal teams, left skew", 100, 100, 1.0, 4.0, -0.5},
		{"Equal teams, right skew", 100, 100, 1.0, 4.0, 0.5},
		{"Equal teams, strong left skew", 100, 100, 1.0, 4.0, -1.0},
		{"Equal teams, strong right skew", 100, 100, 1.0, 4.0, 1.0},

		// Combined effects
		{"Strong team1 + right skew", 200, 100, 1.0, 4.0, 0.5},
		{"Weak team1 + left skew", 100, 200, 1.0, 4.0, -0.5},
	}

	for _, scenario := range scenarios {
		fmt.Printf("=== %s ===\n", scenario.name)
		fmt.Printf("Parameters: x=%.1f, y=%.1f, r=%.1f, stdDevFactor=%.1f, skew=%.1f\n",
			scenario.x, scenario.y, scenario.r, scenario.stdDevFactor, scenario.skew)

		// Calculate theoretical CSF probability
		csfProb := contestSuccessFunction(scenario.x, scenario.y, scenario.r)
		theoreticalMean := minOutput + csfProb*(maxOutput-minOutput)
		theoreticalStdDev := (maxOutput - minOutput) / scenario.stdDevFactor

		fmt.Printf("Theoretical: CSF prob=%.3f, mean=%.3f, stdDev=%.3f\n",
			csfProb, theoreticalMean, theoreticalStdDev)

		// Generate samples
		samples_slice := make([]float64, samples)
		for i := 0; i < samples; i++ {
			samples_slice[i] = csfNormalDistribution(
				scenario.x, scenario.y, scenario.r, minOutput, maxOutput,
				scenario.stdDevFactor, scenario.skew)
		}

		// Calculate statistics
		sort.Float64s(samples_slice)
		mean := calculateMean(samples_slice)
		stdDev := calculateStdDev(samples_slice, mean)
		median := samples_slice[len(samples_slice)/2]
		p25 := samples_slice[len(samples_slice)/4]
		p75 := samples_slice[3*len(samples_slice)/4]

		fmt.Printf("Observed:    mean=%.3f, stdDev=%.3f, median=%.3f\n", mean, stdDev, median)
		fmt.Printf("Percentiles: 25th=%.3f, 75th=%.3f\n", p25, p75)
		fmt.Printf("Range:       min=%.3f, max=%.3f\n", samples_slice[0], samples_slice[len(samples_slice)-1])

		// Show distribution bins
		showDistribution(samples_slice, minOutput, maxOutput, 10)
		fmt.Println()
	}
}

func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func calculateStdDev(data []float64, mean float64) float64 {
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(data)-1))
}

func showDistribution(data []float64, min, max float64, bins int) {
	binWidth := (max - min) / float64(bins)
	binCounts := make([]int, bins)

	for _, v := range data {
		binIndex := int((v - min) / binWidth)
		if binIndex >= bins {
			binIndex = bins - 1
		}
		if binIndex < 0 {
			binIndex = 0
		}
		binCounts[binIndex]++
	}

	fmt.Printf("Distribution (bins of %.2f):\n", binWidth)
	maxCount := 0
	for _, count := range binCounts {
		if count > maxCount {
			maxCount = count
		}
	}

	for i, count := range binCounts {
		binStart := min + float64(i)*binWidth
		binEnd := binStart + binWidth
		percentage := float64(count) / float64(len(data)) * 100

		// Create visual bar
		barLength := int(float64(count) / float64(maxCount) * 30)
		bar := ""
		for j := 0; j < barLength; j++ {
			bar += "█"
		}

		fmt.Printf("  [%.2f-%.2f): %5d (%4.1f%%) %s\n",
			binStart, binEnd, count, percentage, bar)
	}
}
