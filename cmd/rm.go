package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vg/internal/config"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [version]",
	Short: "Remove a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		sdksDir, err := config.GetSdksDir()
		if err != nil {
			fmt.Printf("Error getting sdks dir: %v\n", err)
			os.Exit(1)
		}

		// Check if version exists
		versionPath := filepath.Join(sdksDir, version)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Printf("❌ Go version %s is not installed\n", version)
			fmt.Println("\nRun 'vg list' to see installed versions")
			os.Exit(1)
		}

		// Confirm deletion
		fmt.Printf("Removing Go version %s...\n", version)

		// Show progress animation
		done := make(chan bool)
		go showProgress(done)

		// Perform deletion
		err = os.RemoveAll(versionPath)
		done <- true
		<-done // Wait for animation to finish

		if err != nil {
			fmt.Printf("\n❌ Error removing version: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✅ Successfully removed go%s\n", version)
	},
}

func showProgress(done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			fmt.Print("\r") // Clear the line
			done <- true
			return
		case <-ticker.C:
			fmt.Printf("\r  %s Deleting files...", frames[i%len(frames)])
			i++
		}
	}
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
