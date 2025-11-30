package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vg",
	Short: "vg is a Virtual Go environment manager",
	Long: `vg is a Virtual Go environment manager.

A Fast and Flexible Go Version Manager that helps you manage multiple Go versions per project.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
