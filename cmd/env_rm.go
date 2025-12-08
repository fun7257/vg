package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var rmEnvCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove a virtual environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]

		// 1. Get current Go version
		currentLink, err := config.GetCurrentLink()
		if err != nil {
			fmt.Printf("Error getting current link: %v\n", err)
			os.Exit(1)
		}

		target, err := os.Readlink(currentLink)
		if err != nil {
			fmt.Printf("❌ No Go version is currently active\n")
			os.Exit(1)
		}
		currentVersion := filepath.Base(target)

		// 2. Resolve Environment Path
		envDir, err := config.GetEnvDir(currentVersion, envName)
		if err != nil {
			fmt.Printf("Error getting env dir: %v\n", err)
			os.Exit(1)
		}

		// 3. Check existence
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			fmt.Printf("❌ Environment '%s' does not exist for Go %s\n", envName, currentVersion)
			os.Exit(1)
		}

		// 4. Check if currently active (safeguard)
		// We check if the current GOPATH matches this environment
		currentGopathLink, _ := config.GetCurrentGopathLink()
		targetGopath, _ := os.Readlink(currentGopathLink)

		// Ensure absolute paths for comparison
		absEnvDir, _ := filepath.Abs(envDir)
		absTargetGopath, _ := filepath.Abs(targetGopath)
		// The active environment root is the parent of the GOPATH (as created by vg new)
		activeEnvDir := filepath.Dir(absTargetGopath)

		if activeEnvDir == absEnvDir {
			fmt.Printf("❌ Cannot remove active environment '%s'\n", envName)
			fmt.Println("Please run 'vg env exit' first.")
			os.Exit(1)
		}

		// 5. Remove
		fmt.Printf("Removing environment '%s' (Go %s)...\n", envName, currentVersion)
		if err := os.RemoveAll(envDir); err != nil {
			fmt.Printf("Error removing environment: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Removed environment '%s'\n", envName)
	},
}

func init() {
	envCmd.AddCommand(rmEnvCmd)
}
