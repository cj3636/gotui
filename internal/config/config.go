package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	GithubToken      string           `yaml:"github_token"`
	GitlabToken      string           `yaml:"gitlab_token"`
	WeatherLocation  string           `yaml:"weather_location"`
	RefreshIntervals RefreshIntervals `yaml:"refresh_intervals"`
	GithubRepos      []string         `yaml:"github_repos"`
	GitlabProjects   []string         `yaml:"gitlab_projects"`
	TextFile         string           `yaml:"text_file"`
	MarkdownFile     string           `yaml:"markdown_file"`
	Layout           Layout           `yaml:"layout"`
}

// RefreshIntervals defines how often each widget refreshes (in seconds)
type RefreshIntervals struct {
	Weather int `yaml:"weather"`
	Github  int `yaml:"github"`
	Gitlab  int `yaml:"gitlab"`
	System  int `yaml:"system"`
	IP      int `yaml:"ip"`
}

// Layout defines the grid layout for widgets
type Layout struct {
	Rows int `yaml:"rows"`
	Cols int `yaml:"cols"`
}

// Load reads the configuration from a YAML file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.WeatherLocation == "" {
		config.WeatherLocation = "New York"
	}
	if config.Layout.Rows == 0 {
		config.Layout.Rows = 3
	}
	if config.Layout.Cols == 0 {
		config.Layout.Cols = 3
	}

	return &config, nil
}

// FindConfigFile looks for config.yaml in common locations
func FindConfigFile() (string, error) {
	// Try current directory
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml", nil
	}

	// Try home directory
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "gotui", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	// Return default path
	return "config.yaml", nil
}
