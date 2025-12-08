package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new virtual environment",
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
			fmt.Printf("Please run 'vg use <version>' first\n")
			os.Exit(1)
		}
		currentVersion := filepath.Base(target)

		// 2. Check if env already exists
		envDir, err := config.GetEnvDir(currentVersion, envName)
		if err != nil {
			fmt.Printf("Error getting env dir: %v\n", err)
			os.Exit(1)
		}

		// Check if it exists for this version
		if _, err := os.Stat(envDir); err == nil {
			fmt.Printf("❌ Environment '%s' already exists for Go %s\n", envName, currentVersion)
			os.Exit(1)
		}

		fmt.Printf("Creating virtual environment '%s' using Go %s...\n", envName, currentVersion)

		// 3. Create env directory structure (envs/<version>/<name>)
		if err := os.MkdirAll(envDir, 0755); err != nil {
			fmt.Printf("Error creating env directory: %v\n", err)
			os.Exit(1)
		}

		// Create GOPATH
		gopath := filepath.Join(envDir, "gopath")
		if err := os.MkdirAll(gopath, 0755); err != nil {
			fmt.Printf("Error creating gopath: %v\n", err)
			os.Exit(1)
		}
		for _, subdir := range []string{"src", "bin", "pkg"} {
			if err := os.MkdirAll(filepath.Join(gopath, subdir), 0755); err != nil {
				fmt.Printf("Error creating gopath subdirectory %s: %v\n", subdir, err)
				os.Exit(1)
			}
		}

		// Create GOCACHE
		gocache := filepath.Join(envDir, "gocache")
		if err := os.MkdirAll(gocache, 0755); err != nil {
			fmt.Printf("Error creating gocache: %v\n", err)
			os.Exit(1)
		}

		// Create GOENV
		goenvPath := filepath.Join(envDir, "goenv")
		goenvContent := fmt.Sprintf("# Environment '%s' (Go %s)\n# Managed by vg.\n", envName, currentVersion)
		if err := os.WriteFile(goenvPath, []byte(goenvContent), 0644); err != nil {
			fmt.Printf("Error creating goenv file: %v\n", err)
			os.Exit(1)
		}

		// Save remark if provided
		remark, _ := cmd.Flags().GetString("message")
		if remark != "" {
			if err := os.WriteFile(filepath.Join(envDir, "remark.txt"), []byte(remark), 0644); err != nil {
				fmt.Printf("⚠️  Warning: Failed to save remark: %v\n", err)
			}
		}

		fmt.Printf("✅ Created virtual environment '%s'\n", envName)
		fmt.Printf("\nActivate it with:\n  vg load %s\n", envName)
	},
}
