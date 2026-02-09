package generator

import (
	"os"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
)

// CreateProject creates the entire project structure
func (g *Generator) CreateProject(projectConfig config.ProjectConfig) error {
	baseDir := projectConfig.Name

	// Create directory structure
	dirs := []string{
		"locales",
		"model",
		"repository",
		"service",
		"transport/grpc/proto",
		"transport/http/handler",
		"transport/http/routes",
		"utils",
		"docs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return err
		}
	}

	// Generate base files
	if err := g.generateBaseFiles(baseDir, projectConfig); err != nil {
		return err
	}

	// Generate model-specific files
	for _, model := range projectConfig.Models {
		if err := g.GenerateModelFiles(projectConfig, model); err != nil {
			return err
		}
	}

	// Generate custom services
	for _, service := range projectConfig.Services {
		if err := g.GenerateServiceFile(projectConfig, service); err != nil {
			return err
		}
	}

	return nil
}

// generateBaseFiles generates all base project files
func (g *Generator) generateBaseFiles(baseDir string, projectConfig config.ProjectConfig) error {
	// For now, mapping hardcoded to source paths.
	// In a full implementation, this would be driven entirely by blueprint.yaml
	files := map[string]string{
		"go.mod":                          "base/go.mod.tmpl",
		"main.go":                         "base/main.go.tmpl",
		"README.md":                       "base/README.md.tmpl",          // Need to extract these too
		"Dockerfile":                      "base/Dockerfile.tmpl",         // Need to extract these too
		"docker-compose.yml":              "base/docker-compose.yml.tmpl", // Need to extract these too
		"transport/http/routes/routes.go": "base/routes.go.tmpl",          // Need to extract these too
	}

	for filePath, sourcePath := range files {
		content, err := g.GetTemplateContent("", sourcePath)
		if err != nil {
			// fallback if not extracted yet
			continue
		}
		if err := g.CreateFileFromTemplate(filepath.Join(baseDir, filePath), content, projectConfig); err != nil {
			return err
		}
	}

	return nil
}
