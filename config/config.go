package config

import (
	"log"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var mode = flag.StringP("mode", "m", "dev", "run mode")

func Init() {
	viper.AddConfigPath("config")
	viper.SetConfigType("yml")

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	viper.SetConfigName("config." + *mode)
	err = viper.MergeInConfig()
	if err != nil {
		log.Panic(err)
	}
}
