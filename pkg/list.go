package pkg

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

// List function implements list command
func (k Kindly) ListPackages(ctx context.Context) (s []string, err error) {

	client := github.NewClient(nil)
	source := k.cfg.Source
	source = "https://raw.githubusercontent.com/bojand/ghz/master/testdata/config"
	source = strings.Replace(source, "https://raw.githubusercontent.com/", "", 1)
	source = strings.Replace(source, "/master", "", 1)
	source = strings.TrimSuffix(source, "/")
	sInfo := strings.Split(source, "/")

	// list all organizations for user "willnorris"
	_, dir, _, err := client.Repositories.GetContents(ctx, sInfo[0], sInfo[1], sInfo[2], nil)
	if err != nil {
		return s, err
	}

	for _, n := range dir {
		s = append(s, strings.TrimSuffix(*n.Name, ".go"))
	}

	return s, nil
}
