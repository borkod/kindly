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

// Package cmd is for implementing kindly commands
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/borkod/kindly/config"
)

var cfgFile string

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kindly <command>",
	Version: "0.0.2-rc1",
	Short:   "Kindly installs Linux binaries.",
	Long: `Kindly is a free and open-source software package management system that simplifies the installation of software on Linux.

Kindly downloads pre-compiled binaries for available packages.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra persistent flags are defined here; global for the application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kindly/.kindly.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&cfg.OutBinDir, "OutBinDir", "", "Default binary file output directory (default is $HOME/.kindly/bin/)")
	viper.BindPFlag("OutBinDir", rootCmd.PersistentFlags().Lookup("OutBinDir"))
	rootCmd.PersistentFlags().StringVar(&cfg.ManifestDir, "ManifestDir", "", "Default kindly manifests directory (default is $HOME/.kindly/manifests/)")
	viper.BindPFlag("ManifestDir", rootCmd.PersistentFlags().Lookup("ManifestDir"))
	rootCmd.PersistentFlags().StringVar(&cfg.OutCompletionDir, "OutCompletionDir", "", "Default completions file output directory (default is $HOME/.kindly/completion/)")
	viper.BindPFlag("OutCompletionDir", rootCmd.PersistentFlags().Lookup("OutCompletionDir"))
	rootCmd.PersistentFlags().StringVar(&cfg.OutManDir, "OutManDir", "", "Default man pages output directory (default is $HOME/.kindly/man/)")
	viper.BindPFlag("OutManDir", rootCmd.PersistentFlags().Lookup("OutManDir"))
	//rootCmd.PersistentFlags().BoolVarP(&cfg.UniqueDir, "unique-directory", "", false, "write files into unique directory (default is false)")
	//viper.BindPFlag("unique-directory", rootCmd.PersistentFlags().Lookup("unique-directory"))
	rootCmd.PersistentFlags().StringVar(&cfg.Completion, "completion", "bash", "Completion shell setting")
	viper.BindPFlag("completion", rootCmd.PersistentFlags().Lookup("completion"))
	rootCmd.PersistentFlags().StringVar(&cfg.Source, "Source", "https://raw.githubusercontent.com/borkod/kindly-specs/main/specs/", "Source of package spec files")
	viper.BindPFlag("Source", rootCmd.PersistentFlags().Lookup("Source"))
	rootCmd.PersistentFlags().StringVar(&cfg.OS, "OS", "", "Operating System (default is current OS)")
	viper.BindPFlag("OS", rootCmd.PersistentFlags().Lookup("OS"))
	rootCmd.PersistentFlags().StringVar(&cfg.Arch, "Arch", "", "Architecture (default is current architecture)")
	viper.BindPFlag("Arch", rootCmd.PersistentFlags().Lookup("Arch"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize default values
	cfg.ManifestDir = filepath.Join(home, ".kindly", "manifests")
	cfg.OutBinDir = filepath.Join(home, ".kindly", "bin")
	cfg.OutCompletionDir = filepath.Join(home, ".kindly", "completion")
	cfg.OutManDir = filepath.Join(home, ".kindly", "man")
	cfg.OS = runtime.GOOS
	cfg.Arch = runtime.GOARCH

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".kindly" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".kindly"))
		viper.SetConfigName(".kindly")
	}

	viper.SetEnvPrefix("KINDLY")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Update variables based on any flags or environment variables set by the user
	cfg.ManifestDir = viper.GetString("ManifestDir")
	cfg.OutBinDir = viper.GetString("OutBinDir")
	//cfg.UniqueDir = viper.GetBool("unique-directory")
	cfg.OutManDir = viper.GetString("OutManDir")
	cfg.OS = viper.GetString("OS")
	cfg.Arch = viper.GetString("Arch")
}
