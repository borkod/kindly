package pkg

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

// Update function implements update command
func (k Kindly) Update(ctx context.Context, n string) (err error) {
	filename := filepath.Join(k.cfg.ManifestDir, n+".yaml")
	l := new(pkgManifest)

	// Read package manifest
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, l); err != nil {
		return err
	}

	_, yc, err := k.getValidYConfig(ctx, l.Source, n, false, false)
	if err != nil {
		return err
	}

	if semver.Compare(l.Version, yc.Spec.Version) < 0 {
		if err := k.Install(ctx, l.Source, n, false, false); err != nil {
			return err
		}
	}

	return nil
}
