package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fun7257/vg/internal/config"

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
		// Get symlink paths
		currentLink, err := config.GetCurrentLink()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current link: %v\n", err)
			return
		}

		currentGopathLink, err := config.GetCurrentGopathLink()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current-gopath link: %v\n", err)
			return
		}

		currentGocacheLink, err := config.GetCurrentGocacheLink()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current-gocache link: %v\n", err)
			return
		}

		currentGoenvLink, err := config.GetCurrentGoenvLink()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current-goenv link: %v\n", err)
			return
		}

		// Check if symlinks exist
		if _, err := os.Lstat(currentLink); err != nil {
			fmt.Printf("# vg: No Go version is currently active\n")
			fmt.Printf("# Run 'vg use <version>' to activate a version\n")
			return
		}

		// Get shared GOMODCACHE directory
		gomodcache, err := config.GetGomodcacheDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting gomodcache: %v\n", err)
			return
		}

		// Ensure GOMODCACHE directory exists
		if err := os.MkdirAll(gomodcache, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating gomodcache directory: %v\n", err)
			return
		}

		// Set environment variables pointing to symlinks
		// These symlinks are updated by 'vg use' command
		fmt.Printf("export GOROOT=\"%s\"\n", currentLink)
		fmt.Printf("export GOPATH=\"%s\"\n", currentGopathLink)
		fmt.Printf("export GOCACHE=\"%s\"\n", currentGocacheLink)
		fmt.Printf("export GOENV=\"%s\"\n", currentGoenvLink)
		fmt.Printf("export GOMODCACHE=\"%s\"\n", gomodcache)

		// Set PATH (add GOROOT/bin and GOPATH/bin)
		currentBin := filepath.Join(currentLink, "bin")
		gopathBin := filepath.Join(currentGopathLink, "bin")
		fmt.Printf("export PATH=\"%s:%s:$PATH\"\n", currentBin, gopathBin)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
