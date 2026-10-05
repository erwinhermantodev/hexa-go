package generator

import (
	"fmt"
	"path/filepath"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/naming"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// GenerateServiceFile generates a standalone service file and wires it into
// main.go. A trailing "Service" in the name is dropped, so `add service
// NotificationService` and `add service Notification` are the same.
func (g *Generator) GenerateServiceFile(projectConfig config.ProjectConfig, serviceName string) error {
	serviceName = naming.TrimSuffix(serviceName, "Service")
	if err := naming.Validate(serviceName); err != nil {
		return err
	}
	name := naming.Pascal(serviceName)

	baseDir := projectConfig.Name
	if baseDir == "" {
		baseDir = "."
	}
	servicePath := filepath.Join(baseDir, "service", naming.Snake(name)+".go")
	mainPath := filepath.Join(baseDir, "main.go")

	if err := refuseOverwrite(g.Force, servicePath); err != nil {
		return err
	}

	return atomically([]string{mainPath}, []string{servicePath}, func() error {
		if err := g.GenerateFromSource(servicePath, "core/custom_service.go.tmpl", map[string]interface{}{
			"Config":      projectConfig,
			"ServiceName": name,
		}); err != nil {
			return err
		}
		if err := utils.AddImport(mainPath, projectConfig.ModuleName+"/service"); err != nil {
			return err
		}
		// The blank assignment keeps main.go compiling until the service is used.
		init := fmt.Sprintf("%[1]sService := service.New%[2]sService()\n\t_ = %[1]sService", naming.Camel(name), name)
		return utils.InjectCodeAST(mainPath, "// [SERVICES-INIT]", init)
	})
}
