package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const PATH = "../.env"

type Config struct {
	Env     string `env:"ENV"`
	HTTP    HTTP
	Storage Storage
}

type HTTP struct {
	Addr        string `env:"ADDR"`
	Timeout     string `env:"TIMEOUT"`
	IdleTimeout string `env:"IDLE_TIMEOUT"`
}

type Storage struct {
	Type string `env:"DB_TYPE"`
	URL  string `env:"DB_URL"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(PATH, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
