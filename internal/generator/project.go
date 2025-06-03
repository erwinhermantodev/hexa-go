package generator

import (
	"os"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/templates"
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
	files := map[string]string{
		"go.mod":                          templates.GoModTemplate,
		"README.md":                       templates.ReadmeTemplate,
		"Dockerfile":                      templates.DockerfileTemplate,
		"docker-compose.yml":              templates.DockerComposeTemplate,
		".gitignore":                      templates.GitignoreTemplate,
		".env.example":                    templates.EnvExampleTemplate,
		"Makefile":                        templates.MakefileTemplate,
		"locales/en.json":                 templates.LocaleEnTemplate,
		"locales/id.json":                 templates.LocaleIdTemplate,
		"repository/interfaces.go":        templates.RepositoryInterfacesTemplate,
		"transport/http/routes/routes.go": templates.HttpRoutesTemplate,
		"transport/grpc/server.go":        templates.GrpcServerTemplate,
		"transport/grpc/run.go":           templates.GrpcRunTemplate,
		"utils/codes.go":                  templates.UtilsCodesTemplate,
		"utils/config.go":                 templates.UtilsConfigTemplate,
		"utils/jwt.go":                    templates.UtilsJwtTemplate,
		"utils/messages.go":               templates.UtilsMessagesTemplate,
		"utils/password.go":               templates.UtilsPasswordTemplate,
		"utils/validator.go":              templates.UtilsValidatorTemplate,
		"main.go":                         templates.MainServerTemplate,
	}

	for filePath, tmplContent := range files {
		if err := g.CreateFileFromTemplate(filepath.Join(baseDir, filePath), tmplContent, projectConfig); err != nil {
			return err
		}
	}

	return nil
}
