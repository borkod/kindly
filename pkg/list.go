package pkg

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/borkod/kindly/config"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

// ListPackages function implements list command
func (k Kindly) ListPackages(ctx context.Context, source string, installed bool) (s []string, err error) {

	if installed {
		s, err = k.listInstalled(ctx)
	} else {
		s, err = k.listAvailable(ctx, source)
	}

	return s, err
}

func (k Kindly) listInstalled(ctx context.Context) (s []string, err error) {
	files, err := ioutil.ReadDir(k.cfg.ManifestDir)
	if err != nil {
		return s, err
	}

	// Should we do this or just use file names?
	for _, file := range files {
		l := new(pkgManifest)
		// Read package manifest
		file, err := ioutil.ReadFile(filepath.Join(k.cfg.ManifestDir, file.Name()))
		if err != nil {
			return s, err
		}

		if err := yaml.Unmarshal(file, l); err != nil {
			return s, err
		}

		s = append(s, l.Name+"@"+l.Version)
	}

	/*
		for _, file := range files {
			s = append(s, strings.TrimSuffix(file.Name(), ".yaml"))
		}
	*/

	return s, nil
}

func (k Kindly) listAvailable(ctx context.Context, source string) (s []string, err error) {
	client := github.NewClient(nil)
	sources := k.cfg.Sources
	if len(source) > 0 {
		sources = make(map[string]config.Source)
		sources[source] = k.cfg.Sources[source]
	}

	//source = strings.Replace(source, "https://raw.githubusercontent.com/", "", 1)
	//source = strings.Replace(source, "/main", "", 1)
	//source = strings.TrimSuffix(source, "/")
	//sInfo := strings.Split(source, "/")

	for key, source := range sources {
		if source.Type == "github" {
			_, dir, _, err := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
			if err != nil {
				return s, err
			}

			// Should we read the spec and get the name of the package from spec or use just file name?
			for _, n := range dir {
				_, yc, err := k.getValidYConfig(ctx, key, strings.TrimSuffix(*n.Name, ".yaml"), false, false)
				if err != nil {
					return s, err
				}
				s = append(s, key+"/"+yc.Spec.Name+"@"+yc.Spec.Version)
				//s = append(s, strings.TrimSuffix(*n.Name, ".yaml"))
			}
		}
	}

	return s, nil
}
