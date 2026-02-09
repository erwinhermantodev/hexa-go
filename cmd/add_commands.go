package cmd

import (
	"fmt"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/generator"
	"github.com/erwinhermantodev/hexa-go/internal/prompts"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
	"github.com/spf13/cobra"
)

var addModelCmd = &cobra.Command{
	Use:   "model [model-name]",
	Short: "Add a new model with repository, service, and handler",
	Args:  cobra.ExactArgs(1),
	Run:   addModel,
}

var addServiceCmd = &cobra.Command{
	Use:   "service [service-name]",
	Short: "Add a new service",
	Args:  cobra.ExactArgs(1),
	Run:   addService,
}

var addHandlerCmd = &cobra.Command{
	Use:   "handler [handler-name]",
	Short: "Add a new handler",
	Args:  cobra.ExactArgs(1),
	Run:   addHandler,
}

var addModuleCmd = &cobra.Command{
	Use:   "module [module-name]",
	Short: "Add a new module with co-located repository, service, and handler",
	Args:  cobra.ExactArgs(1),
	Run:   addModule,
}

func init() {
	addModelCmd.Flags().StringSliceP("fields", "f", []string{}, "Model fields (format: name:type:tag:validation)")
	addModelCmd.Flags().BoolP("no-repo", "", false, "Skip repository generation")
	addModelCmd.Flags().BoolP("no-service", "", false, "Skip service generation")
	addModelCmd.Flags().BoolP("no-handler", "", false, "Skip handler generation")
	addModuleCmd.Flags().StringSliceP("fields", "f", []string{}, "Model fields (format: name:type:tag:validation)")
}

func addModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	fields, _ := cmd.Flags().GetStringSlice("fields")
	noRepo, _ := cmd.Flags().GetBool("no-repo")
	noService, _ := cmd.Flags().GetBool("no-service")
	noHandler, _ := cmd.Flags().GetBool("no-handler")

	if !utils.FileExists("go.mod") {
		fmt.Println("❌ No go.mod found. Please run this command in a Go project directory.")
		return
	}

	modelConfig := config.ModelConfig{
		Name:       modelName,
		Fields:     utils.ParseFieldsFromFlags(fields),
		HasRepo:    !noRepo,
		HasService: !noService,
		HasHandler: !noHandler,
	}

	if len(modelConfig.Fields) == 0 {
		modelConfig.Fields = prompts.PromptForModelFields(modelName)
	}

	projectConfig := config.ProjectConfig{
		ModuleName: utils.GetModuleName(),
		Models:     []config.ModelConfig{modelConfig},
	}

	gen := generator.New()
	if err := gen.GenerateModelFiles(projectConfig, modelConfig); err != nil {
		fmt.Printf("❌ Error generating model: %v\n", err)
		return
	}

	fmt.Printf("✅ Model '%s' generated successfully!\n", modelName)
	if modelConfig.HasRepo {
		fmt.Printf("  📝 Generated repository: repository/%s.go\n", strings.ToLower(modelName))
	}
	if modelConfig.HasService {
		fmt.Printf("  🔧 Generated service: service/%s.go\n", strings.ToLower(modelName))
	}
	if modelConfig.HasHandler {
		fmt.Printf("  🌐 Generated handler: transport/http/handler/%s_handler.go\n", strings.ToLower(modelName))
	}
	fmt.Printf("  📋 Generated model: model/%s.go\n", strings.ToLower(modelName))
}

func addService(cmd *cobra.Command, args []string) {
	serviceName := args[0]

	if !utils.FileExists("go.mod") {
		fmt.Println("❌ No go.mod found. Please run this command in a Go project directory.")
		return
	}

	projectConfig := config.ProjectConfig{
		ModuleName: utils.GetModuleName(),
		Services:   []string{serviceName},
	}

	gen := generator.New()
	if err := gen.GenerateServiceFile(projectConfig, serviceName); err != nil {
		fmt.Printf("❌ Error generating service: %v\n", err)
		return
	}

	fmt.Printf("✅ Service '%s' generated successfully!\n", serviceName)
	fmt.Printf("  🔧 Generated: service/%s.go\n", strings.ToLower(serviceName))
}

func addHandler(cmd *cobra.Command, args []string) {
	handlerName := args[0]

	if !utils.FileExists("go.mod") {
		fmt.Println("❌ No go.mod found. Please run this command in a Go project directory.")
		return
	}

	projectConfig := config.ProjectConfig{
		ModuleName: utils.GetModuleName(),
	}

	gen := generator.New()
	if err := gen.GenerateHandlerFile(projectConfig, handlerName); err != nil {
		fmt.Printf("❌ Error generating handler: %v\n", err)
		return
	}

	fmt.Printf("✅ Handler '%s' generated successfully!\n", handlerName)
	fmt.Printf("  🌐 Generated: transport/http/handler/%s_handler.go\n", strings.ToLower(handlerName))
}

func addModule(cmd *cobra.Command, args []string) {
	moduleName := args[0]
	fields, _ := cmd.Flags().GetStringSlice("fields")

	if !utils.FileExists("go.mod") {
		fmt.Println("❌ No go.mod found. Please run this command in a Go project directory.")
		return
	}

	modelConfig := config.ModelConfig{
		Name:       strings.Title(strings.ToLower(moduleName)),
		Fields:     utils.ParseFieldsFromFlags(fields),
		HasRepo:    true,
		HasService: true,
		HasHandler: true,
	}

	if len(modelConfig.Fields) == 0 {
		modelConfig.Fields = prompts.PromptForModelFields(modelConfig.Name)
	}

	projectConfig := config.ProjectConfig{
		ModuleName: utils.GetModuleName(),
		Models:     []config.ModelConfig{modelConfig},
	}

	gen := generator.New()
	if err := gen.GenerateModuleFiles(projectConfig, modelConfig, moduleName); err != nil {
		fmt.Printf("❌ Error generating module: %v\n", err)
		return
	}

	fmt.Printf("✅ Module '%s' generated successfully!\n", moduleName)
	fmt.Printf("  📁 Generated: internal/modules/%s\n", strings.ToLower(moduleName))
}
