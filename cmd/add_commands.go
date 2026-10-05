package cmd

import (
	"fmt"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/generator"
	"github.com/erwinhermantodev/hexa-go/internal/naming"
	"github.com/erwinhermantodev/hexa-go/internal/prompts"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
	"github.com/spf13/cobra"
)

const fieldFlagHelp = `Model field as Name:Type[:GormOptions[:Validation]], repeatable (e.g. -f "Title:string::required,min=5")`

var addModelCmd = &cobra.Command{
	Use:   "model [model-name]",
	Short: "Add a new model with repository, service, and handler",
	Args:  cobra.ExactArgs(1),
	RunE:  addModel,
}

var addServiceCmd = &cobra.Command{
	Use:   "service [service-name]",
	Short: "Add a new service",
	Args:  cobra.ExactArgs(1),
	RunE:  addService,
}

var addHandlerCmd = &cobra.Command{
	Use:   "handler [handler-name]",
	Short: "Add a new handler",
	Args:  cobra.ExactArgs(1),
	RunE:  addHandler,
}

var addModuleCmd = &cobra.Command{
	Use:   "module [module-name]",
	Short: "Add a new module with co-located repository, service, and handler",
	Args:  cobra.ExactArgs(1),
	RunE:  addModule,
}

func init() {
	addModelCmd.Flags().StringArrayP("fields", "f", nil, fieldFlagHelp)
	addModelCmd.Flags().Bool("no-repo", false, "Skip repository generation (also skips service and handler)")
	addModelCmd.Flags().Bool("no-service", false, "Skip service generation (also skips handler)")
	addModelCmd.Flags().Bool("no-handler", false, "Skip handler generation")
	addModuleCmd.Flags().StringArrayP("fields", "f", nil, fieldFlagHelp)

	for _, c := range []*cobra.Command{addModelCmd, addServiceCmd, addHandlerCmd, addModuleCmd} {
		c.Flags().Bool("force", false, "Overwrite existing files and re-register existing wiring")
	}
}

// newGenerator builds a generator honouring --force.
func newGenerator(cmd *cobra.Command) *generator.Generator {
	force, _ := cmd.Flags().GetBool("force")
	gen := generator.New()
	gen.Force = force
	return gen
}

// currentProject describes the project in the working directory.
func currentProject() (config.ProjectConfig, error) {
	module, err := utils.GetModuleName()
	if err != nil {
		return config.ProjectConfig{}, err
	}
	return config.ProjectConfig{ModuleName: module}, nil
}

// modelFromFlags builds a model from -f flags, prompting when none are given.
func modelFromFlags(cmd *cobra.Command, name string) ([]config.FieldConfig, error) {
	specs, _ := cmd.Flags().GetStringArray("fields")
	fields, err := utils.ParseFields(specs)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return prompts.PromptForModelFields(name), nil
	}
	return utils.WithDefaultFields(fields), nil
}

func addModel(cmd *cobra.Command, args []string) error {
	if err := naming.Validate(args[0]); err != nil {
		return err
	}
	name := naming.Pascal(args[0])

	project, err := currentProject()
	if err != nil {
		return err
	}

	fields, err := modelFromFlags(cmd, name)
	if err != nil {
		return err
	}

	noRepo, _ := cmd.Flags().GetBool("no-repo")
	noService, _ := cmd.Flags().GetBool("no-service")
	noHandler, _ := cmd.Flags().GetBool("no-handler")
	model := config.ModelConfig{
		Name:       name,
		Fields:     fields,
		HasRepo:    !noRepo,
		HasService: !noRepo && !noService,
		HasHandler: !noRepo && !noService && !noHandler,
	}
	project.Models = []config.ModelConfig{model}

	if err := newGenerator(cmd).GenerateModelFiles(project, model); err != nil {
		return fmt.Errorf("generating model: %w", err)
	}

	file := naming.Snake(name)
	fmt.Printf("✅ Model '%s' generated successfully!\n", name)
	fmt.Printf("  📋 Generated model: model/%s.go\n", file)
	if model.HasRepo {
		fmt.Printf("  📝 Generated repository: repository/%s.go\n", file)
	}
	if model.HasService {
		fmt.Printf("  🔧 Generated service: service/%s.go\n", file)
	}
	if model.HasHandler {
		fmt.Printf("  🌐 Generated handler: transport/http/handler/%s_handler.go\n", file)
	}
	return nil
}

func addService(cmd *cobra.Command, args []string) error {
	project, err := currentProject()
	if err != nil {
		return err
	}

	name := naming.Pascal(naming.TrimSuffix(args[0], "Service"))
	if err := newGenerator(cmd).GenerateServiceFile(project, args[0]); err != nil {
		return fmt.Errorf("generating service: %w", err)
	}

	fmt.Printf("✅ Service '%s' generated successfully!\n", name)
	fmt.Printf("  🔧 Generated: service/%s.go\n", naming.Snake(name))
	return nil
}

func addHandler(cmd *cobra.Command, args []string) error {
	project, err := currentProject()
	if err != nil {
		return err
	}

	name := naming.Pascal(naming.TrimSuffix(args[0], "Handler"))
	if err := newGenerator(cmd).GenerateHandlerFile(project, args[0]); err != nil {
		return fmt.Errorf("generating handler: %w", err)
	}

	fmt.Printf("✅ Handler '%s' generated successfully!\n", name)
	fmt.Printf("  🌐 Generated: transport/http/handler/%s_handler.go\n", naming.Snake(name))
	return nil
}

func addModule(cmd *cobra.Command, args []string) error {
	if err := naming.ValidatePackage(args[0]); err != nil {
		return err
	}
	name := naming.Pascal(args[0])

	project, err := currentProject()
	if err != nil {
		return err
	}

	fields, err := modelFromFlags(cmd, name)
	if err != nil {
		return err
	}

	model := config.ModelConfig{
		Name:       name,
		Fields:     fields,
		HasRepo:    true,
		HasService: true,
		HasHandler: true,
	}
	project.Models = []config.ModelConfig{model}

	if err := newGenerator(cmd).GenerateModuleFiles(project, model, args[0]); err != nil {
		return fmt.Errorf("generating module: %w", err)
	}

	fmt.Printf("✅ Module '%s' generated successfully!\n", name)
	fmt.Printf("  📁 Generated: internal/modules/%s\n", naming.Package(name))
	return nil
}
