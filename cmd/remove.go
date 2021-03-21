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
	"io/ioutil"
	"log"
	"os"
	"strings"

	kindly "github.com/borkod/kindly/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [name of package]",
	Short: "Removes a previously installed package.",
	Long: `Use to remove a previously installed package.

Optionally, use the --all flag to remove all installed packages.
If set, all other arguments are ignored.

Examples:
	kindly remove gh-cli
	kindly remove -a`,
	Run: func(cmd *cobra.Command, args []string) {
		var k kindly.Kindly
		k.SetConfig(cfg)
		k.SetLogger(log.New(os.Stdout, "", log.Ltime))
		log.SetFlags(log.Ltime)

		if !viper.GetBool("all") && len(args) == 0 {
			log.Fatalln("Must provide a package name as an argument.")
		}

		if viper.GetBool("all") {
			if len(args) > 0 {
				log.Println("Remove All flag is set. Ignoring all other arguments.")
			}
			args = make([]string, 0)

			files, err := ioutil.ReadDir(cfg.ManifestDir)
			if err != nil {
				log.Fatalln(err)
			}

			for _, n := range files {
				if strings.HasSuffix(n.Name(), ".yaml") {
					args = append(args, strings.TrimSuffix(n.Name(), ".yaml"))
				}
			}
		}
		// Iterate over all packages provided as command arguments
		for _, n := range args {

			if cfg.Verbose {
				log.Println("Removing package: ", n)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := k.Remove(ctx, n); err != nil {
				log.Print(string("\u001b[31m"), err, string("\u001b[0m"), "\n")
				continue
			}
		}

		if cfg.Verbose {
			log.Println("Removing complete.")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	removeCmd.Flags().BoolP("all", "a", false, "Removed all installed packages. If this flag is set all other arguments are ignored.")
	if err := viper.BindPFlag("all", removeCmd.Flags().Lookup("all")); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
