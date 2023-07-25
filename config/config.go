package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config -.
type Config struct {
	App `yaml:"app"`
	Log `yaml:"logger"`
}

// App -.
type App struct {
	PathToRAC string `env-required:"true" yaml:"path_to_rac" env:"PATH_TO_RAC"`
	PathTo1C  string `env-required:"true" yaml:"path_to_1cs" env:"PATH_TO_1C"`

	LockCode string `env-required:"true" yaml:"lock_code"`
}

// Log -.
type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
}

func New() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yml"
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
