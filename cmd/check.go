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
	"fmt"
	"log"
	"os"

	kindly "github.com/borkod/kindly/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [name of package]",
	Short: "Check if a package is available.",
	Long: `Check if a package and version is available for your OS.
	
Optionally, outputs the Kindly spec for the package.

Examples:
	kindly check gh-cli
	kindly check gh-cli -o`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var k kindly.Kindly
		k.SetConfig(cfg)
		k.SetLogger(log.New(os.Stdout, "", log.Ltime))
		log.SetFlags(log.Ltime)

		if cfg.Verbose {
			log.Println("Checking packages...")
		}

		// Iterate over all packages provided as command arguments
		for _, n := range args {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			yc, err := k.Check(ctx, n)
			if err != nil {
				log.Println("Package: ", n, string("\u001b[31m"), err, string("\u001b[0m"))
				continue
			}

			// If no YAML output requested, just print one line. Otherwise print YAML.
			if !viper.GetBool("output") {
				log.Println("Package: ", n, string("\u001b[32m"), "OK", string("\u001b[0m"))
			}

			// If YAML output requested, print complete spec YAML
			if viper.GetBool("output") {
				d, err := yaml.Marshal(&yc)
				if err != nil {
					log.Println("ERROR: ", err)
				}
				fmt.Println("# Package: ", n)
				fmt.Printf("---\n%s\n", string(d))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	checkCmd.Flags().BoolP("output", "o", false, "Output YAML spec file.")
	viper.BindPFlag("output", checkCmd.Flags().Lookup("output"))
}
