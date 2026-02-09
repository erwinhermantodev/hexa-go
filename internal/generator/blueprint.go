package generator

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed all:templates/*
var defaultTemplates embed.FS

type Blueprint struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	BasePath     string            `yaml:"base_path"`
	Templates    []TemplateMapping `yaml:"templates"`
	PreGenerate  []string          `yaml:"pre_generate"`
	PostGenerate []string          `yaml:"post_generate"`
}

type TemplateMapping struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func LoadBlueprint(path string) (*Blueprint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var bp Blueprint
	if err := yaml.Unmarshal(data, &bp); err != nil {
		return nil, err
	}

	return &bp, nil
}

func (g *Generator) GetTemplateContent(blueprintDir, sourcePath string) (string, error) {
	// 1. Try external directory
	if blueprintDir != "" {
		fullPath := filepath.Join(blueprintDir, sourcePath)
		content, err := os.ReadFile(fullPath)
		if err == nil {
			return string(content), nil
		}
	}

	// 2. Try embedded defaults
	// Note: since we embed templates/*, the root is 'templates'
	f, err := defaultTemplates.Open(filepath.Join("templates", sourcePath))
	if err != nil {
		return "", fmt.Errorf("template %s not found in embedded FS or external dir: %v", sourcePath, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (g *Generator) ExecuteHooks(hooks []string, dir string) error {
	for _, hook := range hooks {
		cmd := exec.Command("sh", "-c", hook)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute hook '%s': %v", hook, err)
		}
	}
	return nil
}

func (g *Generator) DiscoverPlugins(dir string) ([]string, error) {
	var plugins []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
			plugins = append(plugins, filepath.Join(dir, entry.Name()))
		}
	}
	return plugins, nil
}

func (g *Generator) RunPlugin(pluginPath string, args ...string) error {
	fullArgs := append([]string{"run", pluginPath}, args...)
	cmd := exec.Command("go", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
