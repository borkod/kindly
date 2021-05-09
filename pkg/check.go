package pkg

import (
	"context"
	"strings"

	"github.com/borkod/kindly/config"
	"github.com/google/go-github/github"
)

// Check function checks if the packages passed in args are available TODO variadic function
func (k Kindly) Check(ctx context.Context, s string, n string) (x bool, ks KindlyStruct, err error) {
	client := github.NewClient(nil)
	sources := k.cfg.Sources
	x = false

	if len(s) > 0 {
		sources = make(map[string]config.Source)
		sources[s] = k.cfg.Sources[s]
	}

	for key, source := range sources {
		if source.Type == "github" {
			_, dir, _, err := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
			if err != nil {
				return false, ks, err
			}

			for _, n := range dir {
				_, ks, err = k.getValidYConfig(ctx, key, strings.TrimSuffix(*n.Name, ".yaml"), false, false)
				if err != nil {
					return false, ks, err
				}

			}
		}
	}

	if len(ks.Spec.Name) > 0 {
		x = true
	}

	return x, ks, nil
}
