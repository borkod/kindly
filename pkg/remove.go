package pkg

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Remove function implements remove command
func (k Kindly) Remove(ctx context.Context, p string) (err error) {
	filename := filepath.Join(k.cfg.ManifestDir, p+".yaml")
	l := new(pkgManifest)

	// Read package manifest
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, l); err != nil {
		return err
	}

	for _, n := range l.Bin {
		if k.cfg.Verbose {
			k.logger.Println("Deleting file: ", filepath.Join(k.cfg.OutBinDir, n))
		}
		if err := os.Remove(filepath.Join(k.cfg.OutBinDir, n)); err != nil {
			k.logger.Println("ERROR")
			k.logger.Println(err)
		}
	}

	for _, n := range l.Completion {
		if k.cfg.Verbose {
			k.logger.Println("Deleting file: ", filepath.Join(k.cfg.OutCompletionDir, n))
		}
		if err := os.Remove(filepath.Join(k.cfg.OutCompletionDir, n)); err != nil {
			k.logger.Println("ERROR")
			k.logger.Println(err)
		}
	}

	for _, n := range l.Man {
		if k.cfg.Verbose {
			k.logger.Println("Deleting file: ", filepath.Join(k.cfg.OutManDir, n))
		}
		if err := os.Remove(filepath.Join(k.cfg.OutManDir, n)); err != nil {
			k.logger.Println("ERROR")
			k.logger.Println(err)
		}
	}

	if k.cfg.Verbose {
		k.logger.Println("Deleting file: ", filename)
	}
	if err := os.Remove(filename); err != nil {
		k.logger.Println("ERROR")
		k.logger.Println(err)
	}

	return nil
}
