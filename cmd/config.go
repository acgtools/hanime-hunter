package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Log         *LogConfig
	DLOpt       *DLOption
	ResolverOpt ResolverOpt
}

type LogConfig struct {
	Level string
}

type DLOption struct {
	OutputDir  string
	Quality    string
	Info       bool
	LowQuality bool
	Retry      uint8
}

type ResolverOpt struct {
	Series bool
}

func NewCfg() (*Config, error) {
	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
