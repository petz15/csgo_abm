package models

func InvestDecisionMaking_half(funds float64, curround int, curscoreopo int) float64 {

	if curround == 15 || curround == 30 {
		return funds // Invest half of the funds at round 15 and 30
	} else {
		return funds / 2 // Invest half of the funds in other rounds
	}

}
