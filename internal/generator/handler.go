package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/templates"
)

// GenerateHandlerFile generates a standalone handler file
func (g *Generator) GenerateHandlerFile(projectConfig config.ProjectConfig, handlerName string) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	handlerPath := filepath.Join(baseDir, "transport/http/handler", strings.ToLower(handlerName)+"_handler.go")
	if err := g.CreateFileFromTemplate(handlerPath, templates.CustomHandlerTemplate, map[string]interface{}{
		"Config":      projectConfig,
		"HandlerName": handlerName,
	}); err != nil {
		return err
	}

	// Automated Wiring
	mainPath := filepath.Join(baseDir, "main.go")
	handlerInit := fmt.Sprintf("%[1]sHandler := handler.New%[2]sHandler(%[1]sService, validator)", strings.ToLower(handlerName), handlerName)
	return g.InjectCode(mainPath, "// [HANDLERS-INIT]", handlerInit)
}
