package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add components to existing project",
}

func init() {
	addCmd.AddCommand(addModelCmd, addServiceCmd, addHandlerCmd)
}
