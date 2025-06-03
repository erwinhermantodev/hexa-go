package generator

import (
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/go-starter-cli/internal/config"
	"github.com/erwinhermantodev/go-starter-cli/internal/templates"
)

// GenerateHandlerFile generates a standalone handler file
func (g *Generator) GenerateHandlerFile(projectConfig config.ProjectConfig, handlerName string) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	handlerPath := filepath.Join(baseDir, "transport/http/handler", strings.ToLower(handlerName)+"_handler.go")
	return g.CreateFileFromTemplate(handlerPath, templates.CustomHandlerTemplate, map[string]interface{}{
		"Config":      projectConfig,
		"HandlerName": handlerName,
	})
}
