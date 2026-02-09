package generator

import (
	"fmt"
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
	if err := g.CreateFileFromTemplate(servicePath, templates.CustomServiceTemplate, map[string]interface{}{
		"Config":      projectConfig,
		"ServiceName": serviceName,
	}); err != nil {
		return err
	}

	// Automated Wiring
	mainPath := filepath.Join(baseDir, "main.go")
	serviceInit := fmt.Sprintf("%[1]sService := service.New%[2]sService()", strings.ToLower(serviceName), serviceName)
	return g.InjectCode(mainPath, "// [SERVICES-INIT]", serviceInit)
}
