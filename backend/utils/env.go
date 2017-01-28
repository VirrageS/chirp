package utils

import (
	"os"
)

// GetenvOrDefault retrieves the value of the environment variable named
// by the `key`. If the value is empty `fallback` is returned instead.
func GetenvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
