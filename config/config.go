package config

// Config struct for kindly
type Config struct {
	Verbose          bool
	UniqueDir        bool
	OutBinDir        string
	OutCompletionDir string
	OutManDir        string
	Completion       string
	Source           string
	OS               string
	Arch             string
}
