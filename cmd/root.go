/*
Copyright © 2021 Borko Djurkovic <borkod@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
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
)

var cfgFile string

// Verbose TODO SHould this be declared here?
var Verbose bool

// OutBinDir TODO SHould this be declared here?
var OutBinDir string

// OutCompletionDir TODO SHould this be declared here?
var OutCompletionDir string

// OutCompletionDir TODO SHould this be declared here?
var OutManDir string

// UniqueDir specifies if files should be saved into unique dir
var UniqueDir bool

// Completion specifies completion shell
var Completion string

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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kindly/.kindly.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&OutBinDir, "OutBinDir", "", "Default binary file output directory (default is $HOME/.local/bin/)")
	viper.BindPFlag("OutBinDir", rootCmd.PersistentFlags().Lookup("OutBinDir"))
	rootCmd.PersistentFlags().StringVar(&OutCompletionDir, "OutCompletionDir", "", "Default Completions file output directory (default is $HOME/.local/completion/)")
	viper.BindPFlag("OutCompletionDir", rootCmd.PersistentFlags().Lookup("OutCompletionDir"))
	rootCmd.PersistentFlags().StringVar(&OutManDir, "OutManDir", "", "Default Man Pages output directory (default is $HOME/.local/man/)")
	viper.BindPFlag("OutManDir", rootCmd.PersistentFlags().Lookup("OutManDir"))
	rootCmd.PersistentFlags().BoolVarP(&UniqueDir, "unique-directory", "", false, "write files into unique directory (default is false)")
	viper.BindPFlag("unique-directory", rootCmd.PersistentFlags().Lookup("unique-directory"))
	rootCmd.PersistentFlags().StringVar(&Completion, "completion", "bash", "Completion shell setting")
	viper.BindPFlag("completion", rootCmd.PersistentFlags().Lookup("completion"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	OutBinDir = filepath.Join(home, ".kindly", "bin")
	OutCompletionDir = filepath.Join(home, ".kindly", "completion")
	OutManDir = filepath.Join(home, ".kindly", "man")

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

	OutBinDir = viper.GetString("OutBinDir")
	UniqueDir = viper.GetBool("unique-directory")
	OutManDir = viper.GetString("OutManDir")
}
