package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server Server `yaml:"server"`
		App    App    `yaml:"app"`
		Http   Http   `yaml:"http"`
		Logger Logger `yaml:"logger"`
	}

	Server struct {
		Header string `yaml:"header"`
		Addr   string `env-required:"true" yaml:"addr" env:"SERVER_ADDR"`
	}

	App struct {
		Name      string `yaml:"name"`
		Version   string `yaml:"version"`
		RateLimit string `yaml:"rateLimit"  env:"RATE_LIMIT"`
	}

	Http struct {
		Addr string `env-required:"true" yaml:"addr" env:"HTTP_ADDR"`
	}

	Logger struct {
		Level string `env-required:"true" yaml:"level" env:"LOGGER_LEVEL"`
	}
)

func New() Config {
	dir := os.Getenv("CONF_DIR")

	if len(dir) == 0 {
		_, file, _, _ := runtime.Caller(0)
		dir = filepath.Dir(file)
	}

	var cfg Config

	err := cleanenv.ReadConfig(dir+"/config.yaml", &cfg)
	if err != nil {
		log.Fatalf("config load error: %s", err)
		panic(err)
	}

	return cfg
}
