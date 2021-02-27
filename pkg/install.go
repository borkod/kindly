package pkg

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Install function implements install command
func (k Kindly) Install(ctx context.Context, args []string) {
	if k.cfg.Verbose {
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
		var yc KindlyStruct
		var dl dlInfo

		if dl, yc, err = k.getValidYConfig(ctx, n); err != nil {
			fmt.Println(err)
			continue
		}

		// Applies Version values to the URL template
		if dl.URL, dl.URLSHA, err = executeURL(dl, yc); err != nil {
			continue
		}

		// Downloads package file and package SHA file.
		// Calculates package SHA value
		// Compares package SHA value to SHA value in the SHA file
		if tmpFile, err = k.processFile(ctx, dl, tmpDir); err != nil {
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
				fmt.Println(err)
				continue
			}
		}

		if strings.Contains(tmpFile, "zip") {
			if _, err = Unzip(tmpFile, tmpDir); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				fmt.Println(err)
				continue
			}
		}

		// Copy all extracted bin files from tmpDir into OutBinDir
		for _, n = range yc.Spec.Bin {
			if strings.Contains(strings.ReplaceAll(n, " ", ""), "{{.OS}}") ||
				strings.Contains(strings.ReplaceAll(n, " ", ""), "{{.Arch}}") {
				if n, err = executeBin(n, k.cfg.OS, k.cfg.Arch); err != nil {
					continue
				}
			}
			if k.cfg.OS == "windows" {
				n = n + ".exe"
			}
			if err = copyFile(k.cfg.OutBinDir, tmpDir, n); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				fmt.Println(err)
			}
		}

		// Copy all extracted completion files from tmpDir into OutCompletionDir
		for _, n = range yc.Spec.Completion[k.cfg.Completion] {
			if err = copyFile(k.cfg.OutCompletionDir, tmpDir, n); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				fmt.Println(err)
			}
		}

		// Copy all extracted man pages files from tmpDir into OutManDir
		for _, n = range yc.Spec.Man {
			if err = copyFile(k.cfg.OutManDir, tmpDir, n); err != nil {
				// TODO Write error message
				fmt.Println("ERROR")
				fmt.Println(err)
			}
		}
	}

	if k.cfg.Verbose {
		fmt.Println("\nInstalling files complete.")
	}
}

// Downloads package file and package SHA file.
// Calculates package SHA value
// Compares package SHA value to SHA value in the SHA file
func (k Kindly) processFile(ctx context.Context, dl dlInfo, tmpDir string) (string, error) {

	// Get the data
	if k.cfg.Verbose {
		fmt.Println("\nDownloading file:\t\t", dl.URL)
	}

	const ConnectMaxWaitTime = 1 * time.Second
	const RequestMaxWaitTime = 5 * time.Second

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dl.URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
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

	if k.cfg.Verbose {
		fmt.Println("Download finished.")
	}

	// Calculate SHA256 of downloaded file
	hash := sha256.New()
	if _, err := io.Copy(hash, &buf1); err != nil {
		return "", err
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	if k.cfg.Verbose {
		fmt.Println("Calculated SHA256 value:\t", sum)
	}

	// Get the sha file
	if len(dl.URLSHA) > 1 {
		if k.cfg.Verbose {
			fmt.Println("Downloading SHA256 file:\t", dl.URLSHA)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dl.URLSHA, nil)
		if err != nil {
			return "", err
		}

		respSha, err := client.Do(req)
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
				if strings.Contains(shaLine[1], k.cfg.OS) && strings.Contains(shaLine[1], k.cfg.Arch) {
					newStr = shaLine[0]
				}
			} else {
				newStr = shaLine[0]
			}
		}

		// Get the sha file
		if k.cfg.Verbose {
			fmt.Println("SHA256 file hash value:\t\t", newStr)
		}

		// Check if SHA256 values match
		if newStr != sum {
			return "", errors.New("SHA MISMATCH")
		}
	} else if k.cfg.Verbose {
		fmt.Println("NO SHA FILE PROVIDED. SKIPPING SHA VALUE CHECK")
	}

	// Create the output file in temporary
	urlPath := strings.Split(dl.URL, "/")
	filepath := filepath.Join(tmpDir, urlPath[len(urlPath)-1])

	if k.cfg.Verbose {
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
func executeBin(n string, os string, arch string) (string, error) {
	binT, err := template.New("bin").Parse(n)

	if err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing binary: ", n)
		return "", err
	}

	type binS struct {
		OS   string
		Arch string
	}

	nS := binS{os, arch}

	var buf bytes.Buffer
	if err = binT.Execute(&buf, nS); err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing url: ", n)
		return "", err
	}
	newStr := buf.String()

	if os == "windows" {
		newStr = newStr + ".exe"
	}
	return newStr, nil
}

// Applies Version values to the URL template
func executeURL(dl dlInfo, yc KindlyStruct) (string, string, error) {
	urlT, err := template.New("url").Parse(yc.Spec.Assets[dl.osArch].URL)

	if err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].URL)
		return "", "", err
	}

	urlShaT, err := template.New("urlSha").Parse(yc.Spec.Assets[dl.osArch].ShaURL)
	if err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].ShaURL)
		return "", "", err
	}

	var buf bytes.Buffer
	if err = urlT.Execute(&buf, dl); err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].URL)
		return "", "", err
	}

	url := buf.String()

	buf.Reset()

	if err = urlShaT.Execute(&buf, dl); err != nil {
		// TODO Write error message
		//fmt.Println("Error parsing url: ", yc.Spec.Assets[dl.osArch].ShaURL)
		return "", "", err
	}
	urlSha := buf.String()
	return url, urlSha, nil
}
