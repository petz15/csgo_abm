package output

import (
	"encoding/json"
	"os"
	"runtime"
)

func ExportResultsToJSON(results interface{}, path string) error {
	if path == "" {
		path = "results.json"
	}

	// Marshal to JSON in memory
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(path, jsonData, 0644)

	// Clear the large JSON data from memory
	jsonData = nil

	// Hint to the garbage collector
	runtime.GC()

	return err
}
