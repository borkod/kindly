package pkg

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

// GenerateTemplate function implements template command
func (k Kindly) GenerateTemplate(ctx context.Context, owner string, repo string) (kc KindlyStruct, err error) {
	const goosList = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris windows zos"
	const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm x86_64"

	client := github.NewClient(nil)
	/*
		repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\nREPO INFO")
		fmt.Println(repoInfo.GetName())
		fmt.Println(repoInfo.GetDescription())
		fmt.Println(repoInfo.GetHTMLURL())
		fmt.Println(repoInfo.GetHomepage())
		fmt.Println(repoInfo.Topics)
		fmt.Println(repoInfo.GetLicense().GetSPDXID(), ": ", repoInfo.GetLicense().GetName())

		tags, _, err := client.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			fmt.Println(err)
		}

		release := tags[0]

		fmt.Println("\nRELEASE INFO")
		fmt.Println(release.GetName())

		releaseInfo, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, release.GetName())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\nASSETS")
		for _, n := range releaseInfo.Assets {
			fmt.Println(n.GetBrowserDownloadURL())
			fmt.Println(n.GetContentType())
			fmt.Println(n.GetName())
			fmt.Println()
		}
	*/
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return kc, err
	}
	tags, _, err := client.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return kc, err
	}

	release := tags[0]

	releaseInfo, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, release.GetName())
	if err != nil {
		return kc, err
	}

	//var kc kindly.KindlyStruct

	kc.Spec.Name = repoInfo.GetName()
	kc.Spec.Description = repoInfo.GetDescription()
	kc.Spec.Homepage = repoInfo.GetHomepage()
	kc.Spec.RepoURL = repoInfo.GetHTMLURL()
	kc.Spec.Tags = repoInfo.Topics
	kc.Spec.License = repoInfo.GetLicense().GetSPDXID()
	kc.Spec.Version = release.GetName()
	kc.Spec.Assets = make(map[string]Asset)

	for _, o := range strings.Split(goosList, " ") {
		for _, a := range strings.Split(goarchList, " ") {
			for _, n := range releaseInfo.Assets {
				url := n.GetBrowserDownloadURL()
				if strings.Contains(url, o) && strings.Contains(url, a) {
					if strings.Contains(url, kc.Spec.Version) {
						url = strings.ReplaceAll(url, kc.Spec.Version, "{{.Version}}")
					}
					goArch := o + "_" + a
					if a == "x86_64" {
						goArch = o + "_amd64"
					}
					if _, ok := kc.Spec.Assets[goArch]; !ok {
						kc.Spec.Assets[goArch] = Asset{URL: "", ShaURL: ""}
					}
					if n.GetContentType() == "application/octet-stream" {
						kc.Spec.Assets[goArch] = Asset{URL: kc.Spec.Assets[goArch].URL, ShaURL: url}
					} else {
						kc.Spec.Assets[goArch] = Asset{URL: url, ShaURL: kc.Spec.Assets[goArch].ShaURL}
					}
				}
			}
		}
	}

	return kc, nil

}
