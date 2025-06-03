package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hexa-go",
	Short: "Generate a flexible Go Hexagonal architecture project with custom models, repositories, services, and handlers",
	Long:  `A CLI tool to generate a Go project starter pack with customizable authentication service, REST API, and business logic components.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(generateCmd, addCmd)
}
