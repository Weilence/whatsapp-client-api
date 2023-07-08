package config

import (
	"flag"
	"log"
	"os"

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
		log.Fatalln("port is required")
	}

	_, err := toml.DecodeFile("config.toml", Config)
	if err != nil {
		log.Printf("decode config err: %v\n", err)
	}
}

func Save() {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("open config file err: %v", err)
	}

	err = toml.NewEncoder(f).Encode(Config)
	if err != nil {
		log.Printf("encode config err: %v", err)
	}
}
