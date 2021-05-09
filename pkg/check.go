package pkg

import (
	"context"
	"strings"

	"github.com/borkod/kindly/config"
	"github.com/google/go-github/github"
)

// Check function checks if the packages passed in args are available TODO variadic function
func (k Kindly) Check(ctx context.Context, s string, n string) (ks KindlyStruct, err error) {
	client := github.NewClient(nil)
	sources := k.cfg.Sources
	println(s)
	if len(s) > 0 {
		sources = make(map[string]config.Source)
		sources[s] = k.cfg.Sources[s]
	}

	for key, source := range sources {
		println(key)
		if source.Type == "github" {
			_, dir, _, err := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
			if err != nil {
				return ks, err
			}

			for _, n := range dir {
				_, ks, err = k.getValidYConfig(ctx, key, strings.TrimSuffix(*n.Name, ".yaml"), false, false)
				if err != nil {
					return ks, err
				}

			}
		}
	}

	return ks, nil
}
