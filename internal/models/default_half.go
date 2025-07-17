package models

func InvestDecisionMaking_half(funds int, curround int, curscoreopo int) int {

	if curround == 15 || curround == 30 {
		return funds // Invest half of the funds at round 15 and 30
	} else {
		return funds / 2 // Invest half of the funds in other rounds
	}

}
