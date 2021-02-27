package pkg

import (
	"errors"
	"runtime"
	"strings"

	"golang.org/x/mod/semver"
)

type dlInfo struct {
	Name    string
	Version string
	URL     string
	URLSHA  string
	osArch  string
}

func (k Kindly) getValidYConfig(n string) (yamlConfig, error) {
	var err error
	var yc yamlConfig

	// Pull out package version if provided
	nVer := strings.SplitN(n, "@", 2)

	dl := dlInfo{nVer[0], "", "", "", ""}

	if len(nVer) > 1 {
		dl.Version = semver.Canonical(nVer[1])
		if !semver.IsValid(dl.Version) {
			return yc, errors.New("Invalid package version: " + n)
		}
	}

	// Download package yaml spec and initialize yamlConfig struct
	if yc, err = getYaml(k.cfg.Source + dl.Name + ".yml"); err != nil {
		// TODO Write error message
		return yc, errors.New("ERROR")
	}

	// Check if package is available
	if !(len(yc.Spec.Name) > 0) {
		return yc, errors.New("Unavailable Package: " + dl.Name)
	}

	// Check if requested version is higher value than the available version in the package
	if len(dl.Version) > 0 {
		if semver.Compare(dl.Version, yc.Spec.Version) == 1 {
			return yc, errors.New("Version requested: " + n + "\tLatest version: " + dl.Name + "@" + yc.Spec.Version)
		}
	}

	// If version was not provided in the argument, set it to version in spec file
	if !(len(dl.Version) > 0) {
		dl.Version = yc.Spec.Version
	}

	// processFile Downloads file from url, checks SHA value, and saves it to tmpDir
	// TODO Should the requested OS ARCH be in config or request?
	dl.osArch = runtime.GOOS + "_" + runtime.GOARCH

	// Check if OS architecture is available
	if _, ok := yc.Spec.Assets[dl.osArch]; !ok {
		return yc, errors.New("Unavailable OS Architecture: " + dl.osArch)
	}

	return yc, nil
}
