package strategy

import "math"

func InvestDecisionMaking_anti_allin_v3(ctx StrategyContext_simple) float64 {
	// anti_allin_v3, invests all in the beginning and end of halves/overtime. In between it tries to build up wealth
	Score_to_Win := ctx.GameRules_strategy.HalfLength + (ctx.GameRules_strategy.OTHalfLength * (ctx.OvertimeAmount)) + 1

	pressing_ratio := 1.5 //how much more the team wants to invest in order to win
	overturn_ratio := 0.6 //how much less the team wants to invest in order to save for OT

	//always go all in if pistol round (in not OT rounds) or last round of half
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.IsLastRoundHalf && !ctx.IsOvertime && ctx.ConsecutiveLosses < 1 {
		//also go all in if round after pistol and it is not overtime and we won pistol
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 2 || Score_to_Win-ctx.OwnScore == 2 {
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
		return ctx.Funds
	}

	if !ctx.IsOvertime {
		if ctx.ConsecutiveLosses < 1 {
			//try to invest accoring to the pressing factor

			//technically the loss bonus calculation is not accurate.
			// Because it does not account for the loss bonus calculation method
			// Accordingly, the loss bonus might be higher with the new method (i.e. team wins inbetween, the loss bonus is not reset)
			all_in_loss_bonus := math.Min(float64(ctx.ConsecutiveWins), float64(len(ctx.GameRules_strategy.LossBonus)-1))
			approx_funds := (ctx.GameRules_strategy.LossBonus[int(all_in_loss_bonus)] + ctx.GameRules_strategy.DefaultEquipment) * 5
			approx_saved_eq := 0.0
			if ctx.ConsecutiveWins == 1 {
				approx_saved_eq += float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
			} else {
				approx_saved_eq += float64(ctx.EnemySurvivors) * ctx.GameRules_strategy.LossBonus[int(math.Max(1, all_in_loss_bonus-1))]
			}
			additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

			approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

			return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*pressing_ratio)
		}

		approx_funds := (ctx.GameRules_strategy.RoundOutcomeReward[ctx.RoundEndReason-1] + ctx.GameRules_strategy.DefaultEquipment) * 5
		approx_saved_eq := float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
		additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

		approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

		//if we have enough funds to challenge all in, do it
		if ctx.Funds+ctx.Equipment >= approx_allin_investment*overturn_ratio {
			return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*pressing_ratio)
		}
		//otherwise save until enough funds are built up
		return 0.0

	} else {

		//if somebody is about to win, go all in
		if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
			return ctx.Funds
		} else {
			//try to divide the funds evenly over the OT rounds
			return ctx.Funds / float64(ctx.GameRules_strategy.OTHalfLength)
		}
	}

	return 0.0
}
