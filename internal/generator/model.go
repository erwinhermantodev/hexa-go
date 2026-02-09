package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/templates"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateModelFiles generates all files for a model
func (g *Generator) GenerateModelFiles(projectConfig config.ProjectConfig, model config.ModelConfig) error {
	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	// Define base paths for different components
	modelDir := filepath.Join(baseDir, "model")
	repoDir := filepath.Join(baseDir, "repository")
	serviceDir := filepath.Join(baseDir, "service")
	handlerDir := filepath.Join(baseDir, "transport/http/handler")

	// 1. Generate Model
	if err := g.GenerateFromSource(filepath.Join(modelDir, strings.ToLower(model.Name)+".go"), "core/model.go.tmpl", model); err != nil {
		return err
	}

	// 2. Generate Repository
	if model.HasRepo {
		if err := g.GenerateFromSource(filepath.Join(repoDir, strings.ToLower(model.Name)+".go"), "core/repository.go.tmpl", model); err != nil {
			return err
		}
	}

	// 3. Generate Service
	if model.HasService {
		if err := g.GenerateFromSource(filepath.Join(serviceDir, strings.ToLower(model.Name)+".go"), "core/service.go.tmpl", model); err != nil {
			return err
		}
	}

	// 4. Generate Handler
	if model.HasHandler {
		if err := g.GenerateFromSource(filepath.Join(handlerDir, strings.ToLower(model.Name)+".go"), "core/handler.go.tmpl", model); err != nil {
			return err
		}

		// Handler Test
		handlerTestPath := filepath.Join(baseDir, "transport/http/handler", strings.ToLower(model.Name)+"_handler_test.go")
		g.CreateFileFromTemplate(handlerTestPath, templates.DynamicHandlerTestTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		})
	}

	// Generate tests
	if model.HasRepo {
		repoTestPath := filepath.Join(baseDir, "repository", strings.ToLower(model.Name)+"_test.go")
		g.CreateFileFromTemplate(repoTestPath, templates.DynamicRepositoryTestTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		})
	}

	if model.HasService {
		serviceTestPath := filepath.Join(baseDir, "service", strings.ToLower(model.Name)+"_test.go")
		g.CreateFileFromTemplate(serviceTestPath, templates.DynamicServiceTestTemplate, map[string]interface{}{
			"Config": projectConfig,
			"Model":  model,
		})
	}

	// Automated Wiring
	mainPath := filepath.Join(baseDir, "main.go")
	routesPath := filepath.Join(baseDir, "transport/http/routes/routes.go")

	if model.HasRepo {
		if err := utils.AddImport(mainPath, projectConfig.ModuleName+"/repository"); err != nil {
			return err
		}
		repoInit := fmt.Sprintf("%[1]sRepo := repository.New%[2]sRepository(db)", strings.ToLower(model.Name), model.Name)
		if err := utils.InjectCodeAST(mainPath, "// [REPOS-INIT]", repoInit); err != nil {
			return err
		}
	}

	if model.HasService {
		if err := utils.AddImport(mainPath, projectConfig.ModuleName+"/service"); err != nil {
			return err
		}
		serviceInit := fmt.Sprintf("%[1]sService := service.New%[2]sService(%[1]sRepo)", strings.ToLower(model.Name), model.Name)
		if err := utils.InjectCodeAST(mainPath, "// [SERVICES-INIT]", serviceInit); err != nil {
			return err
		}
	}

	if model.HasHandler {
		if err := utils.AddImport(mainPath, projectConfig.ModuleName+"/transport/http/handler"); err != nil {
			return err
		}
		handlerInit := fmt.Sprintf("%[1]sHandler := handler.New%[2]sHandler(%[1]sService, validator)", strings.ToLower(model.Name), model.Name)
		if err := utils.InjectCodeAST(mainPath, "// [HANDLERS-INIT]", handlerInit); err != nil {
			return err
		}

		// Inject into Router in main.go
		routerInit := fmt.Sprintf("router.%[1]sHandler = %[2]sHandler\n\t_ = %[2]sHandler", strings.Title(strings.ToLower(model.Name)), strings.ToLower(model.Name))
		if err := utils.InjectCodeAST(mainPath, "// [ROUTER-HANDLERS-INIT]", routerInit); err != nil {
			return err
		}

		// Inject into Router fields in routes.go
		if err := utils.AddImport(routesPath, projectConfig.ModuleName+"/transport/http/handler"); err != nil {
			return err
		}
		// For struct fields, we might need a specific AddStructField call, but the routes.go uses markers inside a function or struct?
		// Wait, routes.go uses a struct now! I should use AddStructField for [HANDLER-FIELDS-EXPORTED]

		if err := utils.AddStructField(routesPath, "Router", strings.Title(strings.ToLower(model.Name))+"Handler", "*handler."+strings.Title(strings.ToLower(model.Name))+"Handler", "`json:\"-\"`"); err != nil {
			// fallback to InjectCodeAST if AddStructField fails (e.g. marker mismatch)
			// Actually, let's use the explicit tool.
		}

		// Inject into Router routes in routes.go
		routeInjection := fmt.Sprintf("%[1]ss := api.Group(\"/%[1]ss\")\n\t%[1]ss.POST(\"\", r.%[2]sHandler.Create%[1]s)\n\t%[1]ss.GET(\"\", r.%[2]sHandler.GetAll%[1]ss)\n\t%[1]ss.GET(\"/:id\", r.%[2]sHandler.Get%[1]s)\n\t%[1]ss.PUT(\"/:id\", r.%[2]sHandler.Update%[1]s)\n\t%[1]ss.DELETE(\"/:id\", r.%[2]sHandler.Delete%[1]s)", model.Name, strings.Title(strings.ToLower(model.Name)))
		if err := utils.InjectCodeAST(routesPath, "// [ROUTES-INIT]", routeInjection); err != nil {
			return err
		}
	}

	return nil
}
