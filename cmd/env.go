package cmd

import (
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage virtual environments",
	Long:  `Manage virtual environments for Go versions.`,
}

func init() {
	rootCmd.AddCommand(envCmd)

	// Register subcommands
	envCmd.AddCommand(newCmd)
	envCmd.AddCommand(loadCmd)

	// Flags
	newCmd.Flags().StringP("message", "m", "", "Add a remark to the environment")
}
