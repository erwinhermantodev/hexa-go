package utils

import (
	"fmt"
	"os"
	"strings"
)

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetModuleName extracts the module name from go.mod in the current directory.
func GetModuleName() (string, error) {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("no go.mod found; run this command in the root of a Go project")
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "module")), `"`), nil
		}
	}

	return "", fmt.Errorf("go.mod has no module directive")
}
