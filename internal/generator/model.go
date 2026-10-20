package generator

import (
	"fmt"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/naming"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateModelFiles generates all files for a model and wires them into
// main.go and routes.go. It is all-or-nothing: if any step fails, files it
// created are removed and files it modified are restored.
func (g *Generator) GenerateModelFiles(projectConfig config.ProjectConfig, model config.ModelConfig) error {
	if err := naming.Validate(model.Name); err != nil {
		return err
	}
	if model.HasService && !model.HasRepo {
		return fmt.Errorf("a service needs a repository: drop --no-repo or add --no-service")
	}
	if model.HasHandler && !model.HasService {
		return fmt.Errorf("a handler needs a service: drop --no-service or add --no-handler")
	}
	name := naming.Pascal(model.Name)
	model.Name = name
	model.ModuleName = projectConfig.ModuleName

	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}

	file := naming.Snake(name)
	paths := struct{ model, repo, service, handler, handlerTest, repoTest, serviceTest, main, routes, routesJSON string }{
		model:       filepath.Join(baseDir, "model", file+".go"),
		repo:        filepath.Join(baseDir, "repository", file+".go"),
		service:     filepath.Join(baseDir, "service", file+".go"),
		handler:     filepath.Join(baseDir, "transport/http/handler", file+"_handler.go"),
		handlerTest: filepath.Join(baseDir, "transport/http/handler", file+"_handler_test.go"),
		repoTest:    filepath.Join(baseDir, "repository", file+"_test.go"),
		serviceTest: filepath.Join(baseDir, "service", file+"_test.go"),
		main:        filepath.Join(baseDir, "main.go"),
		routes:      filepath.Join(baseDir, "transport/http/routes/routes.go"),
		routesJSON:  filepath.Join(baseDir, "transport/http/routes/routes.json"),
	}

	created := []string{paths.model}
	if model.HasRepo {
		created = append(created, paths.repo, paths.repoTest)
	}
	if model.HasService {
		created = append(created, paths.service, paths.serviceTest)
	}
	if model.HasHandler {
		created = append(created, paths.handler, paths.handlerTest)
	}

	if err := refuseOverwrite(g.Force, created...); err != nil {
		return err
	}
	if model.HasHandler && !g.Force {
		if err := g.checkRouterField(paths.routes, name+"Handler"); err != nil {
			return err
		}
	}

	return atomically([]string{paths.main, paths.routes, paths.routesJSON}, created, func() error {
		data := map[string]interface{}{"Config": projectConfig, "Model": model}

		if err := g.GenerateFromSource(paths.model, "core/model.go.tmpl", model); err != nil {
			return err
		}
		if model.HasRepo {
			if err := g.GenerateFromSource(paths.repo, "core/repository.go.tmpl", model); err != nil {
				return err
			}
			if err := g.GenerateFromSource(paths.repoTest, "core/repository_test.go.tmpl", data); err != nil {
				return err
			}
		}
		if model.HasService {
			if err := g.GenerateFromSource(paths.service, "core/service.go.tmpl", model); err != nil {
				return err
			}
			if err := g.GenerateFromSource(paths.serviceTest, "core/service_test.go.tmpl", data); err != nil {
				return err
			}
		}
		if model.HasHandler {
			if err := g.GenerateFromSource(paths.handler, "core/handler.go.tmpl", model); err != nil {
				return err
			}
			if err := g.GenerateFromSource(paths.handlerTest, "core/handler_test.go.tmpl", data); err != nil {
				return err
			}
		}
		return wireModel(projectConfig.ModuleName, name, model, paths.main, paths.routes, paths.routesJSON)
	})
}

// wireModel registers the generated repository, service and handler in main.go
// and the HTTP routes in routes.json.
func wireModel(moduleName, name string, model config.ModelConfig, mainPath, routesPath, routesJSONPath string) error {
	camel := naming.Camel(name)
	plural := naming.Plural(name)

	if model.HasRepo {
		if err := utils.AddImport(mainPath, moduleName+"/model"); err != nil {
			return err
		}
		if err := utils.InjectCodeAST(mainPath, "// [MIGRATE]", migrateStmt("model."+name, name)); err != nil {
			return err
		}
		if err := utils.AddImport(mainPath, moduleName+"/repository"); err != nil {
			return err
		}
		init := fmt.Sprintf("%sRepo := repository.New%sRepository(db)", camel, name)
		if err := utils.InjectCodeAST(mainPath, "// [REPOS-INIT]", init); err != nil {
			return err
		}
	}

	if model.HasService {
		if err := utils.AddImport(mainPath, moduleName+"/service"); err != nil {
			return err
		}
		init := fmt.Sprintf("%sService := service.New%sService(%sRepo)", camel, name, camel)
		if err := utils.InjectCodeAST(mainPath, "// [SERVICES-INIT]", init); err != nil {
			return err
		}
	}

	if !model.HasHandler {
		return nil
	}

	handlerImport := moduleName + "/transport/http/handler"
	if err := utils.AddImport(mainPath, handlerImport); err != nil {
		return err
	}
	init := fmt.Sprintf("%sHandler := handler.New%sHandler(%sService, validator)", camel, name, camel)
	if err := utils.InjectCodeAST(mainPath, "// [HANDLERS-INIT]", init); err != nil {
		return err
	}
	assign := fmt.Sprintf("router.%sHandler = %sHandler", name, camel)
	if err := utils.InjectCodeAST(mainPath, "// [ROUTER-HANDLERS-INIT]", assign); err != nil {
		return err
	}

	if err := utils.AddImport(routesPath, handlerImport); err != nil {
		return err
	}
	if err := utils.AddStructField(routesPath, "Router", name+"Handler", "*handler."+name+"Handler", ""); err != nil {
		return err
	}
	field := name + "Handler"
	base := "/" + naming.Kebab(plural)
	return utils.AddRoutes(routesJSONPath, []utils.Route{
		{Method: "POST", Path: base, Handler: field + ".Create" + name},
		{Method: "GET", Path: base, Handler: field + ".GetAll" + plural},
		{Method: "GET", Path: base + "/:id", Handler: field + ".Get" + name},
		{Method: "PUT", Path: base + "/:id", Handler: field + ".Update" + name},
		{Method: "DELETE", Path: base + "/:id", Handler: field + ".Delete" + name},
	})
}

// checkRouterField fails when routes.go already registers a handler under the
// given Router field, which would produce duplicate routes.
func (g *Generator) checkRouterField(routesPath, field string) error {
	exists, err := utils.StructHasField(routesPath, "Router", field)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("a handler named %s is already registered in %s (use --force to wire it again)", field, routesPath)
	}
	return nil
}

// migrateStmt returns the main.go statement that creates the table for a model
// type such as "model.Product" or "order.Order".
func migrateStmt(typ, label string) string {
	return fmt.Sprintf(`if err := db.AutoMigrate(&%s{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate %s")
	}`, typ, label)
}
