package main

import (
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Database struct {
		DSN string `toml:"dsn"`
	} `toml:"database"`
	Runners struct {
		Logging struct {
			File string `toml:"file"`
		} `toml:"logging"`
		Email struct {
			Host     string `toml:"host"`
			Port     int    `toml:"port"`
			Email    string `toml:"email"`
			Password string `toml:"password"`
		} `toml:"email"`
	} `toml:"runners"`
}

func parseConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = toml.Unmarshal(f, config)

	return config, err
}

func GetConfig(paths []string) (*Config, error) {
	var config *Config

	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			continue
		}

		config, err = parseConfig(path)
		if err != nil {
			return nil, err
		}

		break
	}

	if config == nil {
		return nil, os.ErrNotExist
	}

	return config, nil
}
