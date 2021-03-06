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
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists available packages.",
	Long: `Lists available packages.
	
Example:
	kindly list`,
	Run: func(cmd *cobra.Command, args []string) {
		var k kindly.Kindly
		k.SetConfig(cfg)
		k.SetLogger(log.New(os.Stdout, "", log.Ltime))
		log.SetFlags(log.Ltime)

		if cfg.Verbose {
			log.Println("Listing available packages:")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		s, err := k.ListPackages(ctx, viper.GetBool("installed"))
		if err != nil {
			log.Println(err)

		}

		for _, n := range s {
			fmt.Println(n)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolP("installed", "i", false, "List locally installed packages.")
	if err := viper.BindPFlag("installed", listCmd.Flags().Lookup("installed")); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
