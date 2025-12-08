package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var exitEnvCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit virtual environment and return to global context",
	Run: func(cmd *cobra.Command, args []string) {
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
			os.Exit(1)
		}
		currentVersion := filepath.Base(target)

		// 2. Resolve Global Paths for this version
		globalGopath, err := config.GetVersionGopath(currentVersion)
		if err != nil {
			fmt.Printf("Error getting global gopath: %v\n", err)
			os.Exit(1)
		}
		globalGocache, err := config.GetVersionGocache(currentVersion)
		if err != nil {
			fmt.Printf("Error getting global gocache: %v\n", err)
			os.Exit(1)
		}
		globalGoenv, err := config.GetVersionGoenv(currentVersion)
		if err != nil {
			fmt.Printf("Error getting global goenv: %v\n", err)
			os.Exit(1)
		}

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

		// 3. Reset symlinks to global
		if err := updateSymlink(filepath.Join(vgHome, "current-gopath"), globalGopath, "current-gopath"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-gocache"), globalGocache, "current-gocache"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		if err := updateSymlink(filepath.Join(vgHome, "current-goenv"), globalGoenv, "current-goenv"); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Exited virtual environment. Now using global Go %s context.\n", currentVersion)
		fmt.Println("\nEnvironment variables will be updated automatically via symlinks.")
	},
}

func init() {
	envCmd.AddCommand(exitEnvCmd)
}
