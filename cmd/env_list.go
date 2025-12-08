package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var listEnvCmd = &cobra.Command{
	Use:   "list",
	Short: "List virtual environments for the current Go version",
	Run: func(cmd *cobra.Command, args []string) {
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

		// 2. Get envs dir for this version
		// Since config.GetEnvDir takes (version, name), we need to manually list directory
		envsRoot, err := config.GetEnvsDir()
		if err != nil {
			fmt.Printf("Error getting envs dir: %v\n", err)
			os.Exit(1)
		}
		versionEnvsDir := filepath.Join(envsRoot, currentVersion)

		fmt.Printf("Virtual environments for Go %s:\n\n", currentVersion)

		if _, err := os.Stat(versionEnvsDir); os.IsNotExist(err) {
			fmt.Println("  (none)")
			return
		}

		entries, err := os.ReadDir(versionEnvsDir)
		if err != nil {
			fmt.Printf("Error reading envs directory: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "  NAME\tREMARK")

		count := 0
		for _, entry := range entries {
			if entry.IsDir() {
				count++
				name := entry.Name()
				remark := ""

				remarkPath := filepath.Join(versionEnvsDir, name, "remark.txt")
				if data, err := os.ReadFile(remarkPath); err == nil {
					remark = string(data)
					// Truncate if too long (optional, simple first line check)
				}

				_, _ = fmt.Fprintf(w, "  %s\t%s\n", name, remark)
			}
		}
		_ = w.Flush()

		if count == 0 {
			fmt.Println("  (none)")
		}
	},
}

func init() {
	envCmd.AddCommand(listEnvCmd)
}
