package utils

import (
	"fmt"
	"strings"

	"github.com/erwinhermantodev/go-starter-cli/internal/config"
)

// ParseFieldInput parses a field input string into FieldConfig
func ParseFieldInput(input string) config.FieldConfig {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return config.FieldConfig{}
	}

	field := config.FieldConfig{
		Name: parts[0],
		Type: parts[1],
	}

	// Build tag
	var tagParts []string
	var validate string

	for i := 2; i < len(parts); i++ {
		part := parts[i]
		if strings.Contains(part, "=") || part == "required" || part == "email" {
			validate = part
		} else {
			tagParts = append(tagParts, part)
		}
	}

	// Generate default tags
	jsonTag := fmt.Sprintf("json:\"%s\"", strings.ToLower(field.Name))
	var gormTag string

	if field.Name == "ID" {
		gormTag = "gorm:\"primaryKey\""
	} else if len(tagParts) > 0 {
		gormTag = fmt.Sprintf("gorm:\"%s\"", strings.Join(tagParts, ";"))
	}

	if gormTag != "" {
		field.Tag = fmt.Sprintf("`%s %s`", gormTag, jsonTag)
	} else {
		field.Tag = fmt.Sprintf("`%s`", jsonTag)
	}

	field.Validate = validate

	return field
}

// ParseFieldsFromFlags parses field flags into FieldConfig slice
func ParseFieldsFromFlags(fields []string) []config.FieldConfig {
	var fieldConfigs []config.FieldConfig

	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) >= 2 {
			fieldConfig := config.FieldConfig{
				Name: parts[0],
				Type: parts[1],
			}

			if len(parts) > 2 && parts[2] != "" {
				fieldConfig.Tag = parts[2]
			} else {
				jsonTag := fmt.Sprintf("json:\"%s\"", strings.ToLower(fieldConfig.Name))
				fieldConfig.Tag = fmt.Sprintf("`%s`", jsonTag)
			}

			if len(parts) > 3 {
				fieldConfig.Validate = parts[3]
			}

			fieldConfigs = append(fieldConfigs, fieldConfig)
		}
	}

	return fieldConfigs
}
