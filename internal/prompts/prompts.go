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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

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
	var fields []config.FieldConfig

	fmt.Printf("Define fields for %s model:\n", modelName)
	fmt.Println("Format: field_name field_type [gorm_tag] [json_tag] [validation]")
	fmt.Println("Example: Name string required min=2,max=100")
	fmt.Println("Press Enter on empty line to finish.")

	// Add default fields
	fields = append(fields, config.DefaultModelFields()...)

	for {
		input := PromptForInput("Field: ")
		if input == "" {
			break
		}

		field := utils.ParseFieldInput(input)
		if field.Name != "" {
			fields = append(fields, field)
		}
	}

	return fields
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
