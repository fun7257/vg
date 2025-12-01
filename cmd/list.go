package cmd

import (
	"fmt"
	"os"
	"sort"

	"vg/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		sdksDir, err := config.GetSdksDir()
		if err != nil {
			fmt.Printf("Error getting sdks dir: %v\n", err)
			os.Exit(1)
		}

		// Check if sdks directory exists
		if _, err := os.Stat(sdksDir); os.IsNotExist(err) {
			fmt.Println("No Go versions installed yet.")
			return
		}

		// Read directory
		entries, err := os.ReadDir(sdksDir)
		if err != nil {
			fmt.Printf("Error reading sdks directory: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("No Go versions installed yet.")
			return
		}

		// Collect version directories
		var versions []string
		for _, entry := range entries {
			if entry.IsDir() {
				versions = append(versions, entry.Name())
			}
		}

		if len(versions) == 0 {
			fmt.Println("No Go versions installed yet.")
			return
		}

		// Sort versions
		sort.Strings(versions)

		// Display
		fmt.Printf("Installed Go versions (%d):\n", len(versions))
		for _, version := range versions {
			fmt.Printf("  - %s\n", version)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
