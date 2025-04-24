package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

// GetHostname returns the system hostname or "unknown" if it cannot be determined
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Failed to get hostname", "error", err)
		return "unknown"
	}
	return hostname
} 