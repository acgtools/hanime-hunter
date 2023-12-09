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
	OutputDir string
	Info      bool
}

type ResolverOpt struct {
	Series   bool
	PlayList bool
}

func NewCfg() (*Config, error) {
	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
