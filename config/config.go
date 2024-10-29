package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Env    string `envconfig:"ENV" default:"local"`
		DB     DB
		Server Server
	}

	DB struct {
		URL string `envconfig:"DATABASE_URL" required:"true"`
	}

	Server struct {
		Port        string        `envconfig:"SERVER_PORT" default:"8080"`
		Timeout     time.Duration `envconfig:"SERVER_TIMEOUT" default:"4s"`
		IdleTimeout time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
	}
)

var Cfg Config

func Load() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		err := godotenv.Load(configPath)
		if err != nil {
			return fmt.Errorf("cannot load from config file: %s", err.Error())
		}
	}

	err := envconfig.Process("", &Cfg)
	if err != nil {
		return fmt.Errorf("cannot set config: %s", err.Error())
	}
	return nil
}
