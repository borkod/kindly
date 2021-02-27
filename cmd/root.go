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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/borkod/kindly/config"
)

// Config struct for kindly
/*type Config struct {
	Verbose          bool
	OutBinDir        string
	OutCompletionDir string
	OutManDir        string
	UniqueDir        bool
	Completion       string
}*/

var cfgFile string

// Verbose stores verbose flag value TODO SHould this be declared here?
var Verbose bool

// OutBinDir stores the configuration of the directory where binary files should be saved
var OutBinDir string

// OutCompletionDir stores the configuration of the directory where shell completion files should be saved
var OutCompletionDir string

// OutManDir stores the configuration of the directory where man pages files should be saved
var OutManDir string

// UniqueDir specifies if the binary files should be saved into their own unique dir
var UniqueDir bool

// Completion specifies completion shell configuration
var Completion string

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kindly",
	Version: "0.0.1",
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.PersistentFlags().StringVar(&cfg.OutBinDir, "OutBinDir", "", "Default binary file output directory (default is $HOME/.local/bin/)")
	viper.BindPFlag("OutBinDir", rootCmd.PersistentFlags().Lookup("OutBinDir"))
	rootCmd.PersistentFlags().StringVar(&cfg.OutCompletionDir, "OutCompletionDir", "", "Default Completions file output directory (default is $HOME/.local/completion/)")
	viper.BindPFlag("OutCompletionDir", rootCmd.PersistentFlags().Lookup("OutCompletionDir"))
	rootCmd.PersistentFlags().StringVar(&cfg.OutManDir, "OutManDir", "", "Default Man Pages output directory (default is $HOME/.local/man/)")
	viper.BindPFlag("OutManDir", rootCmd.PersistentFlags().Lookup("OutManDir"))
	rootCmd.PersistentFlags().BoolVarP(&cfg.UniqueDir, "unique-directory", "", false, "write files into unique directory (default is false)")
	viper.BindPFlag("unique-directory", rootCmd.PersistentFlags().Lookup("unique-directory"))
	rootCmd.PersistentFlags().StringVar(&cfg.Completion, "completion", "bash", "Completion shell setting")
	viper.BindPFlag("completion", rootCmd.PersistentFlags().Lookup("completion"))
	rootCmd.PersistentFlags().StringVar(&cfg.Source, "Source", "https://b3o-test-bucket.s3.ca-central-1.amazonaws.com/", "Source of packages (default is TODO ???")
	viper.BindPFlag("Source", rootCmd.PersistentFlags().Lookup("Source"))

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
	cfg.OutBinDir = filepath.Join(home, ".kindly", "bin")
	cfg.OutCompletionDir = filepath.Join(home, ".kindly", "completion")
	cfg.OutManDir = filepath.Join(home, ".kindly", "man")

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
	cfg.OutBinDir = viper.GetString("OutBinDir")
	cfg.UniqueDir = viper.GetBool("unique-directory")
	cfg.OutManDir = viper.GetString("OutManDir")
}
