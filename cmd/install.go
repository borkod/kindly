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
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Println("Installing files...")
		}

		// Create a temporary directory where files will be downloaded
		tmpDir, err := ioutil.TempDir("", "kindly_")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Clean up temporary directory
		defer os.RemoveAll(tmpDir)

		// Iterate over all packages provided as command arguments
		for _, n := range args {
			var tmpFile string
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

			// Applies Version values to the URL template
			if dl.URL, dl.URLSHA, err = executeURL(dl, yamlConfig); err != nil {
				continue
			}

			// Downloads package file and package SHA file.
			// Calculates package SHA value
			// Compares package SHA value to SHA value in the SHA file
			if tmpFile, err = processFile(dl, tmpDir); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				fmt.Println(err)
				continue
			}

			// decompress tmpFile into tmpDir
			if strings.Contains(tmpFile, "tar.gz") {
				if err = decompress(tmpDir, tmpFile); err != nil {
					// TODO Write error message
					fmt.Println("ERROR")
					continue
				}
			}

			if strings.Contains(tmpFile, "zip") {
				if _, err = Unzip(tmpFile, tmpDir); err != nil {
					// TODO Write error message
					fmt.Println("ERROR")
					continue
				}
			}

			// Copy all extracted bin files from tmpDir into OutBinDir
			for _, n = range yamlConfig.Spec.Bin {
				if strings.Contains(strings.ReplaceAll(n, " ", ""), "{{.OS}}") ||
					strings.Contains(strings.ReplaceAll(n, " ", ""), "{{.Arch}}") {
					if n, err = executeBin(n); err != nil {
						continue
					}
				}
				if runtime.GOOS == "windows" {
					n = n + ".exe"
				}
				if err = copyFile(OutBinDir, tmpDir, n); err != nil {
					// TODO Write error message
					fmt.Println("ERROR")
					fmt.Println(err)
				}
			}

			// Copy all extracted completion files from tmpDir into OutCompletionDir
			for _, n = range yamlConfig.Spec.Completion[Completion] {
				if err = copyFile(OutCompletionDir, tmpDir, n); err != nil {
					// TODO Write error message
					fmt.Println("ERROR")
					fmt.Println(err)
				}
			}

			// Copy all extracted man pages files from tmpDir into OutManDir
			for _, n = range yamlConfig.Spec.Man {
				if err = copyFile(OutManDir, tmpDir, n); err != nil {
					// TODO Write error message
					fmt.Println("ERROR")
					fmt.Println(err)
				}
			}
		}

		if Verbose {
			fmt.Println("\nInstalling files complete.")
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

// GetYaml downloads the yaml and configures the yamlConfig struct
func GetYaml(arg string) (YamlConfig, error) {
	var yamlConfig YamlConfig

	resp, err := http.Get("https://b3o-test-bucket.s3.ca-central-1.amazonaws.com/" + arg + ".yml")
	if err != nil {
		return yamlConfig, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		fmt.Printf("Error downloading file: %s\n", arg)
		return yamlConfig, err
	}

	yaml.Unmarshal(buf.Bytes(), &yamlConfig)

	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", arg)
		return yamlConfig, err
	}

	return yamlConfig, nil
}

// Downloads package file and package SHA file.
// Calculates package SHA value
// Compares package SHA value to SHA value in the SHA file
func processFile(dl dlInfo, tmpDir string) (string, error) {

	// Get the data
	if Verbose {
		fmt.Println("\nDownloading file:\t\t", dl.URL)
	}

	resp, err := http.Get(dl.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// DO I REALLY NEED TWO COPIES!?
	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)

	if _, err := io.Copy(w, resp.Body); err != nil {
		return "", err
	}

	if Verbose {
		fmt.Println("Download finished.")
	}

	// Calculate SHA256 of downloaded file
	hash := sha256.New()
	if _, err := io.Copy(hash, &buf1); err != nil {
		return "", err
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	if Verbose {
		fmt.Println("Calculated SHA256 value:\t", sum)
	}

	// Get the sha file
	if len(dl.URLSHA) > 1 {
		if Verbose {
			fmt.Println("Downloading SHA256 file:\t", dl.URLSHA)
		}

		respSha, err := http.Get(dl.URLSHA)
		if err != nil {
			return "", err
		}
		defer respSha.Body.Close()

		//buf := new(bytes.Buffer)
		newStr := ""
		//buf.ReadFrom(respSha.Body)
		scanner := bufio.NewScanner(respSha.Body)
		for scanner.Scan() {
			shaLine := strings.SplitN(scanner.Text(), " ", 2)
			if len(shaLine) > 1 {
				if strings.Contains(shaLine[1], runtime.GOOS) && strings.Contains(shaLine[1], runtime.GOARCH) {
					newStr = shaLine[0]
				}
			} else {
				newStr = shaLine[0]
			}
		}

		// Get the sha file
		if Verbose {
			fmt.Println("SHA256 file hash value:\t\t", newStr)
		}

		// Check if SHA256 values match
		if newStr != sum {
			fmt.Println("SHA MISMATCH")
			return "", err
		}
	} else if Verbose {
		fmt.Println("NO SHA FILE PROVIDED. SKIPPING SHA VALUE CHECK")
	}

	// Create the output file in temporary
	urlPath := strings.Split(dl.URL, "/")
	filepath := filepath.Join(tmpDir, urlPath[len(urlPath)-1])

	if Verbose {
		fmt.Println("Writing output file:\t\t", filepath)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, &buf2)
	return filepath, err
}

// Applies OS and Architecture values to the binary file names template
func executeBin(n string) (string, error) {
	binT, err := template.New("bin").Parse(n)

	if err != nil {
		// TODO Write error message
		fmt.Println("Error parsing binary: ", n)
		return "", err
	}

	type binS struct {
		OS   string
		Arch string
	}

	nS := binS{runtime.GOOS, runtime.GOARCH}

	var buf bytes.Buffer
	if err = binT.Execute(&buf, nS); err != nil {
		// TODO Write error message
		fmt.Println("Error parsing url: ", n)
		return "", err
	}
	newStr := buf.String()

	if runtime.GOOS == "windows" {
		newStr = newStr + ".exe"
	}
	return newStr, nil
}

// Applies Version values to the URL template
func executeURL(dl dlInfo, yc YamlConfig) (string, string, error) {
	urlT, err := template.New("url").Parse(yc.Spec.Assets[dl.osArch].URL)

	if err != nil {
		// TODO Write error message
		fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].URL)
		return "", "", err
	}

	urlShaT, err := template.New("urlSha").Parse(yc.Spec.Assets[dl.osArch].ShaURL)
	if err != nil {
		// TODO Write error message
		fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].ShaURL)
		return "", "", err
	}

	var buf bytes.Buffer
	if err = urlT.Execute(&buf, dl); err != nil {
		// TODO Write error message
		fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].ShaURL)
		return "", "", err
	}

	url := buf.String()

	buf.Reset()

	if err = urlShaT.Execute(&buf, dl); err != nil {
		// TODO Write error message
		fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].ShaURL)
		return "", "", err
	}
	urlSha := buf.String()
	return url, urlSha, nil
}
