package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fun7257/vg/internal/config"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current Go version and environment status",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get Go Version (from 'current' symlink)
		currentLink, err := config.GetCurrentLink()
		if err != nil {
			fmt.Printf("Error getting current link: %v\n", err)
			return
		}

		targetGoroot, err := os.Readlink(currentLink)
		goVersion := ""
		if err == nil {
			goVersion = filepath.Base(targetGoroot)
		} else if os.IsNotExist(err) {
			goVersion = "Not set (run 'vg use')"
		} else {
			goVersion = fmt.Sprintf("Error reading link: %v", err)
		}

		// 2. Get Paths (resolve symlinks)
		gopathLink, _ := config.GetCurrentGopathLink()
		targetGopath, _ := os.Readlink(gopathLink)

		gocacheLink, _ := config.GetCurrentGocacheLink()
		targetGocache, _ := os.Readlink(gocacheLink)

		goenvLink, _ := config.GetCurrentGoenvLink()
		targetGoenv, _ := os.Readlink(goenvLink)

		// 3. Determine Environment
		// Structure for envs: .../envs/<version>/<name>/gopath
		// Structure for standard: .../gopaths/<version>
		envName := "(global)"

		// Use GOPATH to detect if we are in an environment
		if strings.Contains(targetGopath, "/envs/") {
			// Path looks like /.../envs/1.25.5/my-env/gopath
			// We want to extract "my-env"
			// parent is my-env, grandparent is 1.25.5, great-grandparent is envs
			dir := filepath.Dir(targetGopath) // .../envs/1.25.5/my-env
			envName = filepath.Base(dir)      // my-env
		}

		// 4. Output
		fmt.Printf("Go Version:  %s\n", goVersion)
		fmt.Printf("Environment: %s\n", envName)

		// Show remark if exists
		if envName != "(global)" {
			envDir := filepath.Dir(targetGopath)
			remarkPath := filepath.Join(envDir, "remark.txt")
			if data, err := os.ReadFile(remarkPath); err == nil && len(data) > 0 {
				fmt.Printf("Remark:      %s\n", string(data))
			}
		}

		fmt.Println()
		fmt.Printf("GOROOT:      %s\n", targetGoroot)
		fmt.Printf("GOPATH:      %s\n", targetGopath)
		fmt.Printf("GOCACHE:     %s\n", targetGocache)
		fmt.Printf("GOENV:       %s\n", targetGoenv)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
