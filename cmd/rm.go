package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [version]",
	Short: "Remove a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		// Normalize version (remove 'go' prefix if present)
		normalizedVersion := strings.TrimPrefix(version, "go")

		sdksDir, err := config.GetSdksDir()
		if err != nil {
			fmt.Printf("Error getting sdks dir: %v\n", err)
			os.Exit(1)
		}

		// Check if version exists
		versionPath := filepath.Join(sdksDir, normalizedVersion)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Printf("❌ Go version %s is not installed\n", version)
			fmt.Println("\nRun 'vg list' to see installed versions")
			os.Exit(1)
		}

		// Check if this version is currently in use
		currentLink, err := config.GetCurrentLink()
		if err != nil {
			fmt.Printf("Error getting current link: %v\n", err)
			os.Exit(1)
		}

		if target, err := os.Readlink(currentLink); err == nil {
			currentVersion := filepath.Base(target)
			if currentVersion == normalizedVersion {
				fmt.Printf("❌ Cannot remove Go %s: it is currently in use\n", normalizedVersion)
				fmt.Printf("\nTo remove this version, first switch to another version:\n")
				fmt.Printf("  vg use <other-version>\n")
				fmt.Printf("  vg rm %s\n", normalizedVersion)
				os.Exit(1)
			}
		}

		// Confirm deletion
		fmt.Printf("Removing Go version %s...\n", normalizedVersion)

		// Show progress animation
		done := make(chan bool)
		go showProgress(done)

		// Perform deletion of SDK
		err = os.RemoveAll(versionPath)
		if err != nil {
			done <- true
			<-done
			fmt.Printf("\n❌ Error removing SDK: %v\n", err)
			os.Exit(1)
		}

		// Remove GOPATH for this version
		gopath, err := config.GetVersionGopath(normalizedVersion)
		if err == nil {
			if _, err := os.Stat(gopath); err == nil {
				if err := os.RemoveAll(gopath); err != nil {
					done <- true
					<-done
					fmt.Printf("\n⚠️  Warning: Error removing GOPATH: %v\n", err)
				}
			}
		}

		// Remove GOENV for this version
		goenvPath, err := config.GetVersionGoenv(normalizedVersion)
		if err == nil {
			if _, err := os.Stat(goenvPath); err == nil {
				if err := os.Remove(goenvPath); err != nil {
					done <- true
					<-done
					fmt.Printf("\n⚠️  Warning: Error removing GOENV: %v\n", err)
				}
			}
		}

		// Remove GOCACHE for this version
		gocache, err := config.GetVersionGocache(normalizedVersion)
		if err == nil {
			if _, err := os.Stat(gocache); err == nil {
				if err := os.RemoveAll(gocache); err != nil {
					done <- true
					<-done
					fmt.Printf("\n⚠️  Warning: Error removing GOCACHE: %v\n", err)
				}
			}
		}

		done <- true
		<-done // Wait for animation to finish

		fmt.Printf("\n✅ Successfully removed Go %s (SDK, GOPATH, GOENV, and GOCACHE)\n", normalizedVersion)
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
