package cmd

import (
	"fmt"
	"os"

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

		fmt.Printf("Installing Go %s...\n", version)
		if err := downloader.DownloadAndInstall(version, distsDir, sdksDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
