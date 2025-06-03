package utils

import (
	"os"
	"strings"
)

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetModuleName extracts module name from go.mod
func GetModuleName() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return "unknown"
}
