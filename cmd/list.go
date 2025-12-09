package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed Go versions and their virtual environments",
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

		// Get envs root dir
		envsRoot, _ := config.GetEnvsDir()

		// Display
		fmt.Printf("Installed Go versions (%d):\n", len(versions))
		for _, version := range versions {
			fmt.Printf("  - %s\n", version)

			// Check for virtual environments
			if envsRoot != "" {
				versionEnvsDir := filepath.Join(envsRoot, version)
				if envEntries, err := os.ReadDir(versionEnvsDir); err == nil {
					var envs []string
					for _, envEntry := range envEntries {
						if envEntry.IsDir() {
							envs = append(envs, envEntry.Name())
						}
					}

					if len(envs) > 0 {
						sort.Strings(envs)
						for _, envName := range envs {
							remark := ""
							remarkPath := filepath.Join(versionEnvsDir, envName, "remark.txt")
							if data, err := os.ReadFile(remarkPath); err == nil {
								remark = strings.TrimSpace(string(data))
							}

							if remark != "" {
								fmt.Printf("      * %s (%s)\n", envName, remark)
							} else {
								fmt.Printf("      * %s\n", envName)
							}
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
