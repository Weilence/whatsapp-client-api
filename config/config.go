package config

import "flag"

var (
	ShowV = flag.Bool("version", false, "show version")
	Env   = flag.String("env", "", "environment")
	Port  = flag.Int("port", 0, "port")
)
