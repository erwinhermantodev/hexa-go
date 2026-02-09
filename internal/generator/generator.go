package generator

import (
	"fmt"
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

// GenerateFromSource generates a file from a source template path
func (g *Generator) GenerateFromSource(targetPath, sourcePath string, data interface{}) error {
	content, err := g.GetTemplateContent("", sourcePath)
	if err != nil {
		return err
	}
	return g.CreateFileFromTemplate(targetPath, content, data)
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

// InjectCode adds code to an existing file before a marker
func (g *Generator) InjectCode(filePath, marker, code string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	strContent := string(content)

	// Find the marker and its indentation
	index := strings.Index(strContent, marker)
	if index == -1 {
		return fmt.Errorf("marker %s not found in %s", marker, filePath)
	}

	// Detect indentation
	startOfLine := strings.LastIndex(strContent[:index], "\n") + 1
	indentation := strContent[startOfLine:index]

	// Avoid duplicate injection
	if strings.Contains(strContent, "\n"+indentation+code) {
		return nil
	}

	newCode := code + "\n" + indentation + marker
	newContent := strings.Replace(strContent, marker, newCode, 1)
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// InjectImport adds an import to a file
func (g *Generator) InjectImport(filePath, importPath string) error {
	marker := "// [IMPORTS]"
	code := fmt.Sprintf("\"%s\"", importPath)
	return g.InjectCode(filePath, marker, code)
}
