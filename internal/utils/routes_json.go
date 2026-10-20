package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Route is one entry of a generated project's routes.json. Path is relative to
// the file's prefix; Handler is "RouterField.Method", for example
// "ProductHandler.CreateProduct".
type Route struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
}

// RouteFile is the content of routes.json.
type RouteFile struct {
	Prefix string  `json:"prefix"`
	Routes []Route `json:"routes"`
}

// AddRoutes appends routes to a routes.json file. A route with the same method
// and path as an existing one is skipped, so wiring twice (--force) does not
// register duplicates.
func AddRoutes(filePath string, routes []Route) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var rf RouteFile
	if err := json.Unmarshal(content, &rf); err != nil {
		return fmt.Errorf("%s is not valid JSON: %w", filePath, err)
	}

	seen := map[string]bool{}
	for _, r := range rf.Routes {
		seen[r.Method+" "+r.Path] = true
	}
	for _, r := range routes {
		if key := r.Method + " " + r.Path; !seen[key] {
			seen[key] = true
			rf.Routes = append(rf.Routes, r)
		}
	}

	out, err := json.MarshalIndent(rf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, append(out, '\n'), 0644)
}
