/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd is for implementing commands
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	kindly "github.com/borkod/kindly/pkg"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template [owner] [repo]",
	Short: "Generate a YAML spec template for a Github repo",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var k kindly.Kindly
		k.SetConfig(cfg)
		k.SetLogger(log.New(os.Stdout, "", log.Ltime))
		log.SetFlags(log.Ltime)

		if cfg.Verbose {
			log.Println("Generating template.")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		owner := args[0]
		repo := args[1]

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
			fmt.Println(err)
		}
		tags, _, err := client.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			fmt.Println(err)
		}

		release := tags[0]

		releaseInfo, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, release.GetName())
		if err != nil {
			fmt.Println(err)
		}

		var kc kindly.KindlyStruct

		kc.Spec.Name = repoInfo.GetName()
		kc.Spec.Description = repoInfo.GetDescription()
		kc.Spec.Homepage = repoInfo.GetHomepage()
		kc.Spec.RepoURL = repoInfo.GetHTMLURL()
		kc.Spec.Tags = repoInfo.Topics
		kc.Spec.License = repoInfo.GetLicense().GetSPDXID()
		kc.Spec.Version = release.GetName()
		kc.Spec.Assets = make(map[string]kindly.Asset)

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
							kc.Spec.Assets[goArch] = kindly.Asset{URL: "", ShaURL: ""}
						}
						if n.GetContentType() == "application/octet-stream" {
							kc.Spec.Assets[goArch] = kindly.Asset{URL: kc.Spec.Assets[goArch].URL, ShaURL: url}
						} else {
							kc.Spec.Assets[goArch] = kindly.Asset{URL: url, ShaURL: kc.Spec.Assets[goArch].ShaURL}
						}
					}
				}
			}
		}
		d, err := yaml.Marshal(&kc)
		if err != nil {
			log.Println("ERROR: ", err)
		}
		fmt.Printf("---\n%s\n", string(d))

	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
