package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/erwinhermantodev/hexa-go/internal/config"
)

// Generator handles code generation
type Generator struct{}

// New creates a new Generator instance
func New() *Generator {
	return &Generator{}
}

// CreateFileFromTemplate creates a file from a template
func (g *Generator) CreateFileFromTemplate(filePath, tmplContent string, data interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpl, err := template.New("file").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		"Title":   strings.Title,
		"contains": func(fields []config.FieldConfig, fieldType string) bool {
			for _, field := range fields {
				if strings.Contains(field.Type, fieldType) {
					return true
				}
			}
			return false
		},
	}).Parse(tmplContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
