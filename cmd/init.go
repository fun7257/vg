package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"vg/internal/config"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate shell configuration",
	Long: `Generate shell configuration to initialize vg environment.
Add the following to your shell profile (e.g., ~/.zshrc or ~/.bashrc):

  eval "$(vg init)"
`,
	Run: func(cmd *cobra.Command, args []string) {
		vgHome, err := config.GetVgHome()
		if err != nil {
			// If we can't get the home dir, we can't really do anything useful.
			// But since this is meant to be evaluated, printing an error to stdout
			// might break the shell. We'll print to stderr.
			fmt.Fprintf(os.Stderr, "Error getting vg home: %v\n", err)
			os.Exit(1)
		}

		currentBin := filepath.Join(vgHome, "current", "bin")

		// Output the export command
		// We use the absolute path to ensure it works regardless of context
		fmt.Printf("export PATH=\"%s:$PATH\"\n", currentBin)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
