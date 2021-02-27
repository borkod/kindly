package config

// Config struct for kindly
type Config struct {
	Verbose          bool
	OutBinDir        string
	OutCompletionDir string
	OutManDir        string
	UniqueDir        bool
	Completion       string
	Source           string
	OS               string
	Arch             string
}
