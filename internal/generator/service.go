package generator

import (
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/templates"
)

// GenerateServiceFile generates a standalone service file
func (g *Generator) GenerateServiceFile(projectConfig config.ProjectConfig, serviceName string) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	servicePath := filepath.Join(baseDir, "service", strings.ToLower(serviceName)+".go")
	return g.CreateFileFromTemplate(servicePath, templates.CustomServiceTemplate, map[string]interface{}{
		"Config":      projectConfig,
		"ServiceName": serviceName,
	})
}
