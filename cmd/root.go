package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is set at build time: go build -ldflags "-X github.com/erwinhermantodev/hexa-go/cmd.version=v0.1.0".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "hexa-go",
	Short: "Generate a flexible Go Hexagonal architecture project with custom models, repositories, services, and handlers",
	// Errors are reported by main; printing usage on top of a generator error is noise.
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version,
	Long:          `A CLI tool to generate a Go project starter pack with customizable authentication service, REST API, and business logic components.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			rootCmd.Version = info.Main.Version
		}
	}
	rootCmd.AddCommand(generateCmd, addCmd)
}
