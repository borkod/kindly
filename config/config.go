package config

// Config struct for kindly
type Config struct {
	Verbose bool
	//	UniqueDir        bool
	ManifestDir      string
	OutBinDir        string
	OutCompletionDir string
	OutManDir        string
	Completion       string
	Sources          map[string]Source
	OS               string
	Arch             string
}

// Source is exported.
type Source struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	Owner         string `yaml:"owner"`
	Repo          string `yaml:"repo"`
	Path          string `yaml:"path"`
	DirectoryPath string `yaml:"directory_path"`
}
