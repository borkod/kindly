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

// KindlyStruct is exported.
type KindlyStruct struct {
	Spec struct {
		Name        string              `yaml:"name"`
		Description string              `yaml:"description"`
		Homepage    string              `yaml:"homepage"`
		RepoURL     string              `yaml:"repo_url"`
		License     string              `yaml:"license"`
		Tags        []string            `yaml:"tags"`
		Version     string              `yaml:"version"`
		Assets      map[string]Asset    `yaml:"assets"`
		Bin         []string            `yaml:"bin"`
		Completion  map[string][]string `yaml:"completion"`
		Man         []string            `yaml:"man"`
	}
}

// Asset is exported.
type Asset struct {
	URL    string `yaml:"url"`
	ShaURL string `yaml:"sha_url"`
}
