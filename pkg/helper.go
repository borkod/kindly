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

// Package pkg is for implementing commands
package pkg

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// GetYaml downloads the yaml and configures the KindlyStruct struct
func getYaml(ctx context.Context, arg string) (KindlyStruct, error) {
	const ConnectMaxWaitTime = 1 * time.Second
	const RequestMaxWaitTime = 5 * time.Second

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}

	var yc KindlyStruct
	buf := new(bytes.Buffer)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, arg, nil)
	if err != nil {
		return yc, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return yc, err
	}
	defer resp.Body.Close()

	if _, err = buf.ReadFrom(resp.Body); err != nil {
		//fmt.Printf("Error downloading file: %s\n", arg)
		return yc, err
	}

	yaml.Unmarshal(buf.Bytes(), &yc)

	if err != nil {
		//fmt.Printf("Error parsing YAML file: %s\n", arg)
		return yc, err
	}

	return yc, nil
}

// decompress decompresses a file
func decompress(dst string, path string) error {

	//if cfg.Verbose {
	//	fmt.Println("Decompressing file:\t\t", path)
	//}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			//if cfg.Verbose {
			//	fmt.Println("Writing file:\t\t\t", target)
			//}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

// copyFile copies file to dst from src
func copyFile(dst string, src string, binName string) error {

	//if cfg.Verbose {
	//	fmt.Println("Copying file:\t\t", binName)
	//}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.Name() == binName && !info.IsDir() {

			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			// Create the output directory
			//if cfg.UniqueDir {
			//	dirPath := filepath.Join(dst, info.Name())
			//	if f, err := os.Stat(dirPath); os.IsNotExist(err) || !f.IsDir() {
			//		os.Mkdir(dirPath, os.ModePerm)
			//	}
			//	dst = dirPath
			//}

			newFile, err := os.Create(filepath.Join(dst, info.Name()))
			if err != nil {
				return err
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders within the zip file (parameter 1) to an output directory (parameter 2)
func Unzip(src string, dest string) ([]string, error) {

	//if cfg.Verbose {
	//	fmt.Println("Unzipping file:\t\t", src)
	//}
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// ExpandPath is helper function to expand file location
func ExpandPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		user, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(user.HomeDir, path[2:])
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abspath
}
