package generator

import (
	"fmt"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/naming"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateHandlerFile generates a standalone handler file and wires it into
// main.go. A trailing "Handler" in the name is dropped.
func (g *Generator) GenerateHandlerFile(projectConfig config.ProjectConfig, handlerName string) error {
	handlerName = naming.TrimSuffix(handlerName, "Handler")
	if err := naming.Validate(handlerName); err != nil {
		return err
	}
	name := naming.Pascal(handlerName)

	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}
	handlerPath := filepath.Join(baseDir, "transport/http/handler", naming.Snake(name)+"_handler.go")
	mainPath := filepath.Join(baseDir, "main.go")

	if err := refuseOverwrite(g.Force, handlerPath); err != nil {
		return err
	}

	return atomically([]string{mainPath}, []string{handlerPath}, func() error {
		if err := g.GenerateFromSource(handlerPath, "core/custom_handler.go.tmpl", map[string]interface{}{
			"Config":      projectConfig,
			"HandlerName": name,
		}); err != nil {
			return err
		}
		if err := utils.AddImport(mainPath, projectConfig.ModuleName+"/transport/http/handler"); err != nil {
			return err
		}
		init := fmt.Sprintf("%[1]sHandler := handler.New%[2]sHandler(validator)\n\t_ = %[1]sHandler", naming.Camel(name), name)
		return utils.InjectCodeAST(mainPath, "// [HANDLERS-INIT]", init)
	})
}
