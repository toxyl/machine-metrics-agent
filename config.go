package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL       string `yaml:"url"`
	Org       string `yaml:"org"`
	Bucket    string `yaml:"bucket"`
	Token     string `yaml:"token"`
	Interval  int    `yaml:"interval"`
	VerifyTLS bool   `yaml:"verify_tls"`
}

func loadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
