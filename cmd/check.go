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
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Println("Checking packages...")
		}

		// Iterate over all packages provided as command arguments
		for _, n := range args {
			var err error
			var yamlConfig YamlConfig

			// Pull out package version if provided
			nVer := strings.SplitN(n, "@", 2)

			dl := dlInfo{nVer[0], "", "", "", ""}

			if len(nVer) > 1 {
				dl.Version = semver.Canonical(nVer[1])
				if !semver.IsValid(dl.Version) {
					fmt.Println("Invalid package version: ", n)
					continue
				}
			}

			// Download package yaml spec and initialize yamlConfig struct
			if yamlConfig, err = GetYaml(dl.Name); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				continue
			}

			// Check if package is available
			if !(len(yamlConfig.Spec.Name) > 0) {
				fmt.Println("Unavailable Package: ", dl.Name)
				continue
			}

			// Check if requested version is higher value than the available version in the package
			if len(dl.Version) > 0 {
				if semver.Compare(dl.Version, yamlConfig.Spec.Version) == 1 {
					fmt.Println("Version requested: ", dl.Version, "Latest version: ", yamlConfig.Spec.Version)
					continue
				}
			}

			// If version was not provided in the argument, set it to version in spec file
			if !(len(dl.Version) > 0) {
				dl.Version = yamlConfig.Spec.Version
			}

			// processFile Downloads file from url, checks SHA value, and saves it to tmpDir
			dl.osArch = runtime.GOOS + "_" + runtime.GOARCH

			// Check if OS architecture is available
			if _, ok := yamlConfig.Spec.Assets[dl.osArch]; !ok {
				fmt.Println("Unavailable OS Architecture: ", dl.osArch)
				continue
			}

			// If no YAML output requested, just print one line. Otherwise print YAML.
			if !viper.GetBool("output") {
				fmt.Println("Package: ", dl.Name, "\t\tVersion: ", dl.Version)
			}

			// If YAML output requested, print complete spec YAML
			if viper.GetBool("output") {
				d, err := yaml.Marshal(&yamlConfig)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}
				fmt.Printf("%s\n\n", string(d))
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
