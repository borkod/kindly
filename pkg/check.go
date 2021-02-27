package pkg

import (
	"fmt"

	"github.com/borkod/kindly/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Kindly struct stores kindly config
type Kindly struct {
	cfg config.Config
}

// SetConfig sets the kindly struct config
func (k *Kindly) SetConfig(c config.Config) {
	k.cfg = c
}

// Check function checks if the packages passed in args are available
func (k Kindly) Check(args []string) {
	if k.cfg.Verbose {
		fmt.Println("Checking packages...")
	}

	// Iterate over all packages provided as command arguments
	for _, n := range args {
		var err error
		var yc yamlConfig

		if yc, err = k.getValidYConfig(n); err != nil {
			fmt.Println(err)
			continue
		}

		// If no YAML output requested, just print one line. Otherwise print YAML.
		if !viper.GetBool("output") {
			fmt.Println("Package: ", n)
		}

		// If YAML output requested, print complete spec YAML
		if viper.GetBool("output") {
			d, err := yaml.Marshal(&yc)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}
			fmt.Println("# Package: ", n)
			fmt.Printf("---\n%s\n", string(d))
		}
	}
}
