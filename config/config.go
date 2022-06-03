package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"whatsapp-client/pkg/utils"
)

var mode = flag.StringP("mode", "m", "dev", "run mode")

func Init() {
	viper.AddConfigPath("config")
	viper.SetConfigType("yml")

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	utils.NoError(err)

	viper.SetConfigName("config." + *mode)
	err = viper.MergeInConfig()
	utils.PrintError(err)
}
