package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TargetURL    string `json:"target_url"`
	TargetDomain string `json:"target_domain"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}
	return cfg
}
