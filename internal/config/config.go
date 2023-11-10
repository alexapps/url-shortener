package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path"  env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"  env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}

// MustLoad - the config must be loaded otherwise - panic. Due to load on startup
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file %s does not exist: %v", configPath, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read file %s, %v", configPath, err)
	}

	return &cfg
}
