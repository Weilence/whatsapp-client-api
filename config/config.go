package config

import (
	"flag"
	"os"

	"log/slog"

	"github.com/BurntSushi/toml"
)

var (
	Env  = flag.String("env", "", "environment")
	Port = flag.Int("port", 0, "port")
)

var Config = &config{}

type config struct {
	Proxy string
}

func Parse() {
	if *Port == 0 {
		panic("port is required")
	}

	_, err := toml.DecodeFile("config.toml", Config)
	if err != nil {
		slog.Error("decode config", "err", err)
	}
}

func Save() {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("open config file", "err", err)
	}

	err = toml.NewEncoder(f).Encode(Config)
	if err != nil {
		slog.Error("encode config", "err", err)
	}
}
