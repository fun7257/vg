package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"vg/internal/config"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

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
		versionPath := filepath.Join(sdksDir, version)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Printf("❌ Go version %s is not installed\n", version)
			fmt.Println("\nRun 'vg list' to see installed versions")
			fmt.Printf("Run 'vg install %s' to install this version\n", version)
			os.Exit(1)
		}

		// Create or update the 'current' symlink
		currentLink := filepath.Join(vgHome, "current")

		// Remove existing symlink if it exists
		if _, err := os.Lstat(currentLink); err == nil {
			if err := os.Remove(currentLink); err != nil {
				fmt.Printf("❌ Error removing old symlink: %v\n", err)
				os.Exit(1)
			}
		}

		// Create new symlink
		if err := os.Symlink(versionPath, currentLink); err != nil {
			fmt.Printf("❌ Error creating symlink: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Switched to Go %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
