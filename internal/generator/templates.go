package generator

import (
	"embed"
	"fmt"
	"io/fs"
)

// defaultTemplates holds every template the CLI renders. They are compiled
// into the binary, so generation never reads templates or runs commands from
// the user's machine.
//
//go:embed all:templates/*
var defaultTemplates embed.FS

// GetTemplateContent returns the embedded template at sourcePath, relative to
// the templates directory (for example "core/model.go.tmpl").
func (g *Generator) GetTemplateContent(sourcePath string) (string, error) {
	content, err := fs.ReadFile(defaultTemplates, "templates/"+sourcePath)
	if err != nil {
		return "", fmt.Errorf("template %s not found: %w", sourcePath, err)
	}
	return string(content), nil
}
