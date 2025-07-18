package output

import (
	"encoding/json"
	"os"
)

func ExportResultsToJSON(results interface{}, path string) error {
	if path == "" {
		path = "results.json"
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
