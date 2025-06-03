package generator

import (
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/go-starter-cli/internal/config"
	"github.com/erwinhermantodev/go-starter-cli/internal/templates"
)

// GenerateModelFiles generates all files for a model
func (g *Generator) GenerateModelFiles(projectConfig config.ProjectConfig, model config.ModelConfig) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	// Generate model file
	modelPath := filepath.Join(baseDir, "model", strings.ToLower(model.Name)+".go")
	if err := g.CreateFileFromTemplate(modelPath, templates.DynamicModelTemplate, map[string]interface{}{
		"Config": projectConfig,
		"Model":  model,
	}); err != nil {
		return err
	}

	// Generate repository if needed
	if model.HasRepo {
		repoPath := filepath.Join(baseDir, "repository", strings.ToLower(model.Name)+".go")
		if err := g.CreateFileFromTemplate(repoPath, templates.DynamicRepositoryTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		}); err != nil {
			return err
		}
	}

	// Generate service if needed
	if model.HasService {
		servicePath := filepath.Join(baseDir, "service", strings.ToLower(model.Name)+".go")
		if err := g.CreateFileFromTemplate(servicePath, templates.DynamicServiceTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		}); err != nil {
			return err
		}
	}

	// Generate handler if needed
	if model.HasHandler {
		handlerPath := filepath.Join(baseDir, "transport/http/handler", strings.ToLower(model.Name)+"_handler.go")
		if err := g.CreateFileFromTemplate(handlerPath, templates.DynamicHandlerTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		}); err != nil {
			return err
		}
	}

	return nil
}
