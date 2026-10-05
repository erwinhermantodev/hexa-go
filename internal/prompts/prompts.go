package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// PromptForInput prompts user for input with given message
func PromptForInput(prompt string) string {
	fmt.Print(prompt)
	line, _ := stdin.ReadString('\n')
	return strings.TrimSpace(line)
}

// stdin is shared across prompts: a fresh reader per call would discard
// buffered input when stdin is piped.
var stdin = bufio.NewReader(os.Stdin)

// PromptForModels prompts user to define custom models
func PromptForModels() []config.ModelConfig {
	var models []config.ModelConfig

	for {
		answer := PromptForInput("Do you want to add a custom model? (y/n): ")
		if strings.ToLower(answer) != "y" {
			break
		}

		modelName := PromptForInput("Enter model name: ")
		fields := PromptForModelFields(modelName)

		hasRepo := strings.ToLower(PromptForInput("Generate repository? (y/n): ")) == "y"
		hasService := strings.ToLower(PromptForInput("Generate service? (y/n): ")) == "y"
		hasHandler := strings.ToLower(PromptForInput("Generate handler? (y/n): ")) == "y"

		models = append(models, config.ModelConfig{
			Name:       modelName,
			Fields:     fields,
			HasRepo:    hasRepo,
			HasService: hasService,
			HasHandler: hasHandler,
		})
	}

	return models
}

// PromptForModelFields prompts user to define fields for a model
func PromptForModelFields(modelName string) []config.FieldConfig {
	fmt.Printf("Define fields for %s model:\n", modelName)
	fmt.Println("Format: Name:Type[:GormOptions[:Validation]]")
	fmt.Println("Example: Title:string::required,min=5")
	fmt.Println("         Slug:string:unique:required")
	fmt.Println("ID, CreatedAt, UpdatedAt and DeletedAt are added automatically.")
	fmt.Println("Press Enter on empty line to finish.")

	var fields []config.FieldConfig
	seen := map[string]bool{}
	for {
		input := PromptForInput("Field: ")
		if input == "" {
			break
		}
		field, err := utils.ParseField(input)
		if err != nil {
			fmt.Printf("  ❌ %v\n", err)
			continue
		}
		if seen[field.Name] {
			fmt.Printf("  ❌ duplicate field %q\n", field.Name)
			continue
		}
		seen[field.Name] = true
		fields = append(fields, field)
	}

	return utils.WithDefaultFields(fields)
}

// PromptForServices prompts user to define custom services
func PromptForServices() []string {
	var services []string

	for {
		answer := PromptForInput("Do you want to add a custom service? (y/n): ")
		if strings.ToLower(answer) != "y" {
			break
		}

		serviceName := PromptForInput("Enter service name: ")
		services = append(services, serviceName)
	}

	return services
}
