package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [env_name]",
	Short: "Load a virtual environment for the current Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]

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

		// 2. Find environment in current version
		envDir, err := config.GetEnvDir(currentVersion, envName)
		if err != nil {
			fmt.Printf("Error getting env path: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			fmt.Printf("❌ Environment '%s' not found for Go %s\n", envName, currentVersion)
			fmt.Printf("Run 'vg new %s' to create it\n", envName)
			os.Exit(1)
		}

		// Verify SDK exists (sanity check, it should matched currentVersion which is active)
		versionPath := filepath.Join(sdksDir, currentVersion)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Printf("❌ Critical Error: Current SDK %s is missing?\n", currentVersion)
			os.Exit(1)
		}

		// Prepare paths
		targetGopath := filepath.Join(envDir, "gopath")
		targetGoenv := filepath.Join(envDir, "goenv")
		targetGocache := filepath.Join(envDir, "gocache")

		// Helper function to update a symlink
		updateSymlink := func(linkPath, targetPath string, linkName string) error {
			if _, err := os.Lstat(linkPath); err == nil {
				if err := os.Remove(linkPath); err != nil {
					return fmt.Errorf("error removing old %s symlink: %w", linkName, err)
				}
			}
			if err := os.Symlink(targetPath, linkPath); err != nil {
				return fmt.Errorf("error creating %s symlink: %w", linkName, err)
			}
			return nil
		}

		// Update symlinks (NOTE: 'current' link to SDK does NOT change, as we are staying on currentVersion)

		if err := updateSymlink(filepath.Join(vgHome, "current-gopath"), targetGopath, "current-gopath"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-gocache"), targetGocache, "current-gocache"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-goenv"), targetGoenv, "current-goenv"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Loaded environment '%s' (Go %s)\n", envName, currentVersion)
		fmt.Println("\nEnvironment variables will be updated automatically via symlinks.")
	},
}
