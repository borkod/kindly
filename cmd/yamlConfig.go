package cmd

// YamlConfig is exported.
type YamlConfig struct {
	Spec struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Homepage    string `yaml:"homepage"`
		Version     string `yaml:"version"`
		Assets      map[string]struct {
			URL    string `yaml:"url"`
			ShaURL string `yaml:"sha_url"`
		}
		Bin        []string            `yaml:"bin"`
		Completion map[string][]string `yaml:"completion"`
		Man        []string            `yaml:"man"`
	}
}
