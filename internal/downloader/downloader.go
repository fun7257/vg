package downloader

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const BaseURL = "https://go.dev/dl/"

func DownloadAndInstall(version, distsDir, sdksDir string) error {
	// 1. Construct URL
	// e.g., go1.25.4.darwin-arm64.tar.gz
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Handle version string: ensure it starts with "go" or just number
	verStr := version
	if !strings.HasPrefix(verStr, "go") {
		verStr = "go" + verStr
	}

	filename := fmt.Sprintf("%s.%s-%s.tar.gz", verStr, goos, goarch)
	url := BaseURL + filename

	// 2. Check if already installed
	// The SDK will be extracted to sdksDir/go<version> usually, or we rename it.
	// Let's say we want sdksDir/<version>
	// Note: The tarball contains a "go" directory at the root.

	// We'll use the raw version number for the directory name, e.g., "1.25.4"
	installPath := filepath.Join(sdksDir, strings.TrimPrefix(verStr, "go"))
	if _, err := os.Stat(installPath); err == nil {
		return fmt.Errorf("version %s is already installed at %s", version, installPath)
	}

	// 3. Download
	if err := os.MkdirAll(distsDir, 0755); err != nil {
		return fmt.Errorf("failed to create dists dir: %w", err)
	}

	filePath := filepath.Join(distsDir, filename)

	// Track if we downloaded it fresh
	downloaded := false

	// Check if file exists and is valid? For now just check existence.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		downloaded = true
		fmt.Printf("Downloading %s...\n", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		switch resp.StatusCode {
		case http.StatusOK:
			// OK
		case http.StatusNotFound:
			return fmt.Errorf("version %s not found", version)
		default:
			return fmt.Errorf("failed to download: %s", resp.Status)
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)

		if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
			return err
		}
	} else {
		fmt.Printf("Archive found at %s, skipping download.\n", filePath)
	}

	// 4. Extract
	fmt.Printf("\nExtracting to %s...\n", installPath)
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	// We need to close f before we might remove it in error handling
	fClosed := false
	defer func() {
		if !fClosed {
			_ = f.Close()
		}
	}()

	if err := extractTarGz(f, installPath); err != nil {
		// Close file handle first so we can remove it if needed
		_ = f.Close()
		fClosed = true

		if !downloaded {
			// If we didn't download it just now, maybe the cache is corrupt.
			// Especially "unexpected EOF" suggests truncation.
			fmt.Printf("Extraction failed (%v). The cached archive might be corrupt.\n", err)
			fmt.Printf("Removing %s and retrying...\n", filePath)

			_ = os.Remove(filePath)
			_ = os.RemoveAll(installPath)

			// Retry
			return DownloadAndInstall(version, distsDir, sdksDir)
		}

		// Cleanup on failure
		_ = os.RemoveAll(installPath)
		return err
	}

	fmt.Println("Done!")
	return nil
}

func extractTarGz(r io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() {
		_ = gzr.Close()
	}()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// The tarball has "go/..." structure. We want to strip "go/".
		name := header.Name
		if strings.HasPrefix(name, "go/") {
			name = strings.TrimPrefix(name, "go/")
		} else {
			// Skip files not in "go/" (unlikely for official builds)
			continue
		}

		if name == "" {
			continue
		}

		target := filepath.Join(destDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				_ = f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		}
	}
	return nil
}
