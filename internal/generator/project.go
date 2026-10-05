package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
)

// projectDirs are created up front so that empty directories exist even when
// no file is generated into them.
var projectDirs = []string{
	"configs",
	"locales",
	"model",
	"repository",
	"service",
	"transport/grpc",
	"transport/http/handler",
	"transport/http/routes",
	"utils",
	"docs",
	"migrations",
}

// baseFiles maps each generated file to its embedded template. Every entry is
// required: a missing template is a bug in the CLI, not something to skip.
var baseFiles = map[string]string{
	"go.mod":                          "base/go.mod.tmpl",
	"main.go":                         "base/main.go.tmpl",
	"README.md":                       "base/README.md.tmpl",
	"Dockerfile":                      "base/Dockerfile.tmpl",
	"docker-compose.yml":              "base/docker-compose.yml.tmpl",
	"Makefile":                        "base/Makefile.tmpl",
	".gitignore":                      "base/gitignore.tmpl",
	".env.example":                    "base/env.example.tmpl",
	".github/workflows/ci.yml":        "base/ci.yml.tmpl",
	"configs/config.yaml":             "base/config.yaml.tmpl",
	"locales/en.json":                 "base/locales/en.json.tmpl",
	"locales/id.json":                 "base/locales/id.json.tmpl",
	"docs/docs.go":                    "base/docs.go.tmpl",
	"transport/http/routes/routes.go": "base/routes.go.tmpl",
	"transport/grpc/server.go":        "base/grpc/server.go.tmpl",
	"transport/grpc/run.go":           "base/grpc/run.go.tmpl",
	"utils/codes.go":                  "base/utils/codes.go.tmpl",
	"utils/config.go":                 "base/utils/config.go.tmpl",
	"utils/database.go":               "base/utils/database.go.tmpl",
	"utils/jwt.go":                    "base/utils/jwt.go.tmpl",
	"utils/messages.go":               "base/utils/messages.go.tmpl",
	"utils/password.go":               "base/utils/password.go.tmpl",
	"utils/response.go":               "base/utils/response.go.tmpl",
	"utils/validator.go":              "base/utils/validator.go.tmpl",
}

// CreateProject creates the entire project structure. If generation fails, the
// project directory is removed so no half-built project is left behind.
func (g *Generator) CreateProject(projectConfig config.ProjectConfig) (err error) {
	baseDir := projectConfig.Name
	if baseDir == "" {
		return fmt.Errorf("project name is required")
	}
	if projectConfig.ModuleName == "" {
		return fmt.Errorf("module name is required")
	}

	if entries, readErr := os.ReadDir(baseDir); readErr == nil && len(entries) > 0 && !g.Force {
		return fmt.Errorf("directory %s already exists and is not empty (use --force to generate into it)", baseDir)
	}
	_, statErr := os.Stat(baseDir)
	createdDir := os.IsNotExist(statErr)

	defer func() {
		if err != nil && createdDir {
			os.RemoveAll(baseDir)
		}
	}()

	for _, dir := range projectDirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return err
		}
	}

	if err := g.generateBaseFiles(baseDir, projectConfig); err != nil {
		return err
	}

	for _, model := range projectConfig.Models {
		if err := g.GenerateModelFiles(projectConfig, model); err != nil {
			return fmt.Errorf("model %s: %w", model.Name, err)
		}
	}

	for _, service := range projectConfig.Services {
		if err := g.GenerateServiceFile(projectConfig, service); err != nil {
			return fmt.Errorf("service %s: %w", service, err)
		}
	}

	return nil
}

// generateBaseFiles renders every entry of baseFiles into baseDir.
func (g *Generator) generateBaseFiles(baseDir string, projectConfig config.ProjectConfig) error {
	for filePath, sourcePath := range baseFiles {
		if err := g.GenerateFromSource(filepath.Join(baseDir, filePath), sourcePath, projectConfig); err != nil {
			return fmt.Errorf("%s: %w", filePath, err)
		}
	}
	return nil
}
