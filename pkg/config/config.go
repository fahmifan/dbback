package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DB struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	OutDir   string `json:"outDir"`
	Postgres DB     `json:"postgres"`
	MySQL    DB     `json:"mysql"`
}

func Load(cfgPath string) (Config, error) {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("open config (%s) : %w", cfgPath, err)
	}

	cfg := Config{}
	err = json.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode json: %w", err)
	}

	return cfg, nil
}
