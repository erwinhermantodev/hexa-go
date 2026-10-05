package generator

import (
	"fmt"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/naming"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateModuleFiles generates all files for a module in a single directory
// and wires it into main.go and routes.go. Like GenerateModelFiles it is
// all-or-nothing.
func (g *Generator) GenerateModuleFiles(projectConfig config.ProjectConfig, model config.ModelConfig, moduleName string) error {
	if err := naming.ValidatePackage(moduleName); err != nil {
		return err
	}
	name := naming.Pascal(moduleName)
	model.Name = name
	model.ModuleName = projectConfig.ModuleName
	packageName := naming.Package(moduleName)

	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}
	modulePath := filepath.Join(baseDir, "internal/modules", packageName)
	mainPath := filepath.Join(baseDir, "main.go")
	routesPath := filepath.Join(baseDir, "transport/http/routes/routes.go")

	files := map[string]string{
		"model.go":      "modules/model.go.tmpl",
		"repository.go": "modules/repository.go.tmpl",
		"service.go":    "modules/service.go.tmpl",
		"handler.go":    "modules/handler.go.tmpl",
	}
	var created []string
	for f := range files {
		created = append(created, filepath.Join(modulePath, f))
	}

	if err := refuseOverwrite(g.Force, created...); err != nil {
		return err
	}
	if !g.Force {
		if err := g.checkRouterField(routesPath, name+"Handler"); err != nil {
			return err
		}
	}

	return atomically([]string{mainPath, routesPath}, created, func() error {
		data := map[string]interface{}{
			"Config":      projectConfig,
			"Model":       model,
			"PackageName": packageName,
		}
		for f, tmpl := range files {
			if err := g.GenerateFromSource(filepath.Join(modulePath, f), tmpl, data); err != nil {
				return err
			}
		}

		moduleImport := fmt.Sprintf("%s/internal/modules/%s", projectConfig.ModuleName, packageName)
		for _, path := range []string{mainPath, routesPath} {
			if err := utils.AddImport(path, moduleImport); err != nil {
				return err
			}
		}

		// Variables are named after the module (blogPost...) while the package
		// qualifier is its lower-case package name (blogpost).
		v := naming.Camel(name)
		inits := []struct{ marker, code string }{
			{"// [REPOS-INIT]", fmt.Sprintf("%sRepo := %s.NewRepository(db)", v, packageName)},
			{"// [SERVICES-INIT]", fmt.Sprintf("%sService := %s.NewService(%sRepo)", v, packageName, v)},
			{"// [HANDLERS-INIT]", fmt.Sprintf("%sHandler := %s.NewHandler(%sService, validator)", v, packageName, v)},
			{"// [ROUTER-HANDLERS-INIT]", fmt.Sprintf("router.%sHandler = %sHandler", name, v)},
		}
		for _, in := range inits {
			if err := utils.InjectCodeAST(mainPath, in.marker, in.code); err != nil {
				return err
			}
		}

		field := name + "Handler"
		if err := utils.AddStructField(routesPath, "Router", field, "*"+packageName+".Handler", ""); err != nil {
			return err
		}
		group := v + "Routes"
		routes := fmt.Sprintf(`%[1]s := api.Group("/%[2]s")
	%[1]s.POST("", r.%[3]s.Create)
	%[1]s.GET("", r.%[3]s.GetAll)
	%[1]s.GET("/:id", r.%[3]s.GetByID)
	%[1]s.PUT("/:id", r.%[3]s.Update)
	%[1]s.DELETE("/:id", r.%[3]s.Delete)`, group, naming.Kebab(naming.Plural(name)), field)
		return utils.InjectCodeAST(routesPath, "// [ROUTES-INIT]", routes)
	})
}
