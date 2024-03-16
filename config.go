package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL      string `yaml:"url"`
	Org      string `yaml:"org"`
	Bucket   string `yaml:"bucket"`
	Token    string `yaml:"token"`
	Interval int    `yaml:"interval"`
}

func loadConfig() (*Config, error) {
	content, err := os.ReadFile("config.yaml")
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
