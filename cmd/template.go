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

	kindly "github.com/borkod/kindly/pkg"
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

		kc, err := k.GenerateTemplate(ctx, owner, repo)
		if err != nil {
			fmt.Println(err)

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
