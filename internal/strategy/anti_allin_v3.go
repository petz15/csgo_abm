package strategy

import "math"

func InvestDecisionMaking_anti_allin_v3(ctx StrategyContext_simple) float64 {
	// anti_allin_v3, invests all in the beginning and end of halves/overtime. In between it tries to build up wealth
	Score_to_Win := ctx.GameRules_strategy.HalfLength + (ctx.GameRules_strategy.OTHalfLength * (ctx.OvertimeAmount)) + 1

	//always go all in if pistol round (in not OT rounds) or last round of half
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.IsLastRoundHalf && !ctx.IsOvertime && ctx.ConsecutiveLosses < 1 {
		//also go all in if round after pistol and it is not overtime and we won pistol
		return ctx.Funds
	}

	if !ctx.IsOvertime {
		if ctx.ConsecutiveLosses < 1 {
			//try to invest twice as much as all in strategy
			all_in_loss_bonus := math.Min(float64(ctx.ConsecutiveWins), float64(len(ctx.GameRules_strategy.LossBonus)))
			approx_funds := (ctx.GameRules_strategy.LossBonus[int(all_in_loss_bonus)] + ctx.GameRules_strategy.DefaultEquipment) * 5
			approx_saved_eq := 0.0
			if ctx.ConsecutiveWins == 1 {
				approx_saved_eq += float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
			} else {
				approx_saved_eq += float64(ctx.EnemySurvivors) * ctx.GameRules_strategy.LossBonus[int(math.Max(1, all_in_loss_bonus-1))]
			}
			additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

			approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

			return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*2)
		}

		approx_funds := (ctx.GameRules_strategy.RoundOutcomeReward[ctx.RoundEndReason-1] + ctx.GameRules_strategy.DefaultEquipment) * 5
		approx_saved_eq := float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
		additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

		approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

		//if we have enough funds to challenge all in, do it
		if ctx.Funds+ctx.Equipment >= approx_allin_investment {
			return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*2)

			//if scores are close, go all in
		} else if Score_to_Win-ctx.OpponentScore == 2 || Score_to_Win-ctx.OwnScore == 2 {
			return ctx.Funds
		} else if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
			return ctx.Funds
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
