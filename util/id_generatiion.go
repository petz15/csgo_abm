package util

import (
	"crypto/rand"
	"io"
	"os"
	"time"
)

// CreateGameID creates a unique identifier for a game session
func CreateGameID() string {
	// Get the hostname
	hostname := "host"
	name, err := os.Hostname()
	if err == nil {
		hostname = name
	}

	// Generate random component
	uniqueID := generateRandomString(4)

	// Format: YYYYMMDD_HHMMSS_hostname_randomID
	ID := time.Now().Format("20060102_150405") + "_" + hostname + "_" + uniqueID

	return ID
}

// generateRandomString generates a random string of specified length
func generateRandomString(max int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
