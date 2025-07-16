package output

import (
	"encoding/json"
	"os"
)

type SimulationResult struct {
	TeamAName       string  `json:"team_a_name"`
	TeamBName       string  `json:"team_b_name"`
	TeamAScore      int     `json:"team_a_score"`
	TeamBScore      int     `json:"team_b_score"`
	TeamAFinancials Financials `json:"team_a_financials"`
	TeamBFinancials Financials `json:"team_b_financials"`
}

type Financials struct {
	StartMoney int `json:"start_money"`
	EndMoney   int `json:"end_money"`
}

func ExportResults(results []SimulationResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(results)
}