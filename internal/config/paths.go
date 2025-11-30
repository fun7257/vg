package config

import (
	"os"
	"path/filepath"
)

const (
	// DistsDirName stores the downloaded tarballs (distributions).
	DistsDirName = "dists"
	SdksDirName  = "sdks"
	VgDirName    = ".vg"
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
