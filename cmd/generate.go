package cmd

import (
	"fmt"
	"os"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/generator"
	"github.com/erwinhermantodev/hexa-go/internal/prompts"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [project-name]",
	Short: "Generate a new Go project",
	Args:  cobra.MaximumNArgs(1),
	Run:   generateProject,
}

func init() {
	generateCmd.Flags().StringP("module", "m", "", "Go module name (e.g., github.com/username/project)")
	generateCmd.Flags().StringP("author", "a", "", "Author name")
	generateCmd.Flags().StringP("description", "d", "", "Project description")
	generateCmd.Flags().BoolP("interactive", "i", false, "Interactive mode for defining models")
	generateCmd.Flags().BoolP("minimal", "", false, "Generate minimal project without auth")
}

func generateProject(cmd *cobra.Command, args []string) {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	} else {
		projectName = prompts.PromptForInput("Enter project name: ")
	}

	moduleName, _ := cmd.Flags().GetString("module")
	if moduleName == "" {
		moduleName = prompts.PromptForInput("Enter module name (e.g., github.com/username/project): ")
	}

	author, _ := cmd.Flags().GetString("author")
	if author == "" {
		author = prompts.PromptForInput("Enter author name: ")
	}

	description, _ := cmd.Flags().GetString("description")
	if description == "" {
		description = prompts.PromptForInput("Enter project description: ")
	}

	interactive, _ := cmd.Flags().GetBool("interactive")
	minimal, _ := cmd.Flags().GetBool("minimal")

	projectConfig := config.ProjectConfig{
		Name:        projectName,
		ModuleName:  moduleName,
		Description: description,
		Author:      author,
		Models:      []config.ModelConfig{},
		Services:    []string{},
	}

	// Add default auth models if not minimal
	if !minimal {
		projectConfig.Models = append(projectConfig.Models,
			config.DefaultUserModel(),
		)
		projectConfig.Services = append(projectConfig.Services, "Auth")
	}

	if interactive {
		projectConfig.Models = append(projectConfig.Models, prompts.PromptForModels()...)
		projectConfig.Services = append(projectConfig.Services, prompts.PromptForServices()...)
	}

	fmt.Printf("Generating project '%s'...\n", projectName)

	gen := generator.New()
	if err := gen.CreateProject(projectConfig); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Project '%s' generated successfully!\n", projectName)
	fmt.Printf("📁 Location: ./%s\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  go run main.go")
}
