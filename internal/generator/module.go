package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateModuleFiles generates all files for a module in a single directory
func (g *Generator) GenerateModuleFiles(projectConfig config.ProjectConfig, model config.ModelConfig, moduleName string) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	modulePath := filepath.Join(baseDir, "internal/modules", strings.ToLower(moduleName))
	packageName := strings.ToLower(moduleName)

	data := map[string]interface{}{
		"Config":      projectConfig,
		"Model":       model,
		"PackageName": packageName,
	}

	// 1. Generate Model
	if err := g.GenerateFromSource(filepath.Join(modulePath, "model.go"), "modules/model.go.tmpl", data); err != nil {
		return err
	}

	// 2. Generate Repository
	if err := g.GenerateFromSource(filepath.Join(modulePath, "repository.go"), "modules/repository.go.tmpl", data); err != nil {
		return err
	}

	// 3. Generate Service
	if err := g.GenerateFromSource(filepath.Join(modulePath, "service.go"), "modules/service.go.tmpl", data); err != nil {
		return err
	}

	// 4. Generate Handler
	if err := g.GenerateFromSource(filepath.Join(modulePath, "handler.go"), "modules/handler.go.tmpl", data); err != nil {
		return err
	}

	// 5. Automated Wiring
	mainPath := filepath.Join(baseDir, "main.go")
	routesPath := filepath.Join(baseDir, "transport/http/routes/routes.go")
	moduleImport := fmt.Sprintf("%s/internal/modules/%s", projectConfig.ModuleName, packageName)

	// Inject Imports
	if err := utils.AddImport(mainPath, moduleImport); err != nil {
		return err
	}
	if err := utils.AddImport(routesPath, moduleImport); err != nil {
		return err
	}

	// Inject Repository Init
	repoInit := fmt.Sprintf("%[1]sRepo := %[1]s.NewRepository(db)", packageName)
	if err := utils.InjectCodeAST(mainPath, "// [REPOS-INIT]", repoInit); err != nil {
		return err
	}

	// Inject Service Init
	serviceInit := fmt.Sprintf("%[1]sService := %[1]s.NewService(%[1]sRepo)", packageName)
	if err := utils.InjectCodeAST(mainPath, "// [SERVICES-INIT]", serviceInit); err != nil {
		return err
	}

	// Inject Handler Init
	handlerInit := fmt.Sprintf("%[1]sHandler := %[1]s.NewHandler(%[1]sService, validator)", packageName)
	if err := utils.InjectCodeAST(mainPath, "// [HANDLERS-INIT]", handlerInit); err != nil {
		return err
	}

	// Inject into Router in main.go
	routerInit := fmt.Sprintf("router.%sHandler = %sHandler", strings.Title(packageName), packageName)
	if err := utils.InjectCodeAST(mainPath, "// [ROUTER-HANDLERS-INIT]", routerInit); err != nil {
		return err
	}

	// Inject into Router fields in routes.go
	handlerFieldName := fmt.Sprintf("%sHandler", strings.Title(strings.ToLower(model.Name)))
	handlerFieldType := fmt.Sprintf("*%s.Handler", packageName)
	if err := utils.AddStructField(routesPath, "Router", handlerFieldName, handlerFieldType, "`json:\"-\"`"); err != nil {
		return err
	}

	// Inject into Router routes in routes.go
	routeInjection := fmt.Sprintf("%[1]ss := api.Group(\"/%[1]ss\")\n\t%[1]ss.POST(\"\", r.%[2]sHandler.Create%[2]s)\n\t%[1]ss.GET(\"\", r.%[2]sHandler.GetAll%[2]ss)\n\t%[1]ss.GET(\"/:id\", r.%[2]sHandler.Get%[2]s)\n\t%[1]ss.PUT(\"/:id\", r.%[2]sHandler.Update%[2]s)\n\t%[1]ss.DELETE(\"/:id\", r.%[2]sHandler.Delete%[2]s)", strings.ToLower(model.Name), strings.Title(packageName))
	if err := utils.InjectCodeAST(routesPath, "// [ROUTES-INIT]", routeInjection); err != nil {
		return err
	}

	return nil
}
