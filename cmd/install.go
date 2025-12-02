package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vg/internal/config"
	"vg/internal/downloader"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		// Normalize version (remove 'go' prefix if present)
		normalizedVersion := strings.TrimPrefix(version, "go")

		distsDir, err := config.GetDistsDir()
		if err != nil {
			fmt.Printf("Error getting dists dir: %v\n", err)
			os.Exit(1)
		}

		sdksDir, err := config.GetSdksDir()
		if err != nil {
			fmt.Printf("Error getting sdks dir: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(filepath.Join(sdksDir, normalizedVersion)); err == nil {
			fmt.Printf("❌ Go %s is already installed\n", normalizedVersion)
			os.Exit(1)
		}

		fmt.Printf("Installing Go %s...\n", normalizedVersion)
		_, err = os.Stat(filepath.Join(distsDir, fmt.Sprintf("go%s.tar.gz", normalizedVersion)))
		if err != nil {
			err = downloader.DownloadAndInstall(normalizedVersion, distsDir, sdksDir)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Archive found at %s, skipping download.\n", filepath.Join(distsDir, fmt.Sprintf("go%s.tar.gz", normalizedVersion)))
		}

		// Create GOPATH directory for this version
		gopath, err := config.GetVersionGopath(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting gopath: %v\n", err)
			os.Exit(1)
		}

		if err := os.MkdirAll(gopath, 0755); err != nil {
			fmt.Printf("Error creating gopath directory: %v\n", err)
			os.Exit(1)
		}

		// Create standard GOPATH subdirectories
		for _, subdir := range []string{"src", "bin", "pkg"} {
			if err := os.MkdirAll(filepath.Join(gopath, subdir), 0755); err != nil {
				fmt.Printf("Error creating gopath subdirectory %s: %v\n", subdir, err)
				os.Exit(1)
			}
		}

		// Initialize GOENV file for this version
		goenvPath, err := config.GetVersionGoenv(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting goenv path: %v\n", err)
			os.Exit(1)
		}

		// Create goenvs directory if it doesn't exist
		goenvsDir := filepath.Dir(goenvPath)
		if err := os.MkdirAll(goenvsDir, 0755); err != nil {
			fmt.Printf("Error creating goenvs directory: %v\n", err)
			os.Exit(1)
		}

		// Create empty GOENV file (GOROOT and GOPATH are set by vg init, not in this file)
		// Users can add custom environment variables here or use 'go env -w KEY=VALUE'
		goenvContent := "# This file is managed by vg.\n# GOROOT and GOPATH are set automatically by 'vg init'.\n# You can add custom environment variables below or use 'go env -w KEY=VALUE'\n"
		if err := os.WriteFile(goenvPath, []byte(goenvContent), 0644); err != nil {
			fmt.Printf("Error creating goenv file: %v\n", err)
			os.Exit(1)
		}

		// Create GOCACHE directory for this version
		gocache, err := config.GetVersionGocache(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting gocache: %v\n", err)
			os.Exit(1)
		}

		if err := os.MkdirAll(gocache, 0755); err != nil {
			fmt.Printf("Error creating gocache directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Created GOPATH: %s\n", gopath)
		fmt.Printf("✅ Created GOENV: %s\n", goenvPath)
		fmt.Printf("✅ Created GOCACHE: %s\n", gocache)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
