package config

import (
	"os"
	"path/filepath"
)

const (
	// DistsDirName stores the downloaded tarballs (distributions).
	DistsDirName      = "dists"
	SdksDirName       = "sdks"
	GopathsDirName    = "gopaths"
	GoenvsDirName     = "goenvs"
	GomodcacheDirName = "gomodcache"
	GocachesDirName   = "gocaches"
	VgDirName         = ".vg"
)

func GetVgHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, VgDirName), nil
}

func GetDistsDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, DistsDirName), nil
}

func GetSdksDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, SdksDirName), nil
}

func GetGopathsDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, GopathsDirName), nil
}

func GetGoenvsDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, GoenvsDirName), nil
}

// GetGomodcacheDir returns the shared GOMODCACHE directory for all versions
func GetGomodcacheDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, GomodcacheDirName), nil
}

// GetVersionGoroot returns the GOROOT path for a specific version
func GetVersionGoroot(version string) (string, error) {
	sdksDir, err := GetSdksDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(sdksDir, version), nil
}

// GetVersionGopath returns the GOPATH path for a specific version
func GetVersionGopath(version string) (string, error) {
	gopathsDir, err := GetGopathsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gopathsDir, version), nil
}

// GetVersionGoenv returns the GOENV file path for a specific version
func GetVersionGoenv(version string) (string, error) {
	goenvsDir, err := GetGoenvsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(goenvsDir, version+".env"), nil
}

// GetGocachesDir returns the directory containing version-specific GOCACHE directories
func GetGocachesDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, GocachesDirName), nil
}

// GetVersionGocache returns the GOCACHE path for a specific version
func GetVersionGocache(version string) (string, error) {
	gocachesDir, err := GetGocachesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gocachesDir, version), nil
}

// GetCurrentLink returns the path to the 'current' symlink (GOROOT)
func GetCurrentLink() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, "current"), nil
}

// GetCurrentGopathLink returns the path to the 'current-gopath' symlink
func GetCurrentGopathLink() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, "current-gopath"), nil
}

// GetCurrentGocacheLink returns the path to the 'current-gocache' symlink
func GetCurrentGocacheLink() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, "current-gocache"), nil
}

// GetCurrentGoenvLink returns the path to the 'current-goenv' symlink
func GetCurrentGoenvLink() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, "current-goenv"), nil
}

const (
	// EnvsDirName stores the virtual environments.
	EnvsDirName = "envs"
)

// GetEnvsDir returns the directory containing virtual environments
func GetEnvsDir() (string, error) {
	vgHome, err := GetVgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(vgHome, EnvsDirName), nil
}

// GetEnvDir returns the directory for a specific environment under a specific Go version
func GetEnvDir(version, name string) (string, error) {
	envsDir, err := GetEnvsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(envsDir, version, name), nil
}
