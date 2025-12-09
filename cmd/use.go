package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fun7257/vg/internal/config"
	"github.com/fun7257/vg/internal/downloader"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		// Normalize version (remove 'go' prefix if present)
		normalizedVersion := strings.TrimPrefix(version, "go")

		sdksDir, err := config.GetSdksDir()
		if err != nil {
			fmt.Printf("Error getting sdks dir: %v\n", err)
			os.Exit(1)
		}

		vgHome, err := config.GetVgHome()
		if err != nil {
			fmt.Printf("Error getting vg home: %v\n", err)
			os.Exit(1)
		}

		// Check if version exists
		versionPath := filepath.Join(sdksDir, normalizedVersion)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			// Try to install it
			fmt.Printf("Go version %s is not installed. Do you want to install it? [y/N] ", normalizedVersion)

			var response string
			_, err := fmt.Scanln(&response)
			if err != nil && err.Error() != "unexpected newline" {
				fmt.Printf("Error scanning response: %v\n", err)
				os.Exit(1)
			}

			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("\nRun 'vg list' to see installed versions")
				fmt.Printf("Run 'vg install %s' to install this version\n", version)
				os.Exit(1)
			}

			distsDir, err := config.GetDistsDir()
			if err != nil {
				fmt.Printf("Error getting dists dir: %v\n", err)
				os.Exit(1)
			}

			// Call downloader
			if err := downloader.DownloadAndInstall(normalizedVersion, distsDir, sdksDir); err != nil {
				fmt.Printf("❌ Failed to install Go %s: %v\n", normalizedVersion, err)
				os.Exit(1)
			}

			fmt.Printf("✅ Automatically installed Go %s\n", normalizedVersion)

		}

		// Get paths for this version
		gopath, err := config.GetVersionGopath(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting gopath: %v\n", err)
			os.Exit(1)
		}

		goenvPath, err := config.GetVersionGoenv(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting goenv path: %v\n", err)
			os.Exit(1)
		}

		gocache, err := config.GetVersionGocache(normalizedVersion)
		if err != nil {
			fmt.Printf("Error getting gocache: %v\n", err)
			os.Exit(1)
		}

		// Verify that GOPATH exists for this version (legacy check)
		if _, err := os.Stat(gopath); os.IsNotExist(err) {
			// Create GOPATH if it doesn't exist (for versions installed before this feature)
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
		}

		// Verify that GOENV exists for this version
		if _, err := os.Stat(goenvPath); os.IsNotExist(err) {
			// Create GOENV if it doesn't exist
			goenvsDir := filepath.Dir(goenvPath)
			if err := os.MkdirAll(goenvsDir, 0755); err != nil {
				fmt.Printf("Error creating goenvs directory: %v\n", err)
				os.Exit(1)
			}
			goenvContent := "# This file is managed by vg.\n# GOROOT and GOPATH are set automatically by 'vg init'.\n# You can add custom environment variables below or use 'go env -w KEY=VALUE'\n"
			if err := os.WriteFile(goenvPath, []byte(goenvContent), 0644); err != nil {
				fmt.Printf("Error creating goenv file: %v\n", err)
				os.Exit(1)
			}
		}

		// Verify that GOCACHE exists for this version
		if _, err := os.Stat(gocache); os.IsNotExist(err) {
			// Create GOCACHE if it doesn't exist
			if err := os.MkdirAll(gocache, 0755); err != nil {
				fmt.Printf("Error creating gocache: %v\n", err)
				os.Exit(1)
			}
		}

		// Helper function to update a symlink
		updateSymlink := func(linkPath, targetPath string, linkName string) error {
			// Remove existing symlink if it exists
			if _, err := os.Lstat(linkPath); err == nil {
				if err := os.Remove(linkPath); err != nil {
					return fmt.Errorf("error removing old %s symlink: %w", linkName, err)
				}
			}
			// Create new symlink
			if err := os.Symlink(targetPath, linkPath); err != nil {
				return fmt.Errorf("error creating %s symlink: %w", linkName, err)
			}
			return nil
		}

		// Update all symlinks
		currentLink, _ := config.GetCurrentLink()
		if err := updateSymlink(currentLink, versionPath, "current"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-gopath"), gopath, "current-gopath"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-gocache"), gocache, "current-gocache"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-goenv"), goenvPath, "current-goenv"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Switched to Go %s\n", normalizedVersion)
		fmt.Println("\nEnvironment variables will be updated automatically via symlinks.")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
