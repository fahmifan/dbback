package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type AliOSS struct {
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type Config struct {
	OutDir   string `json:"outDir"`
	Postgres DB     `json:"postgres"`
	MySQL    DB     `json:"mysql"`
	AliOSS   AliOSS `json:"aliOSS"`
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
