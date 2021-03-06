/*
Copyright © 2021 Borko Djurkovic <borkod@gmail.com>

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
	"log"
	"os"

	kindly "github.com/borkod/kindly/pkg"
	"github.com/spf13/cobra"
)

type dlInfo struct {
	Name    string
	Version string
	URL     string
	URLSHA  string
	osArch  string
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [name of package]",
	Short: "Installs one or many packages.",
	Long: `Installs one or many packages.

Example:
	kindly install gh-cli

You can use @ to specify a semantic version of the package.

Example:
	kindly install gh-cli@v1.0.0

You can provide multiple arguments to install multiple packages.
	
Example:
	kindly install gh-cli ghz`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var k kindly.Kindly
		k.SetConfig(cfg)
		k.SetLogger(log.New(os.Stdout, "", log.Ltime))
		log.SetFlags(log.Ltime)

		// Iterate over all packages provided as command arguments
		for _, n := range args {
			if cfg.Verbose {
				log.Println("Installing package: ", n)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := k.Install(ctx, n); err != nil {
				log.Print(string("\u001b[31m"), err, string("\u001b[0m"), "\n")
				continue
			}
		}

		if cfg.Verbose {
			log.Println("Installing files complete.")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
